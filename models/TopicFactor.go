package models

import (
	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sort"
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
func FindTopicFactorById(id int) TopicFactor {
	o := orm.NewOrm()
	var factor TopicFactor
	o.QueryTable(factor).RelatedSel().Filter("Id", id).One(&factor)
	return factor
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

func UpdateTopicFactorByMap(factorMap map[string]int, topicId int) {
	o := orm.NewOrm()
	var b bytes.Buffer
	b.WriteString("update topic_factor set ")
	for factor, value := range factorMap {
		b.WriteString(factor)
		b.WriteString(" = ")
		b.WriteString(factor)
		b.WriteString(" + (")
		b.WriteString(strconv.Itoa(value))
		b.WriteString("),")
	}
	b.WriteString(" topic_id = ")
	b.WriteString(strconv.Itoa(topicId))
	b.WriteString(" where ")
	b.WriteString(" topic_id = ")
	b.WriteString(strconv.Itoa(topicId))
	o.Raw(b.String()).Exec()
}

/*
	reply.Up = reply.Up + 1
	o.Update(reply, "Up")
*/
func UpdateTopicFactorByTmpFactor(tmpTopicFactor *TmpTopicFactor,topicFactor *TopicFactor) {
	o := orm.NewOrm()
	topicFactor.ShareFactor = topicFactor.ShareFactor + tmpTopicFactor.ShareFactor
	topicFactor.BlogFactor = topicFactor.BlogFactor + tmpTopicFactor.BlogFactor
	topicFactor.WorkFactor = topicFactor.WorkFactor + tmpTopicFactor.WorkFactor
	topicFactor.QAAFactor = topicFactor.QAAFactor + tmpTopicFactor.QAAFactor
	topicFactor.FrameFactor = topicFactor.FrameFactor + tmpTopicFactor.FrameFactor
	topicFactor.NewsFactor = topicFactor.NewsFactor + tmpTopicFactor.NewsFactor
	topicFactor.LangFactor = topicFactor.LangFactor + tmpTopicFactor.LangFactor
	topicFactor.DBFactor = topicFactor.DBFactor + tmpTopicFactor.DBFactor
	topicFactor.OutBagFactor = topicFactor.OutBagFactor + tmpTopicFactor.OutBagFactor
	topicFactor.MatchFactor = topicFactor.MatchFactor + tmpTopicFactor.MatchFactor
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

type topicFactorValue struct {
	Factors []string
	Values  []int
}

func (s *topicFactorValue) Len() int {
	return len(s.Factors)
}

// Swap is part of sort.Interface.
func (s *topicFactorValue) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
	s.Factors[i], s.Factors[j] = s.Factors[j], s.Factors[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *topicFactorValue) Less(i, j int) bool {
	return s.Values[i] < s.Values[j]
}

// 0获取五项最高的因子(特征因子),1获取五项最低的因子(无关因子)
func (uf TopicFactor) GetTopFactorByType(factorType int) map[string]int {
	f := topicFactorValue{}
	factors := []string{"share_factor", "blog_factor", "work_factor", "q_a_a_factor", "frame_factor", "news_factor", "lang_factor", "d_b_factor", "out_bag_factor", "match_factor"}
	values := []int{uf.ShareFactor, uf.BlogFactor, uf.WorkFactor, uf.QAAFactor, uf.FrameFactor, uf.NewsFactor, uf.LangFactor, uf.DBFactor, uf.OutBagFactor, uf.MatchFactor}
	f.Factors = factors
	f.Values = values
	sort.Sort(&f)
	factorMap := make(map[string]int)
	switch factorType {
	case 0:
		for i, val := range f.Factors[5:10] {
			factorMap[val] = f.Values[i+5]
		}
	case 1:
		for i, val := range f.Factors[:5] {
			factorMap[val] = f.Values[i]
		}
	default:
		return nil
	}
	return factorMap
}
