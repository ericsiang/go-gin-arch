// Package dao 定義用戶數據訪問對象和接口
package dao

import (
	"fmt"
	"self_go_gin/container"
	"self_go_gin/domains/user/repository/model"
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
	if !ok {
		return nil, fmt.Errorf("failed to get gorm.DB instance")
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
	logData := map[string]interface{}{
		"account": account,
	}
	var userPO model.User
	err := d.GenericDao.GetDB().Where("account = ?", account).First(&userPO).Error
	if err != nil {
		return nil, fmt.Errorf("UserDaoImpl GetUsersByAccount() data: %s \n %w", logData, err)
	}
	return &userPO, nil
}

// Create 創建用戶
func (d *userDaoImpl) Create(userPO *model.User) (*model.User, error) {
	logData := map[string]interface{}{
		"userPO": userPO,
	}
	err := d.GenericDao.GetDB().Create(userPO).Error
	if err != nil {
		return nil, fmt.Errorf("UserDaoImpl Create() data: %s \n %w", logData, err)
	}
	return userPO, nil
}
