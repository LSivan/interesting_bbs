package utils

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/russross/blackfriday"
	"github.com/xeonx/timeago"
	"time"
)

func FormatTime(time time.Time) string {
	return timeago.English.Format(time)
}

func Markdown(content string) string {
	return string(blackfriday.MarkdownCommon([]byte(utils.NoHtml(content))))
}

func HasPermission(userId int, name string) bool {
	return models.FindPermissionByUserIdAndPermissionName(userId, name)
}

func IsCollect(user models.User, topic models.Topic) bool {
	b, _ := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 1)// 1为收藏
	return b
}
func IsBlack(user models.User, topic models.Topic) bool {
	b, _ := models.FindTopicByUserAndTopicAndActionType(&user, &topic, 2)// 2为拉黑
	return b
}
func init() {
	beego.AddFuncMap("timeago", FormatTime)
	beego.AddFuncMap("markdown", Markdown)
	beego.AddFuncMap("haspermission", HasPermission)
	beego.AddFuncMap("iscollect", IsCollect)
	beego.AddFuncMap("isblack", IsBlack)
}
