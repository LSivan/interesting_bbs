package utils

import (
	"github.com/astaxie/beego"
)

func LogError(action string, err error) {
	if err != nil {
		beego.BeeLogger.Error("%s失败，err：%v\n", action,err)
	}
	beego.BeeLogger.Debug("%s成功\n", action)
}