// Package repository 定義管理員 Repository 介面和實現
package repository

import (
	"fmt"
	"gorm.io/gorm"
	"self_go_gin/domains/admin/entity/model"
	"self_go_gin/domains/admin/repository/dao"
)

// AdminRepositoryInterface 管理員帳號密碼表 Repository 介面
type AdminRepositoryInterface interface {
	GetDB() *gorm.DB
	GetAdminByAccount(account string) (*model.Admins, error)
	CreateAdmin(newAdmin model.Admins) (*model.Admins, error)
}

type adminRepositoryImpl struct {
	dao dao.AdminDaoInterface
}

// NewAdminRepository 建立管理員帳號密碼表 Repository
func NewAdminRepository() (AdminRepositoryInterface, error) {
	dao, err := dao.NewAdminDao()
	if err != nil {
		return nil, fmt.Errorf("AdminRepository NewAdminRepository() : %w", err)
	}

	return &adminRepositoryImpl{
		dao: dao,
	}, nil
}

func (r *adminRepositoryImpl) GetDB() *gorm.DB {
	return r.dao.GetGenericDao().GetDB()
}

func (r *adminRepositoryImpl) GetAdminByAccount(account string) (*model.Admins, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	fmt.Println("logData:", logData)
	admin, err := r.dao.GetAdminByAccount(account)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl GetAdminByAccount() data: %s \n %w", logData, err)
	}

	return admin, err
}

func (r *adminRepositoryImpl) CreateAdmin(newAdmin model.Admins) (*model.Admins, error) {
	logData := map[string]interface{}{
		"newAdmin": newAdmin,
	}
	admin, err := r.dao.GetGenericDao().Create(&newAdmin)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl CreateAdmin() data: %s \n %w", logData, err)
	}
	return admin, nil
}
