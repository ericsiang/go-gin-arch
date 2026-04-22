// Package dao 定義用戶數據訪問對象和接口
package dao

import (
	"self_go_gin/common/msgid"
	"self_go_gin/container"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/user/repository/model"
	"self_go_gin/gin_application/handler"
	apperror "self_go_gin/internal/apperror"
	"self_go_gin/internal/dao"

	"gorm.io/gorm"
)

// UserDaoInterface 用戶數據訪問接口
// DAO 層職責：純粹的資料庫操作，使用 PO（持久化物件）
type UserDaoInterface interface {
	GetGenericDao() dao.GenericDaoInterface[model.User]
	GetUsersByAccount(account string) (*model.User, error)
	Create(userPO *model.User) (*model.User, error)
}

type userDaoImpl struct {
	GenericDao dao.GenericDaoInterface[model.User]
}

// NewUserDao 創建用戶數據訪問對象
func NewUserDao() (UserDaoInterface, error) {
	app := container.GetContainer()
	appDB := app.GetMySQLDB()
	db, ok := appDB.GetDB().(*gorm.DB)
	if !ok || db == nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.DatabaseConnectionFailed,
			handler.ErrGetDBFailed,
			apperror.WithLayer("UserDAO NewUserDao()"),
		)
	}
	return &userDaoImpl{
		GenericDao: dao.NewGenericDAO[model.User](db),
	}, nil
}

func (d *userDaoImpl) GetGenericDao() dao.GenericDaoInterface[model.User] {
	return d.GenericDao
}

// GetUsersByAccount 根據帳號查詢用戶
func (d *userDaoImpl) GetUsersByAccount(account string) (*model.User, error) {
	var userPO model.User
	err := d.GenericDao.GetDB().Where("account = ?", account).First(&userPO).Error
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.QueryFailed,
			err,
			apperror.WithLayer("UserDAO GetUsersByAccount"),
			apperror.WithLogData(map[string]interface{}{
				"account": account,
			}),
		)
	}
	return &userPO, nil
}

// Create 創建用戶
func (d *userDaoImpl) Create(userPO *model.User) (*model.User, error) {
	err := d.GenericDao.GetDB().Create(userPO).Error
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CreateFailed,
			err,
			apperror.WithLayer("UserDAO Create()"),
			apperror.WithLogData(map[string]interface{}{
				"userPO": userPO,
			}),
		)
	}
	return userPO, nil
}
