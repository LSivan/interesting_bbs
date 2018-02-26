package main

import (
	"git.oschina.net/gdou-geek-bbs/cron"
	"git.oschina.net/gdou-geek-bbs/models"
	_ "git.oschina.net/gdou-geek-bbs/routers"
	_ "git.oschina.net/gdou-geek-bbs/templates"
	_ "git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"git.oschina.net/gdou-geek-bbs/engine"
)

func init() {
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("jdbc.username")+":"+beego.AppConfig.String("jdbc.password")+"@/bbs?charset=utf8&parseTime=true&charset=utf8&loc=Asia%2FShanghai", 30)
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
		new(models.UserTopicList),
		new(models.TmpTopicFactor),
		new(models.TmpUserFactor),
	)
	orm.RunSyncdb("default", false, true)
}

func main() {
	// TODO 一些不变的数据使用redis，比如话题等
	// TODO 更多话题/回复/收藏
	// TODO 接入较成熟的全文搜索
	// TODO 导一份默认的sql
	// TODO README.md
	orm.Debug = true
	//ok, err := regexp.MatchString("/topic/edit/[0-9]+", "/topic/edit/123")
	//beego.Debug(ok, err)
	go cron.SetupCron()
	//_,user := models.FindUserById(2)
	//models.FindCollectTopicByUser(&user,7)
	go engine.Indexer.Index()
	beego.Run()
}
