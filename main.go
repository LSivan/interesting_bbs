package main

import (
	"git.oschina.net/gdou-geek-bbs/engine"
	"git.oschina.net/gdou-geek-bbs/models"
	_ "git.oschina.net/gdou-geek-bbs/routers"
	_ "git.oschina.net/gdou-geek-bbs/templates"
	_ "git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
	"git.oschina.net/gdou-geek-bbs/common"
	"git.oschina.net/gdou-geek-bbs/recommend"
	"git.oschina.net/gdou-geek-bbs/cron"
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
	 orm.Debug = true // 开启数据库日志

	// 对话题文章进行索引
	go engine.Indexer.Index()
	// 捕获ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			beego.BeeLogger.Info("捕获Ctrl-c，退出服务器\n")
			engine.Indexer.Exit <- struct{}{}
			<- engine.Indexer.Fin
			os.Exit(0)
		}
	}()
	// redis服务
	common.SetupRedis()
	// 提取文章特征值
	recommend.InitTopicFeature()
	// 提取用户特征值
	recommend.InitUserFeature()
	// 将用户的喜好列表存到redis中
	//recommend.GetUsersFavoriteList()
	// 更改用户兴趣的定时器
	go cron.SetupCron()
	// http服务
	beego.Run()
}
