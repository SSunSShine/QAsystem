package main

import (
	"context"
	"fmt"
	"github.com/SSunSShine/QAsystem/conf"
	"github.com/SSunSShine/QAsystem/database"
	"github.com/SSunSShine/QAsystem/model"
)

var ctx = context.Background()

func initAll(conf *conf.Configuration)  {

	database.RDB.FlushDB(ctx)
	fmt.Println("redis flushDB.")

	if (database.DB.HasTable(&model.User{})) {
		fmt.Println("db has the table user, so drop it.")
		database.DB.DropTable(&model.User{})
	}

	if (database.DB.HasTable(&model.Profile{})) {
		fmt.Println("db has the table profile, so drop it.")
		database.DB.DropTable(&model.Profile{})
	}

	if (database.DB.HasTable(&model.Question{})) {
		fmt.Println("db has the table question, so drop it.")
		database.DB.DropTable(&model.Question{})
	}

	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Profile{})
	database.DB.AutoMigrate(&model.Question{})

	u0 := model.User{Mail: "123456@163.com", Password: "123456", Phone: "13212341234"}
	u0.Create()

	p0 := model.Profile{Name: "admin",Gender: 1, Desc: "This is the first account.", UserID: u0.ID}
	p0.Create()


	q1 := model.Question{Title: "First Question", Desc: "This is the first question !", UserID: u0.ID}
	q1.Create()


	fmt.Println("restarted success !")
}