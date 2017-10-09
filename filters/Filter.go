package filters

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"regexp"
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

// TODO 一些不变的数据使用redis，比如话题等
// TODO 查看他人资料时还可以看到它的收藏
// TODO 更多话题/回复/收藏
var DetailsChangeFactor = func(ctx *context.Context) { // 用户查看话题详情时执行
	withLoginCheck(
		func() {
			ChangeFactor(1, ctx)
		}, ctx)

}

var BlackChangeFactor = func(ctx *context.Context) { // 用户拉黑时执行
	withLoginCheck(
		func() {
			ChangeFactor(-2, ctx)
		}, ctx)
}
var CollectChangeFactor = func(ctx *context.Context) { // 用户收藏时执行
	withLoginCheck(
		func() {
			ChangeFactor(4, ctx)
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

func ChangeFactor(changeValue int, ctx *context.Context) {
	_, user := IsLogin(ctx)
	id := ctx.Input.Param(":id")
	/******** 得到用户以及话题的特征因子和无关因子 *********/
	userFactor, userFeatureFactorMap, userUnusedFactorMap := getUserFactor(user)
	topicFactor, topicFeatureFactorMap, topicUnusedFactorMap := getTopicFactor(models.FindTopicById(utils.MustInt(id)))
	/*
		用户中与ThingFeatureFactor相同的因子，全部加上因子的变化度；
		用户中与ThingUnusedFactor相同的因子，全部减去因子的变化度。
		图书中与UserFeatureFactor相同的因子，全部加上因子的变化度；
		图书中与UserUnusedFactor相同的因子，全部减去因子的变化度。
	*/
	topicFactorChangeMap := make(map[string]int)
	for factor := range userFeatureFactorMap {
		topicFactorChangeMap[factor] = changeValue
	}
	for factor := range userUnusedFactorMap {
		topicFactorChangeMap[factor] = -1 * changeValue
	}
	models.SaveTmpTopicFactorByMap(topicFactorChangeMap, topicFactor.Id)
	userFactorChangeMap := make(map[string]int)
	for factor := range topicFeatureFactorMap {
		userFactorChangeMap[factor] = changeValue
	}
	for factor := range topicUnusedFactorMap {
		userFactorChangeMap[factor] = -1 * changeValue
	}
	models.SaveTmpUserFactorByMap(userFactorChangeMap, userFactor.Id)
}

var getUserFactor = func(user models.User) (models.UserFactor, map[string]int, map[string]int) {
	return models.FindFactorByUser(&user), models.FindFactorByUser(&user).GetTopFactorByType(0), models.FindFactorByUser(&user).GetTopFactorByType(1)
}
var getTopicFactor = func(topic models.Topic) (models.TopicFactor, map[string]int, map[string]int) {
	return models.FindFactorByTopic(&topic), models.FindFactorByTopic(&topic).GetTopFactorByType(0), models.FindFactorByTopic(&topic).GetTopFactorByType(1)
}
