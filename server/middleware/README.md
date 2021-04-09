# Middleware 包
middleware包含mo2中用到的大部分中间件  
所有middleware中的中间件都应该解耦，不要包含不必要的依赖  

## 获取
```bash
go get github.com/Monkey-Mouse/mo2/server/middleware@release-v0.2.1
```


## 内容

### auth middleware
auth是目前唯一的一个中间件，它能实现基于role的身份认证，以及基本的rate limit功能，能在一定程度上抵御ddos攻击  


#### QuickStart
auth中间件实现了一组类似gin的api，它的使用方法类似于gin的router  
使用实例： https://github.com/Monkey-Mouse/mo2-compose/blob/main/watcher.go  
```go
func setupHandlers(c *controller.Controller) {
    // api组里的api在被访问时都会检查用户是否有User身份
	api := middleware.H.Group("/api", "User")
	{
        // 任何Group下的子路由和Group中的role规则是与逻辑
        // 下面这个路由的访问者必须要有User身份，因为在api定义的Group方法里传入了"User"
		api.GET("/logs1", c.Log)
        // 强迫这个路由的访问者在User身份的基础上同时具有Admin身份
		api.GET("/logs", c.Log, "Admin")
	}
    // 可以传多个需要检验的Role，任何平行输入的role在检查时是或逻辑
    // 下方这个api会检查用户是否具有Admin和User其中之一的Role
    middleware.H.GET("/logs2", c.Log, "Admin", "User")
    // Group也是
    api1 := middleware.H.Group("/api1", "User", "Admin")
    
    /* 
    所有需要使用Ratelimit功能保护的api只需要在常规的注册方法后加入WithRateLimit就是了加
    ratelimit功能的方法多接受一参数，就是第三个参数。它是一个数字，
    代表在一个周期内一个ip对该方法请求次数的上限如果超过这个上限，这个ip会被ban。
    周期长度默认10秒，ban时间默认3600秒（1小时）
    也就是说下边这种写法意义是：/logs2这个api最多被一个相同ip在10秒内请求30次，
    如果10秒不到的时间内请求次数达到30，这个ip会被ban 1个小时。
    周期长度和ban的时间可以使用SetupRateLimiter(limitEvery int, unblockevery int, useRedis bool)方法设置
    */
    middleware.H.GETWithRateLimit("/logs2", c.Log, 30, "Admin", "User")
    // 同理，group也有ratelimit版本
    api2 := middleware.H.GroupWithRateLimit("/api2", 30, "User", "Admin")

}
```
注意，在将handler通过类似上方的方法注册**后**，需要手动调用
```go
middleware.H.RegisterMapedHandlers(r, 
    func(ctx *gin.Context) (userInfo middleware.RoleHolder, err error) {
        str, err := ctx.Cookie("jwtToken")
        if err != nil {
            return
        }
        userInfo, err = mo2utils.ParseJwt(str)
        return
    }, mo2utils.UserInfoKey, &middleware.OptionalParams{LimitEvery: 10, Unblockevery: 3600, UseRedis: true})
```
**只有这样中间件和路由才会真正被注册入gin的router中，调用完注册方法后不需要再向gin的router里注册之前的路由！**  
重要方法`func (h handlerMap) RegisterMapedHandlers(r *gin.Engine, GETUserFromCTX FromCTX, userKey string, optional *OptionalParams)`  
参数解释：  
- r：指向需要注册到的`gin.Engine`的指针
- GETUserFromCTX：一个方法，接收gin的context，从里边产生出一个`RoleHolder`接口类型的数据和
- userKey：一个常量，读取的用户信息会被用`ctx.Set(userKey,info)`存在ctx中，方便之后在其它handler中使用  
- optional: 一个指针，用于传入可选参数，用于配置ratelimiter的参数和是否使用redis

> Ratelimit功能默认会使用服务器内存进行缓存，周期长度默认10秒，ban时间默认3600秒（1小时）。
> 它也支持使用redis进行缓存，注意中间件会试图从环境变量`REDIS_URL`里读取redis的url，
> 从`REDIS_PASS`中读取redis密码，从`REDIS_DB`中获得想要使用的redis数据库编号





