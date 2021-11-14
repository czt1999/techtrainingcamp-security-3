# 抓到你了——具备安全防护能力的账号系统

## 产品需求

设计并开发一个登陆注册系统，具备注册账户、登录、登出和注销账户等功能。同时需要对异常用户进行识别和制定相关风控策略进行拦截。

## 接口

以下各个接口均为POST方法

### 注册(api/register)

**请求参数RegisterRequest**

```go
type RegisterRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PhoneNum     string `json:"phone_number"`
	VerifyCode   string `json:"verify_code"`
	security.Env `json:"environment"`
}

type Env struct {
	IP       string `json:"ip"`
	DeviceID string `json:"device_id"`
}
```

通过风控策略的中间层，根据不同的情况进行拦截。接着检查用户名和手机号，如果手机号或用户名已被注册，将无法注册。在注册时保存相应的环境信息来进行风控策略评估。在`redis`里面保存用户的`session`，并且在`user`数据库保存用户信息。

**返回参数**

```
code 			//返回码，0为成功，1为失败
message 		//返回信息
session_ID 
expire_time 	//session过期时间
decision_type 	//决策类型，决定风控
```

### 用户名登录(api/login/name)

**请求参数LoginByNameRequest**

```go
type LoginByNameRequest struct {
	Username     string `form:"username" json:"username"`
	Password     string `form:"password" json:"password"`
	security.Env `json:"environment"`
}
```

将密码转化为`md5`码与`user`数据库当中的用户信息进行对比，如果存在该用户的记录，`redis`里面保存用户的`session`，同时保存用户的环境信息做风控策略。

**返回参数**

```
code 			//返回码，0为成功，1为失败
message 		//返回信息
session_ID 
expire_time 	//session过期时间
decision_type 	//决策类型，决定风控
```

### 请求验证码（api/applycode)

**请求参数**

```go
type GetApplyCodeRequest struct {
	PhoneNum     string `json:"phone_number"`
	security.Env `json:"environment"`
}
```

判断输入的手机号是否合法，若合法，然后判断风控决策等级，如果为合法用户，则随机生成一个六位数，并设置好过期时间，将其保存到`redis`数据库当中。

**返回参数**

```
code 			//返回码，0为成功，1为失败
message 		//返回信息
verify_code 	//六位随机验证码
expire_time  	//验证码过期时间
decision_type 	//决策类型，决定风控
```

### 手机号登录(api/login/phone)

**请求参数**

```go
type LoginByPhoneRequest struct {
	PhoneNum     string `json:"phone_number"`
	VerifyCode   string `json:"verify_code"`
	security.Env `json:"environment"`
}
```

通过判断手机号和对应的验证码是否存在`user`数据库中，如果存在，检查当前的用户登录环境；如果风控策略通过，用户登陆成功，并在`redis`里设置好相应的`session`和过期时间。否则登录失败。

**返回参数**

```
code 			//返回码，0为成功，1为失败
message 		//返回信息
session_ID 
expire_time 	//session过期时间
decision_type 	//决策类型，决定风控
```

### 登出与注销（api/logout)

**请求参数**

```go
type LogoutRequest struct {
	SessionID    string `json:"session_id"`
	ActionType   int    `json:"action_type"`
	security.Env `json:"environment"`
}
```

登出与注销是两种不同的方式，通过`ActionType`来判断是登出还是注销，如果`ActionType`为1，则为登出操作，用户必须先登录，即`sessionID`对应的记录存在，才能进行登出操作，删除用户的的`session`信息。如果`ActionType`为2，则为注销操作，同样的，用户必须先登录才能注销，注销要将用户的账号删除，处理删除用户的`session`信息外，同时也要将用户的注册信息一并删除。

**返回参数**

```
code 		//返回码，0为成功，1为失败
message 	//返回信息
```

### 通过session得到用户信息（api/user/name)

**请求参数**

```go
type GetUsernameRequest struct {
	SessionID    string `json:"session_id"`
	security.Env `json:"environment"`
}
```

在`session`信息当中查找用户，如果用户`session`存在，则返回用户名，显示用户已登录。实现通过`sessionID`来保持登录状态。

**返回参数**

```
code		//返回码，0为成功，1为失败
message 	//返回信息
username 	//用户名
```



## 风控策略

- 如果用户在一段时间内（T1）请求超过规定次数（N1），将该用户判定存在低风险，将触发滑块验证操作
- 如果用户在一段时间内（T2）触发滑块验证超过规定次数（N2），将该用户判定存在中风险，将触发暂时封禁操作
- 如果同一设备ID在一段时间内（T3）注册/登录用户数量达到规定上限（N3），将该设备ID和相关IP判定为中风险，将触发暂时封禁操作
- 如果同一设备IP或ID在一段时间内（T4）判定中风险达到规定次数（N4），将被判定为高风险，将触发永久封禁操作
- 默认 (T1=2秒,N1=5) , (T2=1小时,N2=10) , (T3=2天,N3=3) , (T4=14天,N4=3)
- 默认暂时封禁的时长TempBlock=1天
- 用户在注销账户后24小时内不能使用相同手机号注册新账号



## 代码架构

### 后端服务

设置风控拦截中间件，并在api路由服务之后设置评估用户的风控等级。

```go
//如果ip，设备号在拦截名单内，不提供api服务，直接断开连接
r.Use(CheckBlocked)

// API
apiG := r.Group("/api")
apiG.POST("/register", api.Register)
apiG.POST("/login/name", api.LoginByName)
apiG.POST("/login/phone", api.LoginByPhone)
apiG.POST("/logout", api.Logout)
apiG.POST("/applycode", api.GetApplyCode)
apiG.POST("/user/name", api.GetUsername)

//根据用户的行为判断是否提升对该用户的风控等级
apiG.Use(AddBlocked)
```



### 前端服务

```
index.html 		// 用户进行登录和注册的页面
main.html 		// 登录成功后的主页面，登出和注销在该页面进行
error.html 		// 发生错误后的页面
```

`index.html`设置了登录界面，包括账户密码登录方式和手机号登录方式，以及可以切换到注册界面，需要填写用户名，密码，手机号，并且获得验证码后才能进行注册。

`main.html`显示出已登录的用户名，提供登出和注销按钮，点击登出或者注销都会返回`index.html`页面

相关按钮设置了冷却时间，避免频繁点击误触发风控拦截。

均使用`ajax`方法传递`json`类型数据来与后端进行交互。



### 代码目录

```
-conf
	app.ini 			//配置文件
-models
	models.go 			//数据库初始化
	user.go 			//User表的定义和增删查改方法
-pkg
	-app
		response.go  	//响应参数的定义
	-gredis
		redis.go  		//redis的初始化和一些方法
	-settings
		settings.go 	//加载配置文件的方法
	-util
		md5.go
-routers
	-api
		applycode.go 	//获取验证码的实现
		user.go 		//接口的实现
	middleware.go 		//中间件的实现
	router.go 			//路由初始化
	router_test.go 		//路由单元测试
-security 				//风控策略的实现
	decision.go         //风控等级判断
	security.go         //风控相关规则的初始化
	session.go          //会话管理接口
-template 				//前端页面的元素
	-pages
		-css 
		-js
		-json
		demo.html
		error.html
		index.html
		main.html
main.go					//web应用主程序
```
