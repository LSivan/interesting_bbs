package utils

import (
    "time"
    "github.com/xeonx/timeago"
    "github.com/russross/blackfriday"
    "github.com/astaxie/beego"
    "git.oschina.net/gdou-geek-bbs/models"
    "git.oschina.net/gdou-geek-bbs/utils"
)

func FormatTime(time time.Time) string {
    return timeago.Chinese.Format(time)
}

func Markdown(content string) string {
    return string(blackfriday.MarkdownCommon([]byte(utils.NoHtml(content))))
}

func HasPermission(userId int, name string) bool {
    return models.FindPermissionByUserIdAndPermissionName(userId, name)
}

func init() {
    beego.AddFuncMap("timeago", FormatTime)
    beego.AddFuncMap("markdown", Markdown)
    beego.AddFuncMap("haspermission", HasPermission)
}