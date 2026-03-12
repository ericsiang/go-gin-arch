package repository

import (
	"fmt"
	"self_go_gin/domains/user/entity/model"
	"self_go_gin/domains/user/repository/dao"

	"gorm.io/gorm"
)

// UserRepositoryInterface 用戶倉庫接口
type UserRepositoryInterface interface {
	GetDB() *gorm.DB
	GetUsersByAccount(account string) (*model.User, error)
	CreateUser(newUser *model.User) (*model.User, error)
}

// userRepositoryImpl 用戶倉庫實現
type userRepositoryImpl struct {
	dao dao.UserDaoInterface
}

// NewUserRepository 創建用戶倉庫
func NewUserRepository() (UserRepositoryInterface, error) {
	dao, err := dao.NewUserDao()
	if err != nil {
		return nil, fmt.Errorf("UserRepository NewUserRepository(): %w", err)
	}
	return &userRepositoryImpl{
		dao: dao,
	}, nil
}

func (r *userRepositoryImpl) GetDB() *gorm.DB {
	return r.dao.GetGenericDao().GetDB()
}

func (r *userRepositoryImpl) GetUsersByAccount(account string) (*model.User, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	user, err := r.dao.GetUsersByAccount(account)
	if err != nil {
		return nil, fmt.Errorf("UserRepositoryImpl GetUsersByAccount() data: %s \n %w", logData, err)
	}

	return user, err
}

func (r *userRepositoryImpl) CreateUser(newUser *model.User) (*model.User, error) {
	logData := map[string]interface{}{
		"newUser": newUser,
	}
	user, err := r.dao.GetGenericDao().Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("UserRepositoryImpl CreateUser() data: %s \n %w", logData, err)
	}
	return user, nil
}
