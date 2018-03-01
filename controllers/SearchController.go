package controllers

import (
	"fmt"
	"git.oschina.net/gdou-geek-bbs/engine"
	"git.oschina.net/gdou-geek-bbs/filters"
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"log"
	"strconv"
)

type SearchController struct {
	beego.Controller
}

var searcher = engine.Searcher{}

func (c *SearchController) Search() {
	q := c.GetString("q")
	if q == "" {
		c.Redirect("/", 302)
	}
	p, _ := strconv.Atoi(c.Ctx.Input.Query("p"))
	if p == 0 {
		p = 1
	}
	result, err := searcher.Search(q, p)
	if err != nil {
		log.Println("err : ", err)
	}

	c.Data["PageTitle"] = "\"" + q + "\"的搜索结果"
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	c.Data["q"] = q
	// 做一次反查
	if result.Hits.Len() != 0 {
		topicIDS := make([]int, 0, len(result.Hits))
		for _, v := range result.Hits {
			topicIDS = append(topicIDS, utils.MustInt(v.ID))
		}
		topics := models.FindTopicByIDS(topicIDS)
		list := make([]interface{}, 0, len(topics))
		for _, v := range result.Hits {
			for _, topic := range topics {

				if utils.MustInt(v.ID) == topic.Id {
					topic.Content = ""
					fmt.Println("v.Fragments", v.Fragments)
					for _, fragments := range v.Fragments {
						for _, fragment := range fragments {
							topic.Content += fmt.Sprintf("...%s...", fragment)
						}

					}
					list = append(list, topic)
				}
			}
		}

		c.Data["Page"] = utils.PageUtil(
			int(result.Total),
			p,
			beego.AppConfig.DefaultInt("page.size", 8),
			list,
		)
	}
	c.Layout = "layout/layout.tpl"
	c.TplName = "search.tpl"
}
