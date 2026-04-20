// Package dao 定義用戶數據訪問對象和接口
package dao

import (
	"fmt"
	"self_go_gin/common/msgid"
	"self_go_gin/container"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/user/repository/model"
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
		return nil, apperror.NewAppErrorWithLogData(
			msgid.Fail,
			appmsg.DAODatabaseConnectionFailed,
			fmt.Errorf("failed to get gorm.DB instance"),
			nil,
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
		// 記錄不存在不是錯誤，直接返回
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		// 只包裝真正的數據庫錯誤
		return nil, apperror.NewAppErrorWithLogData(
			msgid.Fail,
			appmsg.DAOQueryRecordsFailed,
			fmt.Errorf("UserDaoImpl GetUsersByAccount() account: %s, error: %w", account, err),
			map[string]interface{}{
				"account": account,
			},
		)
	}
	return &userPO, nil
}

// Create 創建用戶
func (d *userDaoImpl) Create(userPO *model.User) (*model.User, error) {
	err := d.GenericDao.GetDB().Create(userPO).Error
	if err != nil {
		return nil, apperror.NewAppErrorWithLogData(
			msgid.Fail,
			appmsg.DAOCreateRecordFailed,
			fmt.Errorf("UserDaoImpl Create() error: %w", err),
			map[string]interface{}{
				"userPO": userPO,
			},
		)
	}
	return userPO, nil
}
