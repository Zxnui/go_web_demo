{{define "header"}}
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/css/bootstrap.min.css">
	<link rel="stylesheet" type="text/css" href="/css/this.css">
	<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
	<script src="/js/jquery-1.11.3.min.js"></script>
	<script type="text/javascript" src="/js/bootstrap.min.js"></script>
{{end}}

{{define "menu_list"}}
<a class="navbar-brand" href="/"><img src="/img/icon.png" class="web-icon"></a>
<div>
	<ul class="nav navbar-nav">
		<li {{if .isGames}}class="active"{{end}}><a href="/game">游戏</a></li>
		<li {{if .isDate}}class="active"{{end}}><a href="/date">数据</a></li>
	</ul>
</div>

<div class="pull-right">
	<ul class="nav navbar-nav">
		{{if .isLogin}}
		<li><a href="/login?exit=true">退出</a></li>
		{{else}}
		<li><a href="/login">登录/注册</a></li>
		{{end}}
	</ul>
</div>
{{end}}