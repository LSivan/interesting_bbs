<div class="row">
  <div class="col-md-9">
    {{if .CurrentUserInfo}}
    <div class="panel panel-default">
      <div class="panel-body">
        <div class="media">
          <div class="media-left">
            <img src="{{.CurrentUserInfo.Avatar}}" class="avatar-lg" alt="{{.CurrentUserInfo.Username}}">
          </div>
          <div class="media-body">
            <h3 style="margin-top: 0">{{.CurrentUserInfo.Username}}</h3>
            {{if .CurrentUserInfo.Signature}}<p><i class="gray">{{.CurrentUserInfo.Signature}}</i></p>{{end}}
            {{if .CurrentUserInfo.Url}}
            <div>主页: <a href="{{.CurrentUserInfo.Url}}" target="_blank">{{.CurrentUserInfo.Url}}</a></div>
            {{end}}
            <div>入驻时间: {{.CurrentUserInfo.InTime | timeago}}</div>
          </div>
        </div>
      </div>
    </div>
    <div class="panel panel-default">
      <div class="panel-heading">{{.CurrentUserInfo.Username}}创建的话题</div>
      <div class="panel-body">
        {{range .Topics}}
        <div class="media">
          <div class="media-body">
            <div class="title">
              <a href="/topic/{{.Id}}">{{.Title}}</a>
            </div>
            <p>
              <a href="/?tab={{.Section.Id}}">{{.Section.Name}}</a>
              <span>•</span>
              <span><a href="/user/{{.User.Username}}">{{.User.Username}}</a></span>
              <span class="hidden-sm hidden-xs">•</span>
              <span class="hidden-sm hidden-xs">{{.ReplyCount}}个回复</span>
              <span class="hidden-sm hidden-xs">•</span>
              <span class="hidden-sm hidden-xs">{{.View}}次浏览</span>
              <span>•</span>
              <span>{{.InTime | timeago}}</span>
              {{if .LastReplyUser}}
                <span>•</span>
                <span>最后回复来自 <a href="/user/{{.LastReplyUser.Username}}">{{.LastReplyUser.Username}}</a></span>
              {{end}}
            </p>
          </div>
        </div>
        <div class="divide mar-top-5"></div>
        {{end}}
      </div>
      {{ if ge (.Topics|len) 7 }}
      <div class="panel-footer">
        <a href="/user/{{.CurrentUserInfo.Username}}/topics">{{.CurrentUserInfo.Username}}更多话题&gt;&gt;</a>
      </div>
      {{else if eq (.Topics|len) 0 }}
        <div class="panel-footer">
          <a href="javaScript:void(0);">暂无话题～(ˇˍˇ)～</a>
        </div>
      {{end}}
    </div>
    <div class="panel panel-default">
      <div class="panel-heading">{{.CurrentUserInfo.Username}}回复的话题</div>
      <table class="table table-striped">
        <tbody>
        {{range .Replies}}
        <tr>
          <td>
            {{.InTime | timeago}}
            回复了
            <a href="/user/{{.User.Username}}">{{.User.Username}}</a>
            创建的话题 › <a href="/topic/{{.Topic.Id}}">{{.Topic.Title}}</a>
          </td>
        </tr>
        <tr>
          <td><p>{{str2html (.Content | markdown)}}</p></td>
        </tr>
        {{end}}
        </tbody>
      </table>
      {{ if ge (.Replies|len) 7 }}
      <div class="panel-footer">
        <a href="/user/{{.CurrentUserInfo.Username}}/replies">{{.CurrentUserInfo.Username}}更多回复&gt;&gt;</a>
      </div>
      {{else if eq (.Replies|len) 0 }}
          <div class="panel-body"></div>
          <div class="panel-footer">
            <a href="javaScript:void(0);">暂无回复～(ˇˍˇ)～</a>
          </div>
      {{ end}}
    </div>
    <div class="panel panel-default">
          <div class="panel-heading">{{.CurrentUserInfo.Username}}收藏的话题</div>
          <div class="panel-body">
            {{range .Collects}}
            <div class="media">
                      <div class="media-body">
                        <div class="title">
                          <a href="/topic/{{.Id}}">{{.Title}}</a>
                        </div>
                        <p>
                          <a href="/?tab={{.Section.Id}}">{{.Section.Name}}</a>
                          <span>•</span>
                          <span><a href="/user/{{.User.Username}}">{{.User.Username}}</a></span>
                          <span class="hidden-sm hidden-xs">•</span>
                          <span class="hidden-sm hidden-xs">{{.ReplyCount}}个回复</span>
                          <span class="hidden-sm hidden-xs">•</span>
                          <span class="hidden-sm hidden-xs">{{.View}}次浏览</span>
                          <span>•</span>
                          <span>{{.InTime | timeago}}</span>
                          {{if .LastReplyUser}}
                            <span>•</span>
                            <span>最后回复来自 <a href="/user/{{.LastReplyUser.Username}}">{{.LastReplyUser.Username}}</a></span>
                          {{end}}
                        </p>
                      </div>
                    </div>
                    <div class="divide mar-top-5"></div>
            {{end}}
          </div>
          {{ if ge (.Collects|len) 7 }}
            <div class="panel-footer">
              <a href="/user/{{.CurrentUserInfo.Username}}/replies">{{.CurrentUserInfo.Username}}更多的收藏&gt;&gt;</a>
            </div>
          {{else if eq (.Collects|len) 0 }}
            <div class="panel-footer">
              <a href="javaScript:void(0);">暂无收藏～(ˇˍˇ)～</a>
            </div>
          {{end}}
        </div>
    {{else}}
    <div class="panel panel-default">
      <div class="panel-body">用户不存在</div>
    </div>
    {{end}}
    {{ if eq .UserInfo.Id .CurrentUserInfo.Id}}
    <div class="panel panel-default">
      <div class="panel-heading">{{.CurrentUserInfo.Username}}的黑名单话题</div>
      <div class="panel-body">
        {{range .Blacks}}
        <div class="media">
                  <div class="media-body">
                    <div class="title">
                      <a href="/topic/{{.Id}}">{{.Title}}</a>
                    </div>
                    <p>
                      <a href="/?tab={{.Section.Id}}">{{.Section.Name}}</a>
                      <span>•</span>
                      <span><a href="/user/{{.User.Username}}">{{.User.Username}}</a></span>
                      <span class="hidden-sm hidden-xs">•</span>
                      <span class="hidden-sm hidden-xs">{{.ReplyCount}}个回复</span>
                      <span class="hidden-sm hidden-xs">•</span>
                      <span class="hidden-sm hidden-xs">{{.View}}次浏览</span>
                      <span>•</span>
                      <span>{{.InTime | timeago}}</span>
                      {{if .LastReplyUser}}
                        <span>•</span>
                        <span>最后回复来自 <a href="/user/{{.LastReplyUser.Username}}">{{.LastReplyUser.Username}}</a></span>
                      {{end}}
                    </p>
                  </div>
                </div>
                <div class="divide mar-top-5"></div>
        {{end}}
      </div>
      {{ if ge (.Blacks|len) 7 }}
        <div class="panel-footer">
          <a href="/user/{{.CurrentUserInfo.Username}}/replies">{{.CurrentUserInfo.Username}}更多的黑名单话题&gt;&gt;</a>
        </div>
      {{else if eq (.Blacks|len) 0 }}
        <div class="panel-footer">
          <a href="javaScript:void(0);">暂无黑名单   ～(ˇˍˇ)～</a>
        </div>
      {{end}}
    </div>
    {{ end }}
  </div>
  <div class="col-md-3 hidden-sm hidden-xs">

  </div>
</div>