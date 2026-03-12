package dao

import (
	"fmt"
	"self_go_gin/domains/user/entity/model"
	"self_go_gin/infra/orm/gorm_mysql"
	"self_go_gin/internal/dao"
)

// UserDaoInterface 用戶數據訪問接口
type UserDaoInterface interface {
	GetGenericDao() dao.GenericDaoInterface[model.User]
	GetUsersByAccount(account string) (*model.User, error)
}

type userDaoImpl struct {
	GenericDao dao.GenericDaoInterface[model.User]
}

// NewUserDao 創建用戶數據訪問對象
func NewUserDao() (UserDaoInterface, error) {
	db ,err :=gorm_mysql.GetMysqlDB()
	if err != nil {
		return nil, fmt.Errorf("UserDao NewUserDao() GetMysqlDB fail : %w", err)
	}
	return &userDaoImpl{
		GenericDao: dao.NewGenericDAO[model.User](db),
	}, nil
}

func (d *userDaoImpl) GetGenericDao() dao.GenericDaoInterface[model.User] {
	return d.GenericDao
}
func (d *userDaoImpl) GetUsersByAccount(account string) (*model.User, error) {
	logData := map[string]interface{}{
		"account": account,
	}
	var user model.User
	err := d.GenericDao.GetDB().Where("account =?", account).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("UserDaoImpl GetUsersByAccount() data: %s \n %w", logData, err)
	}
	return &user, err
}
