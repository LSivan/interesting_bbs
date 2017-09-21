package models

import (
	"github.com/astaxie/beego/orm"
	"sort"
	"bytes"
	"strconv"
)

//分享 博客 招聘 问答 框架 新闻 语言 数据库 外包 比赛
type UserFactor struct {
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
	User         *User `orm:"rel(fk)"`
}

func FindFactorByUser(user *User) UserFactor {
	o := orm.NewOrm()
	var factor UserFactor
	o.QueryTable(factor).RelatedSel().Filter("User", user).One(&factor)
	return factor
}

func SaveUserFactor(userFactor *UserFactor) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(userFactor)
	return id
}

func UpdateUserFactor(userFactor *UserFactor) {
	o := orm.NewOrm()
	o.Update(userFactor)
}

func (UserFactor) New() *UserFactor {
	return &UserFactor{
		ShareFactor : 10,
		BlogFactor  :10,
		WorkFactor  :10,
		QAAFactor    :10,
		FrameFactor  :10,
		NewsFactor   :10,
		LangFactor   :10,
		DBFactor     :10,
		OutBagFactor :10,
		MatchFactor  :10,
	}
}

type userFactorValue struct {
	Factors []string
	Values []int
}

func (s *userFactorValue) Len() int {
	return len(s.Factors)
}

// Swap is part of sort.Interface.
func (s *userFactorValue) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
	s.Factors[i], s.Factors[j] = s.Factors[j], s.Factors[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *userFactorValue) Less(i, j int) bool {
	return s.Values[i]< s.Values[j]
}

// 0获取五项最高的因子(特征因子),1获取五项最低的因子(无关因子)
func (uf UserFactor) GetTopFactorByType(factorType int) map[string]int {
	f := userFactorValue{}
	factors := []string{"share_factor","blog_factor","work_factor","q_a_a_factor","frame_factor","news_factor","lang_factor","d_b_factor","out_bag_factor","match_factor"}
	values := []int{uf.ShareFactor,uf.BlogFactor,uf.WorkFactor,uf.QAAFactor,uf.FrameFactor,uf.NewsFactor,uf.LangFactor,uf.DBFactor,uf.OutBagFactor,uf.MatchFactor}
	f.Factors = factors
	f.Values = values
	sort.Sort(&f)
	factorMap := make(map[string]int)
	switch factorType {
	case 0:
		for i,val := range f.Factors[5:] {
			factorMap[val] = f.Values[i]
		}
	case 1:
		for i,val := range f.Factors[:5] {
			factorMap[val] = f.Values[i]
		}
	default:
		return nil
	}
	return factorMap
}

func UpdateUserFactorByMap(factorMap map[string]int, userId int) {
	o := orm.NewOrm()
	var b bytes.Buffer
	b.WriteString("update user_factor set ")
	for factor,value := range factorMap {
		b.WriteString(factor)
		b.WriteString(" = ")
		b.WriteString(factor)
		b.WriteString(" + (")
		b.WriteString(strconv.Itoa(value))
		b.WriteString("),")
	}
	b.WriteString(" user_id = ")
	b.WriteString(strconv.Itoa(userId))
	b.WriteString(" where ")
	b.WriteString(" user_id = ")
	b.WriteString(strconv.Itoa(userId))
	o.Raw(b.String()).Exec()
}