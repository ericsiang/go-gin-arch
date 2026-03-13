// Package migrate 負責數據庫結構的自動遷移，確保數據庫表結構與應用程式中的模型定義保持一致。
package migrate

import (
	admin_model "self_go_gin/domains/admin/entity/model"
	user_model "self_go_gin/domains/user/entity/model"
	"self_go_gin/infra/orm/gorm_mysql"
)

// Migrate 自動遷移數據庫結構
func Migrate() {
	db, err := gormysql.GetMysqlDB()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&user_model.User{})
	panicErr(err)
	err = db.AutoMigrate(&admin_model.Admins{})
	panicErr(err)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
