// Package migrate 負責數據庫結構的自動遷移，確保數據庫表結構與應用程式中的模型定義保持一致。
package migrate

import (
	"self_go_gin/container"
	admin_model "self_go_gin/domains/admin/repository/model"
	user_model "self_go_gin/domains/user/repository/model"

	"gorm.io/gorm"
)

// Migrate 自動遷移數據庫結構
func Migrate() {
	app := container.GetContainer()
	appDB := app.GetMySQLDB()
	db, ok := appDB.GetDB().(*gorm.DB)
	if !ok {
		panic("failed to get gorm.DB instance")
	}
	err := db.AutoMigrate(&user_model.User{})
	panicErr(err)
	err = db.AutoMigrate(&admin_model.Admin{})
	panicErr(err)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
