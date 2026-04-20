// Package repository 定義管理員 Repository 介面和實現
package repository

import (
	"fmt"
	"self_go_gin/common/msgid"

	"self_go_gin/domains/admin/entity"
	"self_go_gin/domains/admin/repository/dao"
	"self_go_gin/domains/admin/repository/model"
	"self_go_gin/domains/common/valueobj"
	apperror "self_go_gin/internal/apperror"

	"gorm.io/gorm"
)

// AdminRepository 管理員帳號密碼表接口
type AdminRepository interface {
	GetDB() *gorm.DB
	GetAdminByAccount(account string) (*entity.Admin, error)
	CreateAdmin(newAdmin *entity.Admin) (*entity.Admin, error)
}

type adminRepositoryImpl struct {
	dao dao.AdminDao
}

// NewAdminRepository 建立管理員帳號密碼表 Repository
func NewAdminRepository() (AdminRepository, error) {
	dao, err := dao.NewAdminDao()
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			"初始化管理員倉庫失敗",
			fmt.Errorf("AdminRepository NewAdminRepository() : %w", err),
			nil,
		)
	}

	return &adminRepositoryImpl{
		dao: dao,
	}, nil
}

func (r *adminRepositoryImpl) GetDB() *gorm.DB {
	return r.dao.GetGenericDao().GetDB()
}

// GetAdminByAccount 根據帳號查詢管理員
func (r *adminRepositoryImpl) GetAdminByAccount(account string) (*entity.Admin, error) {
	// 從 DAO 層取得 PO
	adminPO, err := r.dao.GetAdminByAccount(account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			"查詢管理員失敗",
			fmt.Errorf("AdminRepositoryImpl GetAdminByAccount() account: %s, error: %w", account, err),
			map[string]interface{}{
				"account": account,
			},
		)
	}

	// PO -> 領域模型轉換
	admin, err := r.modelToDomain(adminPO)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			"管理員數據轉換失敗",
			fmt.Errorf("AdminRepositoryImpl GetAdminByAccount() convert PO to domain failed: %w", err),
			map[string]interface{}{
				"adminPO": adminPO,
			},
		)
	}

	return admin, nil
}

// CreateAdmin 創建管理員
func (r *adminRepositoryImpl) CreateAdmin(newAdmin *entity.Admin) (*entity.Admin, error) {
	// 領域模型 -> PO 轉換
	adminPO := r.domainTomodel(newAdmin)

	// 儲存到資料庫
	createdPO, err := r.dao.Create(adminPO)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			"建立管理員失敗",
			fmt.Errorf("AdminRepositoryImpl CreateAdmin() error: %w", err),
			map[string]interface{}{
				"adminPO": adminPO,
			},
		)
	}

	// PO -> 領域模型轉換
	admin, err := r.modelToDomain(createdPO)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			"管理員數據轉換失敗",
			fmt.Errorf("AdminRepositoryImpl CreateAdmin() convert PO to domain failed: %w", err),
			map[string]interface{}{
				"createdPO": createdPO,
			},
		)
	}

	return admin, nil
}

// ============ 轉換方法（私有） ============

// domainTomodel 領域模型轉換為持久化物件
func (r *adminRepositoryImpl) domainTomodel(admin *entity.Admin) *model.Admin {
	return &model.Admin{
		GormModel: admin.GormModel,
		Account:   admin.GetAccount(),
		Password:  admin.GetPasswordHash(),
	}
}

// modelToDomain 持久化物件轉換為領域模型
func (r *adminRepositoryImpl) modelToDomain(po *model.Admin) (*entity.Admin, error) {
	// 重建值物件
	account, err := valueobj.NewAccount(po.Account)
	if err != nil {
		// 資料庫中的資料應該是有效的，如果出錯可能是資料損壞
		return nil, apperror.NewAppError(
			msgid.Fail,
			"管理員數據損壞",
			fmt.Errorf("invalid account in database: %w", err),
			map[string]interface{}{
				"account": po.Account,
			},
		)
	}

	password := valueobj.NewPasswordFromHash(po.Password)

	// 重建聚合根
	admin := entity.ReconstructAdmin(po.ID, account, password, po.GormModel)

	return admin, nil
}
