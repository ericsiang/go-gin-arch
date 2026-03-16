// Package repository 定義管理員 Repository 介面和實現
package repository

import (
	"fmt"
	"self_go_gin/domains/admin/entity/model"
	"self_go_gin/domains/admin/repository/dao"
	"self_go_gin/domains/common/valueobj"

	"gorm.io/gorm"
)

// AdminRepository 管理員帳號密碼表接口
type AdminRepository interface {
	GetDB() *gorm.DB
	GetAdminByAccount(account string) (*model.Admins, error)
	CreateAdmin(newAdmin *model.Admins) (*model.Admins, error)
}

type adminRepositoryImpl struct {
	dao dao.AdminDao
}

// NewAdminRepository 建立管理員帳號密碼表 Repository
func NewAdminRepository() (AdminRepository, error) {
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

// GetAdminByAccount 根據帳號查詢管理員
func (r *adminRepositoryImpl) GetAdminByAccount(account string) (*model.Admins, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	
	// 從 DAO 層取得 PO
	adminPO, err := r.dao.GetAdminByAccount(account)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl GetAdminByAccount() data: %s \n %w", logData, err)
	}

	// PO -> 領域模型轉換
	admin, err := r.poToDomain(adminPO)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl GetAdminByAccount() convert PO to domain failed: %w", err)
	}

	return admin, nil
}

// CreateAdmin 創建管理員
func (r *adminRepositoryImpl) CreateAdmin(newAdmin *model.Admins) (*model.Admins, error) {
	logData := map[string]interface{}{
		"newAdmin": newAdmin,
	}
	
	// 領域模型 -> PO 轉換
	adminPO := r.domainToPO(newAdmin)
	
	// 儲存到資料庫
	createdPO, err := r.dao.Create(adminPO)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl CreateAdmin() data: %s \n %w", logData, err)
	}
	
	// PO -> 領域模型轉換
	admin, err := r.poToDomain(createdPO)
	if err != nil {
		return nil, fmt.Errorf("AdminRepositoryImpl CreateAdmin() convert PO to domain failed: %w", err)
	}
	
	return admin, nil
}

// ============ 轉換方法（私有） ============

// domainToPO 領域模型轉換為持久化物件
func (r *adminRepositoryImpl) domainToPO(admin *model.Admins) *dao.AdminPO {
	return &dao.AdminPO{
		GormModel: admin.GormModel,
		Account:   admin.GetAccount(),
		Password:  admin.GetPasswordHash(),
	}
}

// poToDomain 持久化物件轉換為領域模型
func (r *adminRepositoryImpl) poToDomain(po *dao.AdminPO) (*model.Admins, error) {
	// 重建值物件
	account, err := valueobj.NewAccount(po.Account)
	if err != nil {
		// 資料庫中的資料應該是有效的，如果出錯可能是資料損壞
		return nil, fmt.Errorf("invalid account in database: %w", err)
	}
	
	password := valueobj.NewPasswordFromHash(po.Password)
	
	// 重建聚合根
	admin := model.ReconstructAdmins(po.ID, account, password, po.GormModel)
	
	return admin, nil
}
