// Package dao 定義用戶數據訪問對象和接口
package dao

import (
	"fmt"
	"self_go_gin/container"
	"self_go_gin/internal/dao"
)

// UserDaoInterface 用戶數據訪問接口
// DAO 層職責：純粹的資料庫操作，使用 PO（持久化物件）
type UserDaoInterface interface {
	GetGenericDao() dao.GenericDaoInterface[UserPO]
	GetUsersByAccount(account string) (*UserPO, error)
	Create(userPO *UserPO) (*UserPO, error)
}

type userDaoImpl struct {
	GenericDao dao.GenericDaoInterface[UserPO]
}

// NewUserDao 創建用戶數據訪問對象
func NewUserDao() (UserDaoInterface, error) {
	app :=container.GetContainer()
	db :=app.GetDB()
	return &userDaoImpl{
		GenericDao: dao.NewGenericDAO[UserPO](db),
	}, nil
}

func (d *userDaoImpl) GetGenericDao() dao.GenericDaoInterface[UserPO] {
	return d.GenericDao
}

// GetUsersByAccount 根據帳號查詢用戶
func (d *userDaoImpl) GetUsersByAccount(account string) (*UserPO, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	var userPO UserPO
	err := d.GenericDao.GetDB().Where("account = ?", account).First(&userPO).Error
	if err != nil {
		return nil, fmt.Errorf("UserDaoImpl GetUsersByAccount() data: %s \n %w", logData, err)
	}
	return &userPO, nil
}

// Create 創建用戶
func (d *userDaoImpl) Create(userPO *UserPO) (*UserPO, error) {
	logData := map[string]interface{}{
		"userPO": userPO,
	}
	err := d.GenericDao.GetDB().Create(userPO).Error
	if err != nil {
		return nil, fmt.Errorf("UserDaoImpl Create() data: %s \n %w", logData, err)
	}
	return userPO, nil
}
