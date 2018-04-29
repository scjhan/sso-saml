<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />   
 	<title>TbInc</title>
 	
 	<style type="text/css">
 		body,html{
 			width: 100%;
 			height: 100%;
 			padding: 0;
 			margin: 0;
 			min-width: 800px;
 			min-height: 600px;
 		}
 		.headerTop{
 			height: 45px;
 			background-color: #333333;
 			color: white;
 			line-height: 40px;
 		}
 		.content{
 			background-color: #fafafa;
 			background-image: url("static/img/welcome.gif");
 			background-position: center;
			background-repeat: no-repeat;
			background-attachment: fixed;
			position: relative;;
			bottom: 0px;
			margin: 0px;
 			text-align: center;
 		}
 		@supports (width:calc()) or (width:calc(100% - 230px)){
 			.content{
 				height: calc(100% - 230px);
 			}
 		}
 		.fltR{
 			float: right;
 		}
 		.btn{
 			font-size: 16px;
 			color: white;
 			cursor: pointer;
 			opacity: 0.8;
 			height: 28px;
 			text-align: center;
 			vertical-align: middle;
 			margin: 1px 5px;
 			display: inline-block;
 			border-radius: 3px;
 			line-height: 28px;

 			transition-property: opacity;
 			transition-duration: 0.3s;

 			-webkit-transition-property: opacity;
 			-webkit-transition-duration: 0.3s;

 			-ms-transition-property: opacity;
 			-ms-transition-duration: 0.3s;

 			-moz-transition-property: opacity;
 			-moz-transition-duration: 0.3s;

 			-o-transition-property: opacity;
 			-o-transition-duration: 0.3s;
 		}
 		.hello-text{

 		}
 		.btn:hover{
 			color: #b9b9b9
 		}
		footer {
			line-height: 0.0;
			text-align: center;
			padding: 50px 0;
			color: #999;
		}
		.content footer {
			position: absolute;
			left: 0px;
			right: 0px;
			bottom: 0px;
			margin: 0 auto;
		}
 	</style>

 	<script src="static/js/jquery-3.3.1.min.js"></script>

	<script type="text/javascript">
		function onSignOutBtnClick() {
			$.ajax({
				type : "GET",
				url : "/logout",
				success : function(msg) {
					console.log(msg);
					window.location.href="/index";
				}
			});
		}
	</script>
</head>
<body>
	<div class="headerTop">
		<div style="float: right;height: 40px;line-height: 40px;margin-right: 20px;">
			<a class="btn">Hello {{.UserName}}</a>
			<a class="btn" onclick="onSignOutBtnClick()">Sign out</a>
		</div>
	</div>
	<div style="background-color: #fafafa;">
		<div style="height:80px;"></div>
		<div style="text-align: center;font-size: 50px;font-weight:bold;">TB Inc.</div>
		<div style="text-align: center;font-size: 24px;">Say hello to the visitor.</div>
	</div>
	<div class="content">
		<footer>
			<div class="author">
				Copyright © ChenJunhan 2018
			</div>
		</footer>
	</div>
</body>
</html>