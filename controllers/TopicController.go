package controllers

import (
	"git.oschina.net/gdou-geek-bbs/engine"
	"git.oschina.net/gdou-geek-bbs/filters"
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"strconv"
	"git.oschina.net/gdou-geek-bbs/recommend"
	"git.oschina.net/gdou-geek-bbs/common"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Create() {
	beego.ReadFromRequest(&c.Controller)
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	c.Data["PageTitle"] = "发布话题"
	c.Data["Sections"] = models.FindAllSection()
	c.Layout = "layout/layout.tpl"
	c.TplName = "topic/create.tpl"
}

func (c *TopicController) Save() {
	flash := beego.NewFlash()
	title, content, sid := c.Input().Get("title"), c.Input().Get("content"), c.Input().Get("sid")
	if len(title) == 0 || len(title) > 120 {
		flash.Error("话题标题不能为空且不能超过120个字符")
		flash.Store(&c.Controller)
		c.Redirect("/topic/create", 302)
	} else if len(sid) == 0 {
		flash.Error("请选择话题版块")
		flash.Store(&c.Controller)
		c.Redirect("/topic/create", 302)
	} else {
		s, _ := strconv.Atoi(sid)
		section := models.Section{Id: s}
		_, user := filters.IsLogin(c.Ctx)
		topic := models.Topic{Title: title, Content: content, Section: &section, User: &user}
		id := models.SaveTopic(&topic)
		topic.Id = int(id)
		engine.Indexer.InsertChan <- &topic
		//topicFactor := models.TopicFactor{}.New(section.Id)
		//topicFactor.Topic = &topic
		//models.SaveTopicFactor(topicFactor)
		go func(t *models.Topic) {
			feature := recommend.GetTopicFeature(t)
			err := common.Redis.HSet("topic-feature",strconv.Itoa(t.Id),feature).Err()
			if err != nil {
				// 一般分词没有结果,导致feature.tokens为[NaN,NaN]而不能存进redis中,从而过滤无关的文章,TODO 可以用来更新关键字列表
				beego.BeeLogger.Info("recommend.ID:%v, err :%v",feature.ID, err)
			}
			recommend.ChangeUserFeature(user.Id,4,t)
		}(&topic)

		c.Redirect("/topic/"+strconv.FormatInt(id, 10), 302)
	}
}

func (c *TopicController) Detail() {
	id := c.Ctx.Input.Param(":id")
	tid, _ := strconv.Atoi(id)
	if tid > 0 {
		c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
		topic := models.FindTopicById(tid)
		models.IncrView(&topic) //查看+1
		c.Data["PageTitle"] = topic.Title
		c.Data["Topic"] = topic
		c.Data["Replies"] = models.FindReplyByTopic(&topic)
		c.Layout = "layout/layout.tpl"
		c.TplName = "topic/detail.tpl"
	} else {
		c.Ctx.WriteString("话题不存在")
	}
}

func (c *TopicController) Edit() {
	beego.ReadFromRequest(&c.Controller)
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if id > 0 {
		topic := models.FindTopicById(id)
		c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
		c.Data["PageTitle"] = "编辑话题"
		c.Data["Sections"] = models.FindAllSection()
		c.Data["Topic"] = topic
		c.Layout = "layout/layout.tpl"
		c.TplName = "topic/edit.tpl"
	} else {
		c.Ctx.WriteString("话题不存在")
	}
}

func (c *TopicController) Update() {
	flash := beego.NewFlash()
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	title, content, sid := c.Input().Get("title"), c.Input().Get("content"), c.Input().Get("sid")
	if len(title) == 0 || len(title) > 120 {
		flash.Error("话题标题不能为空且不能超过120个字符")
		flash.Store(&c.Controller)
		c.Redirect("/topic/edit/"+strconv.Itoa(id), 302)
	} else if len(sid) == 0 {
		flash.Error("请选择话题版块")
		flash.Store(&c.Controller)
		c.Redirect("/topic/edit/"+strconv.Itoa(id), 302)
	} else {
		s, _ := strconv.Atoi(sid)
		section := models.Section{Id: s}
		topic := models.FindTopicById(id)
		topic.Title = title
		topic.Content = content
		topic.Section = &section
		models.UpdateTopic(&topic)
		engine.Indexer.UpdateChan <- &topic
		c.Redirect("/topic/"+strconv.Itoa(id), 302)
	}
}

func (c *TopicController) Delete() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if id > 0 {
		topic := models.FindTopicById(id)
		models.DeleteTopic(&topic)
		models.DeleteReplyByTopic(&topic)
		engine.Indexer.DeleteChan <- &topic
		c.Redirect("/", 302)
	} else {
		c.Ctx.WriteString("话题不存在")
	}
}

func (c *TopicController) Collect() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	result := utils.Result{Code: 500, Description: "话题不存在"}
	c.Data["json"] = &result
	if id > 0 {

		topic := models.FindTopicById(id)
		_, user := filters.IsLogin(c.Controller.Ctx)
		b, _ := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 1)
		if !b { // 确保不存在
			models.SaveUserTopic(&models.UserTopicList{Topic: &topic, User: &user, ActionType: 1}) // 1为收藏
			models.IncrCollectCount(&topic)
			result := utils.Result{Code: 200, Description: "成功"}
			c.Data["json"] = &result
		}
	}
	c.ServeJSON()
}

func (c *TopicController) CancelCollect() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	result := utils.Result{Code: 500, Description: "话题不存在"}
	c.Data["json"] = &result
	if id > 0 {
		topic := models.FindTopicById(id)
		_, user := filters.IsLogin(c.Controller.Ctx)
		b, userTopicList := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 1)
		if b { // 确保存在
			models.DeleteUserTopic(userTopicList)
			models.ReduceCollectCount(&topic)
			result := utils.Result{Code: 200, Description: "成功"}
			c.Data["json"] = &result
		}
	}
	c.ServeJSON()
}

func (c *TopicController) Black() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	result := utils.Result{Code: 500, Description: "话题不存在"}
	c.Data["json"] = &result
	if id > 0 {

		topic := models.FindTopicById(id)
		_, user := filters.IsLogin(c.Controller.Ctx)
		b, _ := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 2)
		if !b { // 确保不存在
			models.SaveUserTopic(&models.UserTopicList{Topic: &topic, User: &user, ActionType: 2}) // 1为收藏
			models.IncrCollectCount(&topic)
			result := utils.Result{Code: 200, Description: "成功"}
			c.Data["json"] = &result
		}
	}
	c.ServeJSON()
}

func (c *TopicController) CancelBlack() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	result := utils.Result{Code: 500, Description: "话题不存在"}
	c.Data["json"] = &result
	if id > 0 {
		topic := models.FindTopicById(id)
		_, user := filters.IsLogin(c.Controller.Ctx)
		b, userTopicList := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 2)
		if b { // 确保存在
			models.DeleteUserTopic(userTopicList)
			models.ReduceCollectCount(&topic)
			result := utils.Result{Code: 200, Description: "成功"}
			c.Data["json"] = &result
		}
	}
	c.ServeJSON()
}
