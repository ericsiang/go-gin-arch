// Package repository 定義用戶倉庫接口和實現
package repository

import (
	"context"
	"self_go_gin/common/msgid"
	"self_go_gin/container"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/gin_application/handler"
	apperror "self_go_gin/internal/apperror"
	"time"

	"self_go_gin/domains/user/entity"
	"self_go_gin/domains/user/repository/dao"
	"self_go_gin/domains/user/repository/model"

	"github.com/redis/go-redis/v9"
)

// UserRepository 用戶接口
type UserRepository interface {
	GetUsersByAccount(account string) (*entity.User, error)
	CreateUser(newUser *entity.User) (*entity.User, error)
}

// userRepositoryImpl 用戶倉庫實現
type userRepositoryImpl struct {
	dao dao.UserDaoInterface
}

// NewUserRepository 創建用戶倉庫
func NewUserRepository() (UserRepository, error) {
	dao, err := dao.NewUserDao()
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.InitFailed,
			err,
			apperror.WithLayer("UserRepository NewUserRepository()"),
		)
	}
	return &userRepositoryImpl{
		dao: dao,
	}, nil
}

// GetUsersByAccount 根據帳號查詢用戶
func (r *userRepositoryImpl) GetUsersByAccount(account string) (*entity.User, error) {
	// 先確認緩存
	app := container.GetContainer()
	rdb := app.GetRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	exist, err := rdb.Get(ctx, "user:exists:"+account).Result()
	if err != nil && err != redis.Nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.InitFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl GetUsersByAccount() Check Redis"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}
	// 表示 redis 內帳號已存在
	if exist != "" {
		return nil, apperror.NewAppError(
			msgid.ResourceExist,
			appmsg.CheckAccountExistFailed,
			handler.ErrResourceExist,
			apperror.WithLayer("UserRepositoryImpl GetUsersByAccount() redis get()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}

	userModel, err := r.dao.GetUsersByAccount(account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.QueryFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl GetUsersByAccount() GetUsersByAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}

	user, err := r.modelToDomain(userModel)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.DataConversionFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl GetUsersByAccount() modelToDomain()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = rdb.SetNX(ctx, "user:exists:"+account, 1, 10*time.Minute).Result()
	if err != nil && err != redis.Nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.RedisSetFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl GetUsersByAccount() SetNX()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}

	return user, nil
}

// CreateUser 創建用戶
func (r *userRepositoryImpl) CreateUser(newUser *entity.User) (*entity.User, error) {
	userModel := r.domainToModel(newUser)
	// 儲存到資料庫
	createdPO, err := r.dao.Create(userModel)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CreateFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl CreateUser() Create()"),
			apperror.WithLogData(map[string]interface{}{
				"account": newUser.GetAccount(),
			}),
		)
	}

	user, err := r.modelToDomain(createdPO)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.DataConversionFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl CreateUser() modelToDomain()"),
			apperror.WithLogData(map[string]interface{}{
				"account": newUser.GetAccount(),
			}),
		)
	}

	app := container.GetContainer()
	rdb := app.GetRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = rdb.SetNX(ctx, "user:exists:"+newUser.GetAccount(), 1, 10*time.Minute).Result()
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CreateFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl CreateUser() SetNX()"),
			apperror.WithLogData(map[string]interface{}{
				"account": newUser.GetAccount(),
			}),
		)
	}

	return user, nil
}

// ============ 轉換方法（私有） ============

// domainToModel 領域模型轉換為持久化物件
func (r *userRepositoryImpl) domainToModel(user *entity.User) *model.User {
	return &model.User{
		GormModel: user.GormModel,
		Account:   user.GetAccount(),
		Password:  user.GetPasswordHash(),
	}
}

// modelToDomain 持久化物件轉換為領域模型
func (r *userRepositoryImpl) modelToDomain(model *model.User) (*entity.User, error) {
	// 重建值物件
	account, err := valueobj.NewAccount(model.Account)
	if err != nil {
		// 資料庫中的資料應該是有效的，如果出錯可能是資料損壞
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.DataConversionFailed,
			err,
			apperror.WithLayer("UserRepositoryImpl modelToDomain() NewAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account_from_db": model.Account,
			}),
		)
	}

	password := valueobj.NewPasswordFromHash(model.Password)

	// 重建聚合根
	user := entity.ReconstructUser(model.ID, account, password, model.GormModel)

	return user, nil
}
