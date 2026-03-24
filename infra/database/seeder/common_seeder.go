// Package seeder 提供資料庫種子數據的創建功能，用於在開發和測試環境中快速生成初始數據。
package seeder

import (
	"self_go_gin/container"
	"self_go_gin/domains/common/valueobj"

	admin_model "self_go_gin/domains/admin/repository/model"
	user_model "self_go_gin/domains/user/repository/model"
	"strconv"
)

// CreateUser 創建用戶資料
func CreateUser() {
	app := container.GetContainer()
	db := app.GetDB()
	seeder := NewSeeder(db)
	if err := seeder.Clear("users"); err != nil {
		panic(err)
	}
	var users []*user_model.User

	// 使用 DDD 方式創建用戶
	for i := 1; i < 4; i++ {
		account, err := valueobj.NewAccount("user" + strconv.Itoa(i))
		if err != nil {
			panic("Seeder CreateUser() create account fail: " + err.Error())
		}
		password, err := valueobj.NewPasswordFromPlainText("123456")
		if err != nil {
			panic("Seeder CreateUser() create password fail: " + err.Error())
		}

		user := &user_model.User{
			Account:  account.String(),
			Password: password.Hash(),
		}
		users = append(users, user)
	}

	err := db.Create(&users).Error
	if err != nil {
		panic("Seeder CreateUser() Create fail")
	}
}

// CreateAdmin 創建管理員資料
func CreateAdmin() {
	app := container.GetContainer()
	db := app.GetDB()
	seeder := NewSeeder(db)
	if err := seeder.Clear("admins"); err != nil {
		panic(err)
	}
	var admins []*admin_model.Admin

	// 使用 DDD 方式創建管理員
	for i := 1; i < 4; i++ {
		account, err := valueobj.NewAccount("admin" + strconv.Itoa(i))
		if err != nil {
			panic("Seeder CreateAdmin() create account fail: " + err.Error())
		}
		password, err := valueobj.NewPasswordFromPlainText("123456")
		if err != nil {
			panic("Seeder CreateAdmin() create password fail: " + err.Error())
		}
		admin := &admin_model.Admin{
			Account:  account.String(),
			Password: password.Hash(),
		}
		admins = append(admins, admin)
	}

	err := db.Create(&admins).Error
	if err != nil {
		panic("Seeder CreateUser() Create fail")
	}
}
