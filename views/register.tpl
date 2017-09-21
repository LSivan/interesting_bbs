<div class="row">
  <div class="col-md-6">
    <div class="panel panel-default">
      <div class="panel-heading">注册</div>
      <div class="panel-body">
        {{template "components/flash_error.tpl" .}}
        <form action="/register" method="post">
          <div class="form-group">
            <label for="username">用户名</label>
            <input type="text" id="username" name="username" class="form-control" placeholder="用户名">
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input type="password" id="password" name="password" class="form-control" placeholder="密码">
          </div>
          <div class="form-group">
            <label for="title">感兴趣的版块（最多选取6个哦）</label>
              <div class="col-sm-offset-0 col-sm-12">
                <div class="checkbox">
                  <label>
                    {{range $index, $elem := .Sections}}
                      <label>
                      <input type="checkbox" name="sections" value="{{$elem.Id}}" id="title" class="section">{{$elem.Name}}
                      </label>
                    {{end}}
                  </label>
                </div>
              </div>
            </div>
          <input type="submit" class="btn btn-sm btn-default" value="注册"> <a href="/login">去登录</a>
        </form>
      </div>
    </div>
  </div>
</div>

<script>
$('input:checkbox').click(function () {
    var index = 0;
    $('input[name="sections"]').each(function(){
        if ($(this).is(":checked")){
            index = index + 1;
        }
    });
    if (index>6){
        return false;
    }
});

</script>