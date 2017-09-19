package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
)
//分享 博客 招聘 问答 框架 新闻 语言 数据库 外包 比赛
type TopicFactor struct {
	Id           int   `orm:"pk;auto"`
	ShareFactor  int   `orm:"default(10)"`
	BlogFactor   int   `orm:"default(10)"`
	WorkFactor   int   `orm:"default(10)"`
	QAAFactor    int   `orm:"default(10)"`
	FrameFactor  int   `orm:"default(10)"`
	NewsFactor   int   `orm:"default(10)"`
	LangFactor   int   `orm:"default(10)"`
	DBFactor     int   `orm:"default(10)"`
	OutBagFactor int   `orm:"default(10)"`
	MatchFactor  int   `orm:"default(10)"`
	Topic        *Topic `orm:"rel(fk)"`
}

func FindFactorByTopic(topic *Topic) TopicFactor {
	o := orm.NewOrm()
	var factor TopicFactor
	o.QueryTable(factor).RelatedSel().Filter("Topic", topic).One(&factor)
	return factor
}

func SaveTopicFactor(topicFactor *TopicFactor) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(topicFactor)
	return id
}

func UpdateTopicFactor(topicFactor *TopicFactor) {
	o := orm.NewOrm()
	o.Update(topicFactor)
}

func(TopicFactor)New(factorType int) *TopicFactor {
	switch factorType {
	case 1:
		factors := beego.AppConfig.String("default_share_factor_attribute")
		strings.Split(factors,",")
	}

	return nil
}