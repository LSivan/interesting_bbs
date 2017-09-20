<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>{{.PageTitle}}</title>
  <link rel="stylesheet" href="/static/bootstrap/css/bootstrap.min.css"/>
  <link rel="stylesheet" href="/static/css/pybbs.css">
  <script src="//cdn.bootcss.com/jquery/2.2.2/jquery.min.js"></script>
  <script src="//cdn.bootcss.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>
</head>
<body>
<div class="wrapper">
  <nav class="navbar navbar-inverse">
    <div class="container">
      <div class="navbar-header">
        <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar"
                aria-expanded="false" aria-controls="navbar">
          <span class="sr-only">Toggle navigation</span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
        </button>
        <a class="navbar-brand" style="color:#fff;" href="/">海大开发者社区</a>
      </div>
      <div id="navbar" class="navbar-collapse collapse header-navbar">
        <ul class="nav navbar-nav">
            <li {{ if .IsIndex}} class="active" {{ end }}>
                 <a href="/">首页</a>
            </li>
            <li {{ if .IsFavorite}} class="active" {{ end }}>
                 <a href="/favorite">猜你喜欢</a>
            </li>
            <li {{ if .favorite}} class="active" {{ end }}>
                 <a href="#">七天热门</a>
            </li>
            <li class="dropdown">
                 <a href="#" class="dropdown-toggle" data-toggle="dropdown">资源<strong class="caret"></strong></a>
                <ul class="dropdown-menu">
                    <li>
                         <a href="#">Java</a>
                    </li>
                    <li>
                         <a href="http://docscn.studygolang.com/" target="_blank">Golang入门</a>
                    </li>
                    <li>
                         <a href="#">Javascript</a>
                    </li>
                    <li class="divider">
                    </li>
                    <li>
                         <a href="#">Separated link</a>
                    </li>
                </ul>
            </li>
      	</ul>
        <ul class="nav navbar-nav navbar-right">
          <li>
            <a href="/about">关于</a>
          </li>
          {{if .IsLogin}}
          <li>
            <a href="/user/{{.UserInfo.Username}}">
              {{.UserInfo.Username}}
            </a>
          </li>
          <li>
            <a href="javascript:;" class="dropdown-toggle" data-toggle="dropdown"
               data-hover="dropdown">
              设置
              <span class="caret"></span>
            </a>
            <span class="dropdown-arrow"></span>
            <ul class="dropdown-menu">
              <li><a href="/user/setting">个人资料</a></li>
              <li><a href="/user/setting">我的收藏</a></li>
              {{if haspermission .UserInfo.Id "user:list"}}<li><a href="/user/list">用户管理</a></li>{{end}}
              {{if haspermission .UserInfo.Id "role:list"}}<li><a href="/role/list">角色管理</a></li>{{end}}
              {{if haspermission .UserInfo.Id "permission:list"}}<li><a href="/permission/list">权限管理</a></li>{{end}}
              <li><a href="/logout">退出</a></li>
            </ul>
          </li>
          {{else}}
          <li><a href="/login">登录</a></li>
          <li><a href="/register">注册</a></li>
          {{end}}
        </ul>
      </div>
    </div>
  </nav>
  <div class="container">
    {{.LayoutContent}}
  </div>
</div>
<div class="container" >
  <br>
  <div class="text-center">
    ©2018 Powered by <a href="javascript:void(0);" target="_blank">海大开发者社区</a>
  </div>
  <br>
</div>
</body>
</html>