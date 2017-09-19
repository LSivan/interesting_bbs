package models

import "github.com/astaxie/beego/orm"
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