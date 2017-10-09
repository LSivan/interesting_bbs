package models

import (
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Topic struct {
	Id            int       `orm:"pk;auto"`
	Title         string    `orm:"unique"`
	Content       string    `orm:"type(text);null"`
	InTime        time.Time `orm:"auto_now_add;type(datetime)"`
	User          *User     `orm:"rel(fk)"`
	Section       *Section  `orm:"rel(fk)"`
	View          int       `orm:"default(0)"`
	ReplyCount    int       `orm:"default(0)"`
	CollectCount  int       `orm:"default(0)"`
	LastReplyUser *User     `orm:"rel(fk);null"`
	LastReplyTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func SaveTopic(topic *Topic) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(topic)
	return id
}

func FindTopicById(id int) Topic {
	o := orm.NewOrm()
	var topic Topic
	o.QueryTable(topic).RelatedSel().Filter("Id", id).One(&topic)
	return topic
}

func PageTopic(p int, size int, section *Section) utils.Page {
	o := orm.NewOrm()
	var topic Topic
	var list []Topic
	qs := o.QueryTable(topic)
	if section.Id > 0 {
		qs = qs.Filter("Section", section)
	}
	count, _ := qs.Limit(-1).Count()
	qs.RelatedSel().OrderBy("-InTime").Limit(size).Offset((p - 1) * size).All(&list)
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))
	return utils.PageUtil(c, p, size, list)
}

func FavoritePageTopic(p int, size int, user *User) utils.Page {
	o := orm.NewOrm()
	var topic Topic
	qs := o.QueryTable(topic)
	count, _ := qs.Limit(-1).Count()
	s := "select  topic.id,topic.title,topic.content,topic.in_time,topic.user_id,topic.section_id,topic.`view`,topic.reply_count,topic.last_reply_user_id,topic.last_reply_time," +
		"(uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor + " +
		"uf.share_factor * tf.share_factor) factor " +
		" FROM user_factor uf,topic_factor tf ,topic " +
		" where uf.user_id = ? and " +
		" tf.topic_id = topic.id and " +
		" topic.id not in(select topic_id from user_topic_list where user_id = ? and action_type = 2) "+ // 拉黑的文章不展示在"猜你喜欢"中
		" order by factor desc,topic.id desc " +
		" limit ? offset ?  "
	var topics []Topic
	_, err := o.Raw(s, user.Id,user.Id, size, (p-1)*size).QueryRows(&topics)
	if err != nil {
		logs.Debug("FavoritePageTopic: ", err)
	}
	for i, t := range topics {
		b, user := FindUserById(t.User.Id)
		if b {
			t.User = &user
		}
		b, section := FindSectionById(t.Section.Id)
		if b {
			t.Section = &section
		}
		b, lastReplyUser := FindUserById(t.User.Id)
		if b {
			t.LastReplyUser = &lastReplyUser
		}
		(topics)[i] = t // 至关重要的一步,TODO 可以写个博客来研究下
	}
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))
	return utils.PageUtil(c, p, size, &topics)
}

func IncrView(topic *Topic) {
	o := orm.NewOrm()
	topic.View = topic.View + 1
	o.Update(topic, "View")
}

func IncrReplyCount(topic *Topic) {
	o := orm.NewOrm()
	topic.ReplyCount = topic.ReplyCount + 1
	o.Update(topic, "ReplyCount")
}

func ReduceReplyCount(topic *Topic) {
	o := orm.NewOrm()
	topic.ReplyCount = topic.ReplyCount - 1
	o.Update(topic, "ReplyCount")
}

func IncrCollectCount(topic *Topic) {
	o := orm.NewOrm()
	topic.CollectCount = topic.CollectCount + 1
	o.Update(topic, "CollectCount")
}

func ReduceCollectCount(topic *Topic) {
	o := orm.NewOrm()
	topic.CollectCount = topic.CollectCount - 1
	o.Update(topic, "CollectCount")
}

func FindTopicByUser(user *User, limit int) []*Topic {
	o := orm.NewOrm()
	var topic Topic
	var topics []*Topic
	o.QueryTable(topic).RelatedSel().Filter("User", user).OrderBy("-LastReplyTime", "-InTime").Limit(limit).All(&topics)
	return topics
}

func FindCollectTopicByUser(user *User, limit int )[]*Topic{
	o := orm.NewOrm()
	var topic Topic
	var topics []*Topic
	//  utl where utl.topic_id=? and utl.action_type = 1
	o.QueryTable(topic).RelatedSel("UserTopicList").Filter("User", user).OrderBy("-LastReplyTime", "-InTime").Limit(limit).All(&topics)
	return topics
}

func UpdateTopic(topic *Topic) {
	o := orm.NewOrm()
	o.Update(topic)
}

func DeleteTopic(topic *Topic) {
	o := orm.NewOrm()
	o.Delete(topic)
}

func DeleteTopicByUser(user *User) {
	o := orm.NewOrm()
	o.Raw("delete from topic where user_id = ?", user.Id).Exec()
}
