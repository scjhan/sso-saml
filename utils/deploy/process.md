访问主站 http://tm.com/index
    查询cookie有没有uid，有则表示已经登录，返回显示用户名的index页面
    没有的话返回带有登录按钮的unknown页面，登录链接为http://tm.com/login/?return_to=http://tm.com/index

http://idp.com/check_login?
    return_to=http://tm.com/login?
        token=tokenvalue&
        return_to=http://tm.com/index

[首次登陆]
访问      http://tm.com/index
跳转      http://tm.com/login?return_to=http%3A%2F%2Ftm.com%2Findex
重定向    http://idp.com/check_login?return_to=http%3A%2F%2Ftm.com%2Flogin%3Ftoken%3Dtokenvalue%26return_to%3Dhttp%253A%252F%252Ftm.com%252Findex
返回      http://idp.com/login_page
访问      http://idp.com/login
重定向    http://tm.com/login?token=tokenvalue&return_to=http%3A%2F%2Ftm.com%2Findex
访问      http://idp.com/verify_token?token=tokenvalue
创建session
重定向    http://tm.com/index

[第二次登录]
访问      http://tb.com/index
跳转      http://tb.com/login?return_to=http%3A%2F%2Ftb.com%2Findex
重定向    http://idp.com/check_login?return_to=http%3A%2F%2Ftb.com%2Flogin%3Ftoken%3Dtokenvalue%26return_to%3Dhttp%253A%252F%252Ftb.com%252Findex
重定向    http://tb.com/login?token=tokenvalue&return_to=http%3A%2F%2Ftb.com%2Findex
访问      http://idp.com/verify_token?token=tokenvalue
创建session
重定向    http://tm.com/index

[sp session]
sessionid

[redis]
sessionid -> uid name

[idp session]
sessionid

[idp redis]
sessionid -> token

[logout]
http://idp.com/logout?raw=url