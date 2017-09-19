package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

//分享 博客 招聘 问答 框架 新闻 语言 数据库 外包 比赛
type TopicFactor struct {
	Id           int    `orm:"pk;auto"`
	ShareFactor  int    `orm:"default(10)"`
	BlogFactor   int    `orm:"default(10)"`
	WorkFactor   int    `orm:"default(10)"`
	QAAFactor    int    `orm:"default(10)"`
	FrameFactor  int    `orm:"default(10)"`
	NewsFactor   int    `orm:"default(10)"`
	LangFactor   int    `orm:"default(10)"`
	DBFactor     int    `orm:"default(10)"`
	OutBagFactor int    `orm:"default(10)"`
	MatchFactor  int    `orm:"default(10)"`
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

func (TopicFactor) New(factorType int) *TopicFactor {
	mustInt := func(val string) int {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
		return 0
	}
	fn := func(factors []string) *TopicFactor {
		factor := TopicFactor{}
		factor.ShareFactor = mustInt(factors[0])
		factor.BlogFactor = mustInt(factors[1])
		factor.WorkFactor = mustInt(factors[2])
		factor.QAAFactor = mustInt(factors[3])
		factor.FrameFactor = mustInt(factors[4])
		factor.NewsFactor = mustInt(factors[5])
		factor.LangFactor = mustInt(factors[6])
		factor.DBFactor = mustInt(factors[7])
		factor.OutBagFactor = mustInt(factors[8])
		factor.MatchFactor = mustInt(factors[9])
		return &factor
	}
	switch factorType {
	case 1:
		factors := strings.Split(beego.AppConfig.String("constant.default_share_factor_attribute"), ",")
		return fn(factors)
	case 2:
		factors := strings.Split(beego.AppConfig.String("constant.default_blog_factor_attribute"), ",")
		return fn(factors)
	case 3:
		factors := strings.Split(beego.AppConfig.String("constant.default_work_factor_attribute"), ",")
		return fn(factors)
	case 4:
		factors := strings.Split(beego.AppConfig.String("constant.default_QAA_factor_attribute"), ",")
		return fn(factors)
	case 5:
		factors := strings.Split(beego.AppConfig.String("constant.default_frame_factor_attribute"), ",")
		return fn(factors)
	case 6:
		factors := strings.Split(beego.AppConfig.String("constant.default_news_factor_attribute"), ",")
		return fn(factors)
	case 7:
		factors := strings.Split(beego.AppConfig.String("constant.default_lang_factor_attribute"), ",")
		return fn(factors)
	case 8:
		factors := strings.Split(beego.AppConfig.String("constant.default_DB_factor_attribute"), ",")
		return fn(factors)
	case 9:
		factors := strings.Split(beego.AppConfig.String("constant.default_outbag_factor_attribute"), ",")
		return fn(factors)
	case 10:
		factors := strings.Split(beego.AppConfig.String("constant.default_match_factor_attribute"), ",")
		return fn(factors)
	default:
		return nil
	}
}
