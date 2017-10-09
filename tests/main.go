package main

import (
	"github.com/robfig/cron"
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego/logs"
)

func main() {
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		list := models.FindUserFactorChangeSum()
		for _,v := range list {
			logs.Debug("v  : ",v)
		}
	})
	c.Start()

	select{}
}
