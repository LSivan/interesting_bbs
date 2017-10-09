package cron

import (
	"github.com/robfig/cron"
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego/logs"
)
var c * cron.Cron

func init() {
	c = cron.New()
}

func SetupCron(){
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, changeTopicFactor)
	c.AddFunc(spec, changeUserFactor)
	c.Start()

	select{}
}
var changeTopicFactor = func() {
	list := models.FindTopicFactorChangeSum()
	for _,v := range list {
		logs.Debug("v.TopicFactor.Id : ",v.TopicFactor.Id,"v : ",v)
	}
}

var changeUserFactor = func() {
	list := models.FindUserFactorChangeSum()
	for _,v := range list {
		logs.Debug("v.UserFactor.Id : ",v.UserFactor.Id,"v : ",v)
	}
}