package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type UserTopicList struct {
	Id         int       `orm:"pk;auto"`
	InTime     time.Time `orm:"auto_now_add;type(datetime)"`
	User       *User     `orm:"rel(fk)"`
	Topic      *Topic    `orm:"rel(fk)"`
	ActionType int       `orm:"default(1)"` // 1为收藏，2为拉黑
}

func SaveUserTopic(userTopicList *UserTopicList) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(userTopicList)
	return id
}
func FindTopicByUserAndTopicAndActionType(user *User, topic *Topic, actionType int) (bool, *UserTopicList) {
	o := orm.NewOrm()
	var userTopicList UserTopicList
	err := o.QueryTable(userTopicList).RelatedSel().Filter("User", user).Filter("Topic", topic).Filter("ActionType", actionType).One(&userTopicList)

	return err != orm.ErrNoRows, &userTopicList
}

func FindAllUserTopicList() ( *[]UserTopicList) {
	o := orm.NewOrm()
	var userTopicList []UserTopicList
	o.QueryTable(UserTopicList{}).RelatedSel().All(&userTopicList)
	return  &userTopicList
}

func DeleteUserTopic(userTopicList *UserTopicList) {
	o := orm.NewOrm()
	o.Delete(userTopicList)
}
