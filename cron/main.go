package cron

import (
	"git.oschina.net/gdou-geek-bbs/models"
	//"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	"git.oschina.net/gdou-geek-bbs/recommend"
	"github.com/astaxie/beego"
)

var c *cron.Cron

func init() {
	c = cron.New()
}

func SetupCron() {
	//spec := "02 02 04 * * ?" // 每天凌晨的4:02:02进行因子的变化
	beego.BeeLogger.Info("每天凌晨的4:02:02进行因子的变化")
	spec := "*/10 * * * * ?"
	c.AddFunc(spec, recommend.GetUsersFavoriteList)
	//c.AddFunc(spec, changeUserFactor)
	c.Start()
	select {}
}


var changeTopicFactor = func() {
	var updateTopicFactor = func(tmpTopicFactor models.TmpTopicFactor){
		topicFactor := models.FindTopicFactorById(tmpTopicFactor.TopicFactor.Id)
		models.UpdateTopicFactorByTmpFactor(&tmpTopicFactor,&topicFactor)
	}

	list := models.FindTopicFactorChangeSum()
	for _, v := range list {
		//logs.Debug("v.TopicFactor.Id : ", v.TopicFactor.Id, "v : ", v)
		updateTopicFactor(v)
	}
	models.ClearTmpTopicFactor() // 把临时表的数据清掉,避免第二天用旧的数据进行修改傻花
}

var changeUserFactor = func() {
	var updateUserFactor = func(tmpUserFactor models.TmpUserFactor){
		userFactor := models.FindUserFactorById(tmpUserFactor.UserFactor.Id)
		models.UpdateUserFactorByTmpFactor(&tmpUserFactor,&userFactor)
	}

	list := models.FindUserFactorChangeSum()
	for _, v := range list {
		//logs.Debug("v.UserFactor.Id : ", v.UserFactor.Id, "v : ", v)
		updateUserFactor(v)
	}
	models.ClearTmpUserFactor()
}
