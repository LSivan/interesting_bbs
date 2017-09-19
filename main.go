package main

import (
	_ "git.oschina.net/gdou-geek-bbs/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"git.oschina.net/gdou-geek-bbs/models"
	_ "github.com/go-sql-driver/mysql"
	_ "git.oschina.net/gdou-geek-bbs/utils"
	_ "git.oschina.net/gdou-geek-bbs/templates"
)

func init(){
    orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("jdbc.username") + ":" + beego.AppConfig.String("jdbc.password") + "@/pybbs-go?charset=utf8&parseTime=true&charset=utf8&loc=Asia%2FShanghai", 30)
    orm.RegisterModel(
		new(models.User),
		new(models.Topic),
		new(models.Section),
		new(models.Reply),
		new(models.ReplyUpLog),
        new(models.Role),
        new(models.Permission),
		new(models.UserFactor),
		new(models.TopicFactor),
	)
    orm.RunSyncdb("default", false, true)
}

func main() {
    orm.Debug = true
    //ok, err := regexp.MatchString("/topic/edit/[0-9]+", "/topic/edit/123")
    //beego.Debug(ok, err)

	beego.Run()
}

