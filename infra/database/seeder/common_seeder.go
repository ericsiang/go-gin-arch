// Package seeder 提供資料庫種子數據的創建功能，用於在開發和測試環境中快速生成初始數據。
package seeder

import (
	admin_model "self_go_gin/domains/admin/entity/model"
	user_model "self_go_gin/domains/user/entity/model"
	"self_go_gin/infra/orm/gorm_mysql"
	"self_go_gin/util/bcryptencap"
	"strconv"
)

// CreateUser 創建用戶資料
func CreateUser() {
	db, err := gormysql.GetMysqlDB()
	if err != nil {
		panic(err)
	}
	seeder := NewSeeder(db)
	if err := seeder.Clear("users"); err != nil {
		panic(err)
	}
	var users []*user_model.User
	//密碼加密
	bcryptPassword, err := bcryptencap.GenerateFromPassword("123456")
	if err != nil {
		panic("Seeder CreateUser() bcrypt fail")
	}

	for i := 1; i < 4; i++ {
		users = append(users, &user_model.User{
			Account:  "user" + strconv.Itoa(i),
			Password: string(bcryptPassword),
		})
	}

	err = db.Create(&users).Error
	if err != nil {
		panic("Seeder CreateUser() Create fail")
	}
}

// CreateAdmin 創建管理員資料
func CreateAdmin() {
	db, err := gormysql.GetMysqlDB()
	if err != nil {
		panic(err)
	}
	seeder := NewSeeder(db)
	if err := seeder.Clear("admins"); err != nil {
		panic(err)
	}
	var admins []*admin_model.Admins
	//密碼加密
	bcryptPassword, err := bcryptencap.GenerateFromPassword("123456")
	if err != nil {
		panic("Seeder CreateAdmin() bcrypt fail")
	}

	for i := 1; i < 4; i++ {
		admins = append(admins, &admin_model.Admins{
			Account:  "admin" + strconv.Itoa(i),
			Password: string(bcryptPassword),
		})
	}

	err = db.Create(&admins).Error
	if err != nil {
		panic("Seeder CreateUser() Create fail")
	}
}
