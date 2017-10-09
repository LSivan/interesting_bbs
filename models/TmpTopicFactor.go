package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"bytes"
	"github.com/astaxie/beego/logs"
)

//分享 博客 招聘 问答 框架 新闻 语言 数据库 外包 比赛
type TmpTopicFactor struct {
	Id           int          `orm:"pk;auto"`
	ShareFactor  int          `orm:"default(10)"`
	BlogFactor   int          `orm:"default(10)"`
	WorkFactor   int          `orm:"default(10)"`
	QAAFactor    int          `orm:"default(10)"`
	FrameFactor  int          `orm:"default(10)"`
	NewsFactor   int          `orm:"default(10)"`
	LangFactor   int          `orm:"default(10)"`
	DBFactor     int          `orm:"default(10)"`
	OutBagFactor int          `orm:"default(10)"`
	MatchFactor  int          `orm:"default(10)"`
	TopicFactor  *TopicFactor `orm:"rel(fk)"`
}

func FindTopicFactorChangeSum() []TmpTopicFactor {
	o := orm.NewOrm()
	s := "select sum(share_factor) share_factor," +
		"sum(blog_factor) blog_factor," +
		"sum(work_factor) work_factor," +
		"sum(q_a_a_factor) q_a_a_factor," +
		"sum(frame_factor) frame_factor," +
		"sum(news_factor) news_factor," +
		"sum(lang_factor) lang_factor," +
		"sum(d_b_factor) d_b_factor," +
		"sum(out_bag_factor) out_bag_factor," +
		"sum(match_factor) match_factor," +
		"topic_factor_id " +
		"from tmp_topic_factor group by topic_factor_id"
	var factors []TmpTopicFactor
	_, err := o.Raw(s).QueryRows(&factors)
	if err != nil {
		logs.Debug("FindUserFactorChangeSumByFactorId: ", err)
	}
	return factors
}

func SaveTmpTopicFactorByMap(factorValue map[string]int,topicFactorId int) int64 {
	o := orm.NewOrm()
	var b bytes.Buffer
	b.WriteString("insert into tmp_topic_factor SET id = null")
	for factor,value := range factorValue {
		b.WriteString(", ")
		b.WriteString(factor)
		b.WriteString(" = ")
		b.WriteString(strconv.Itoa(value))
	}
	b.WriteString(", topic_factor_id = ")
	b.WriteString(strconv.Itoa(topicFactorId))
	res, err := o.Raw(b.String()).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		return num
	}
	return 0
}
