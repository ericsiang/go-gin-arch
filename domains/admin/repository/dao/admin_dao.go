// Package dao 定義管理員帳號密碼表的數據訪問對象和接口
package dao

import (
	"fmt"
	"self_go_gin/container"
	"self_go_gin/internal/dao"
)

// AdminDao 管理員帳號密碼表 DAO 介面
type AdminDao interface {
	GetGenericDao() dao.GenericDaoInterface[AdminPO]
	GetAdminByAccount(account string) (*AdminPO, error)
	Create(adminPO *AdminPO) (*AdminPO, error)
}

type adminDaoImpl struct {
	GenericDao dao.GenericDaoInterface[AdminPO]
}

// NewAdminDao 建立管理員帳號密碼表 DAO
func NewAdminDao() (AdminDao, error) {
	app := container.GetContainer()
	db := app.GetDB()
	return &adminDaoImpl{
		GenericDao: dao.NewGenericDAO[AdminPO](db),
	}, nil
}

// GetGenericDao 取得通用 DAO 實例
func (d *adminDaoImpl) GetGenericDao() dao.GenericDaoInterface[AdminPO] {
	return d.GenericDao
}

// GetAdminByAccount 根據帳號查詢管理員
func (d *adminDaoImpl) GetAdminByAccount(account string) (*AdminPO, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	var adminPO AdminPO
	err := d.GenericDao.GetDB().Where("account = ?", account).First(&adminPO).Error
	if err != nil {
		return nil, fmt.Errorf("AdminDaoImpl GetAdminByAccount() data: %s \n %w", logData, err)
	}
	return &adminPO, nil
}

// Create 創建管理員
func (d *adminDaoImpl) Create(adminPO *AdminPO) (*AdminPO, error) {
	logData := map[string]interface{}{
		"adminPO": adminPO,
	}
	err := d.GenericDao.GetDB().Create(adminPO).Error
	if err != nil {
		return nil, fmt.Errorf("AdminDaoImpl Create() data: %s \n %w", logData, err)
	}
	return adminPO, nil
}
