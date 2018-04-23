package models

import (
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"sync"
	"fmt"
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
		" topic.id not in(select topic_id from user_topic_list where user_id = ? and action_type = 2) " + // 拉黑的文章不展示在"猜你喜欢"中
		" order by factor desc,topic.id desc " +
		" limit ? offset ?  "
	var topics []*Topic
	_, err := o.Raw(s, user.Id, user.Id, size, (p-1)*size).QueryRows(&topics)
	if err != nil {
		logs.Debug("FavoritePageTopic: ", err)
	}
	fillTopicFields(&topics)
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
func FindTopicFrom(offset int, limit int) []*Topic {
	o := orm.NewOrm()
	var topic Topic
	var topics []*Topic
	o.QueryTable(topic).RelatedSel().Offset(offset).Limit(limit).All(&topics)
	return topics
}
func CountTopicFromID(id int) int {
	o := orm.NewOrm()
	var topic Topic
	count, err := o.QueryTable(topic).Filter("id__gt", id).Count()
	if err == nil {
		return int(count)
	}
	return 0
}
func FindTopicByIDS(IDS []int) []*Topic {
	o := orm.NewOrm()
	var topics []*Topic
	placeHolder := ""
	for i := 0;i<len(IDS);i++ {
		if i != 0 {
			placeHolder += ","
		}
		placeHolder += "?"
	}
	sql := fmt.Sprintf("SELECT * FROM bbs.topic where id in (%s) ORDER BY field(id,%s)",placeHolder,placeHolder)
	o.Raw(sql,IDS,IDS).QueryRows(&topics)
	//o.QueryTable(Topic{}).RelatedSel().Filter("id__in", IDS).All(&topics)
	fillTopicFields(&topics)
	return topics
}
func FindTopicByUser(user *User, limit int) []*Topic {
	o := orm.NewOrm()
	var topic Topic
	var topics []*Topic
	o.QueryTable(topic).RelatedSel().Filter("User", user).OrderBy("-LastReplyTime", "-InTime").Limit(limit).All(&topics)
	return topics
}

func findTopicByUserAndActionType(user *User, actionType int, limit int, p int) []*Topic {
	o := orm.NewOrm()
	var topics []*Topic
	sql := "select " +
		"topic.`id`, topic.`title`, topic.`content`, topic.`in_time`, topic.`user_id`, " +
		"topic.`section_id`, topic.`view`, topic.`reply_count`, topic.`last_reply_user_id`, topic.`last_reply_time`, topic.`collect_count` " +
		"from topic inner join user_topic_list " +
		"where topic.id = user_topic_list.topic_id and " +
		"user_topic_list.action_type = ? and " +
		"user_topic_list.user_id = ? " +
		"order by topic.in_time desc " +
		"limit ? offset ?"
	o.Raw(sql, actionType, user.Id, limit, limit*(p-1)).QueryRows(&topics)
	fillTopicFields(&topics)
	return topics
}

var fillTopicFields = func(topics *[]*Topic) {
	if topics != nil && len(*topics) > 0 {
		userIDS := make([]int, 0, len(*topics)*2)
		sectionIDS := make([]int, 0, len(*topics))
		for _, t := range *topics {
			userIDS = append(userIDS, t.User.Id)
			if t.LastReplyUser != nil {
				userIDS = append(userIDS, t.LastReplyUser.Id)
			}
			sectionIDS = append(sectionIDS, t.Section.Id)
		}
		users := FindUserByIDS(userIDS)
		sections := FindSectionByIDS(sectionIDS)
		wg := &sync.WaitGroup{}
		for j, topic := range *topics {
			wg.Add(1)
			go func(t *Topic, i int) {
				defer wg.Done()
				for _,user := range users {
					if t.User.Id == user.Id {
						t.User = user
					}
					if t.LastReplyUser != nil {
						if t.LastReplyUser.Id == user.Id {
							t.LastReplyUser = user
						}
					}
				}
				for _, section := range sections {
					if t.Section.Id == section.Id {
						t.Section = section
					}
				}
				(*topics)[i] = t // t只是循环的时候的一个变量，必须要赋值到topics上
			}(topic,j)
		}
		wg.Wait()
	}
}

func FindCollectTopicByUser(user *User, limit int, p int) []*Topic {
	return findTopicByUserAndActionType(user, 1, limit, p)
}
func FindBlackTopicByUser(user *User, limit int, p int) []*Topic {
	return findTopicByUserAndActionType(user, 2, limit, p)
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
