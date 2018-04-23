package filters

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"regexp"
	"git.oschina.net/gdou-geek-bbs/recommend"
	"strconv"
)

func IsLogin(ctx *context.Context) (bool, models.User) {
	token, flag := ctx.GetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"))
	var user models.User
	if flag {
		flag, user = models.FindUserByToken(token)
	}
	return flag, user
}

var HasPermission = func(ctx *context.Context) {
	ok, user := IsLogin(ctx)
	if !ok {
		ctx.Redirect(302, "/login")
	} else {
		permissions := models.FindPermissionByUser(user.Id)
		url := ctx.Request.RequestURI
		beego.Debug("url: ", url)
		var flag = false
		for _, v := range permissions {
			if a, _ := regexp.MatchString(v.Url, url); a {
				flag = true
				break
			}
		}
		if !flag {
			ctx.WriteString("你没有权限访问这个页面")
		}
	}
}

var withLoginCheck = func(fn func(), ctx *context.Context) {
	ok, _ := IsLogin(ctx)

	if ok { // 用户已登录才进行
		reg := regexp.MustCompile(`/topic/([0-9]+)`)
		if reg.MatchString(ctx.Input.URI()) { // 如果是查看话题详情
			flag := ctx.Input.Query("flag")
			if flag == "true" { // 则看是不是在猜你喜欢页面过来的请求
				fn()
			}
		} else { // 不匹配，直接执行
			fn()
		}
	}
}

var DetailsChangeFactor = func(ctx *context.Context) { // 用户查看话题详情时执行
	withLoginCheck(
		func() {
			ChangeFactor(0.5, ctx)
		}, ctx)
}
var BlackChangeFactor = func(ctx *context.Context) { // 用户拉黑时执行
	withLoginCheck(
		func() {
			ChangeFactor(-2, ctx)
		}, ctx)
}
var CancelBlackChangeFactor = func(ctx *context.Context) { // 用户取消拉黑时执行
	withLoginCheck(
		func() {
			ChangeFactor(2, ctx)
		}, ctx)
}
var CollectChangeFactor = func(ctx *context.Context) { // 用户收藏时执行
	withLoginCheck(
		func() {
			ChangeFactor(3, ctx)
		}, ctx)
}
var ReplyChangeFactor = func(ctx *context.Context) { // 用户查看话题详情时执行
	withLoginCheck(
		func() {
			ChangeFactor(2, ctx)
		}, ctx)
}
var FilterNoCache = func(ctx *context.Context) {
	ctx.Output.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Output.Header("Pragma", "no-cache")
	ctx.Output.Header("Expires", "0")
}

var FilterUser = func(ctx *context.Context) {
	ok, _ := IsLogin(ctx)
	if !ok {
		ctx.Redirect(302, "/login")
	}
}

func ChangeFactor(changeValue float64, ctx *context.Context) {
	_, user := IsLogin(ctx)
	id,err := strconv.Atoi(ctx.Input.Param(":id"))
	if err != nil {
		return
	}
	topic := models.FindTopicById(id)
	if &topic == nil {
		return
	}
	recommend.ChangeUserFeature(user.Id,changeValue,&topic)

}
