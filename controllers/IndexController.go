package controllers

import (
	"git.oschina.net/gdou-geek-bbs/filters"
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego"
	"github.com/sluu99/uuid"
	"strconv"
	"git.oschina.net/gdou-geek-bbs/common"
	"git.oschina.net/gdou-geek-bbs/recommend"
	"git.oschina.net/gdou-geek-bbs/utils"
)

type IndexController struct {
	beego.Controller
}

//首页
func (c *IndexController) Index() {
	c.Data["PageTitle"] = "首页"
	c.Data["IsIndex"] = true
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	p, _ := strconv.Atoi(c.Ctx.Input.Query("p"))
	if p == 0 {
		p = 1
	}
	size, _ := beego.AppConfig.Int("page.size")
	s, _ := strconv.Atoi(c.Ctx.Input.Query("s"))
	c.Data["S"] = s
	section := models.Section{Id: s}
	c.Data["Page"] = models.PageTopic(p, size, &section)
	c.Data["Sections"] = models.FindAllSection()
	c.Layout = "layout/layout.tpl"
	c.TplName = "index.tpl"
}

//登录页
func (c *IndexController) LoginPage() {
	IsLogin, _ := filters.IsLogin(c.Ctx)
	if IsLogin {
		c.Redirect("/", 302)
	} else {
		beego.ReadFromRequest(&c.Controller)
		u := models.FindPermissionByUser(1)
		beego.Debug(u)
		c.Data["PageTitle"] = "登录"
		c.Layout = "layout/layout.tpl"
		c.TplName = "login.tpl"
	}
}

//验证登录
func (c *IndexController) Login() {
	flash := beego.NewFlash()
	username, password := c.Input().Get("username"), c.Input().Get("password")
	if flag, user := models.Login(username, password); flag {
		c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), user.Token, 30*24*60*60, "/", beego.AppConfig.String("cookie.domain"), false, true)
		c.Redirect("/", 302)
	} else {
		flash.Error("用户名或密码错误")
		flash.Store(&c.Controller)
		c.Redirect("/login", 302)
	}
}

//注册页
func (c *IndexController) RegisterPage() {
	IsLogin, _ := filters.IsLogin(c.Ctx)
	if IsLogin {
		c.Redirect("/", 302)
	} else {
		beego.ReadFromRequest(&c.Controller)
		c.Data["Sections"] = models.FindAllSection()
		c.Data["PageTitle"] = "注册"
		c.Layout = "layout/layout.tpl"
		c.TplName = "register.tpl"
	}
}

//验证注册
func (c *IndexController) Register() {
	flash := beego.NewFlash()
	username, password := c.Input().Get("username"), c.Input().Get("password")
	if len(username) == 0 || len(password) == 0 {
		flash.Error("用户名或密码不能为空")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
	} else if flag, _ := models.FindUserByUserName(username); flag {
		flash.Error("用户名已被注册")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
	} else {
		var token = uuid.Rand().Hex()
		user := models.User{Username: username, Password: password, Avatar: "/static/imgs/avatar.png", Token: token}
		models.SaveUser(&user)
		/** 默认普通用户的角色 **/
		commonUserId, _ := beego.AppConfig.Int("constant.common_user_id")
		models.SaveUserRole(user.Id, commonUserId)
		/** 根据用户选的感兴趣的模块赋用户因子 **/
		//userFactor := models.UserFactor{}.New(sections)
		//userFactor.User = &user
		//models.SaveUserFactor(userFactor)
		/** 赋予用户默认的特征值 **/
		// 放回到redis中
		common.Redis.HSet("user-feature", strconv.Itoa(user.Id), recommend.DefaultUserFeature)
		c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), token, 30*24*60*60, "/", beego.AppConfig.String("cookie.domain"), false, true)
		c.Redirect("/", 302)
	}
	c.Redirect("/", 302)
}

//登出
func (c *IndexController) Logout() {
	c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), "", -1, "/", beego.AppConfig.String("cookie.domain"), false, true)
	c.Redirect("/", 302)
}

//关于
func (c *IndexController) About() {
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	c.Data["PageTitle"] = "关于"
	c.Layout = "layout/layout.tpl"
	c.TplName = "about.tpl"
}

// 猜你喜欢
func (c *IndexController) Favorite() {

	c.Data["PageTitle"] = "猜你喜欢"
	isLogin, user := filters.IsLogin(c.Controller.Ctx)
	c.Data["IsLogin"], c.Data["UserInfo"] = isLogin,user
		c.Data["IsFavorite"] = true
	p, _ := strconv.Atoi(c.Ctx.Input.Query("p"))
	if p == 0 {
		p = 1
	}
	size, _ := beego.AppConfig.Int("page.size")
	s, _ := strconv.Atoi(c.Ctx.Input.Query("s"))
	c.Data["S"] = s
	sc := common.Redis.HGet("user-favorite",strconv.Itoa(user.Id))
	if err := sc.Err(); err != nil {
		beego.BeeLogger.Error("get %s favorite list err : %v", user.Id, err)
	}
	ft := &recommend.FavoriteTopic{}
	if err := ft.UnmarshalBinary([]byte(sc.Val())); err != nil {
		beego.BeeLogger.Error("get %s favorite list err : %v", user.Id, err)
	}
	var topics []*models.Topic
	if (p) * size > len(ft.TopicIDS) && (p-1)*size > len(ft.TopicIDS) {
		topics = make([]*models.Topic,0)
	}else if (p) * size > len(ft.TopicIDS) && (p-1)*size < len(ft.TopicIDS) {
		topics = models.FindTopicByIDS(ft.TopicIDS[(p-1)*size:])
	}else{
		topics = models.FindTopicByIDS(ft.TopicIDS[(p-1)*size:(p)*size])
	}
	c.Data["Page"] = utils.PageUtil(len(ft.TopicIDS), p, size, &topics)
	//c.Data["Sections"] = models.FindAllSection()
	c.Layout = "layout/layout.tpl"
	c.TplName = "favorite.tpl"
}
