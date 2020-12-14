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

	if (database.DB.HasTable(&model.Answer{})) {
		fmt.Println("db has the table answer, so drop it.")
		database.DB.DropTable(&model.Answer{})
	}

	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Profile{})
	database.DB.AutoMigrate(&model.Question{})
	database.DB.AutoMigrate(&model.Answer{})

	u0 := model.User{Mail: "123456@163.com", Password: "123456", Phone: "13212341234"}
	u0.Create()

	p0 := model.Profile{Name: "admin",Gender: 1, Desc: "This is the first account.", UserID: u0.ID}
	p0.Create()


	q0 := model.Question{Title: "First Question", Desc: "This is the first question !", UserID: u0.ID}
	q0.Create()

	a0 := model.Answer{Content: "This is the first answer!", QuestionID: q0.ID, UserID: u0.ID}
	a0.Create()

	fmt.Println("restarted success !")
}