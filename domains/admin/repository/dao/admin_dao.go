// Package dao 定義管理員帳號密碼表的數據訪問對象和接口
package dao

import (
	"fmt"
	"self_go_gin/container"
	"self_go_gin/domains/admin/repository/model"
	"self_go_gin/internal/dao"

	"gorm.io/gorm"
)

// AdminDao 管理員帳號密碼表 DAO 介面
type AdminDao interface {
	GetGenericDao() dao.GenericDaoInterface[model.Admin]
	GetAdminByAccount(account string) (*model.Admin, error)
	Create(adminPO *model.Admin) (*model.Admin, error)
}

type adminDaoImpl struct {
	GenericDao dao.GenericDaoInterface[model.Admin]
}

// NewAdminDao 建立管理員帳號密碼表 DAO
func NewAdminDao() (AdminDao, error) {
	app := container.GetContainer()
	appDB := app.GetMySQLDB()
	db, ok := appDB.GetDB().(*gorm.DB) // 獲取底層 DB 實例
	if !ok {
		return nil, fmt.Errorf("failed to get gorm.DB instance")
	}
	return &adminDaoImpl{
		GenericDao: dao.NewGenericDAO[model.Admin](db),
	}, nil
}

// GetGenericDao 取得通用 DAO 實例
func (d *adminDaoImpl) GetGenericDao() dao.GenericDaoInterface[model.Admin] {
	return d.GenericDao
}

// GetAdminByAccount 根據帳號查詢管理員
func (d *adminDaoImpl) GetAdminByAccount(account string) (*model.Admin, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	var adminPO model.Admin
	err := d.GenericDao.GetDB().Where("account = ?", account).First(&adminPO).Error
	if err != nil {
		return nil, fmt.Errorf("AdminDaoImpl GetAdminByAccount() data: %s \n %w", logData, err)
	}
	return &adminPO, nil
}

// Create 創建管理員
func (d *adminDaoImpl) Create(adminPO *model.Admin) (*model.Admin, error) {
	logData := map[string]interface{}{
		"adminPO": adminPO,
	}
	err := d.GenericDao.GetDB().Create(adminPO).Error
	if err != nil {
		return nil, fmt.Errorf("AdminDaoImpl Create() data: %s \n %w", logData, err)
	}
	return adminPO, nil
}
