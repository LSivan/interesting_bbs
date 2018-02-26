package controllers

import (
	"github.com/astaxie/beego"
	"git.oschina.net/gdou-geek-bbs/engine"
	"log"
	"fmt"
)

type SearchController struct {
	beego.Controller
}

var searcher engine.Searcher = engine.Searcher{}

func (c *SearchController) Search() {
	q := c.GetString("q")
	result, err := searcher.Search(q)
	if err != nil {
		log.Println("err : ",err)
	}
	fmt.Println(result)
	//for _, v := range result.Hits {
	//
	//}
	c.Data["json"] = &result
	c.ServeJSON()
}
