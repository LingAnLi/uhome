package main

import (
    "github.com/julienschmidt/httprouter"
    "github.com/micro/go-log"
    "github.com/micro/go-web"
    "net/http"
    "uhome/UHomeWeb/handler"
    _ "uhome/UHomeWeb/models"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.UHomeWeb"),
                web.Version("latest"),
                web.Address(":8888"),
        )
	// initialise service
	//初始化服务
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

    //使用路由中间件映射页面
    rou:=httprouter.New()

    rou.NotFound=http.FileServer(http.Dir("html"))
    //获取地区请求
    rou.GET("/api/v1.0/areas",handler.GetArea)
    //获取session
    rou.GET("/api/v1.0/session",handler.GetSession)
    //获取首页图
    rou.GET("/api/v1.0/house/index",handler.GetIndex)
    //获取验证码
    rou.GET("/api/v1.0/imagecode/:uuid",handler.GetImageCd)
    //获取短信验证码
    rou.GET("/api/v1.0/smscode/:mobile",handler.GetSmsCd)
    //注册
    rou.POST("/api/v1.0/users",handler.PostRet)
    //登陆
    rou.POST("/api/v1.0/sessions",handler.PostLogin)
    //退出
    rou.DELETE("/api/v1.0/session",handler.DeleteSession)
    //获取用户信息
    rou.GET("/api/v1.0/user",handler.GetUserInfo)
    //头像
    rou.POST("/api/v1.0/user/avatar",handler.PostAvatar)
    //用户认证
    rou.GET("/api/v1.0/user/auth",handler.GetUserAuth)
    //用户实名认证
    rou.POST("/api/v1.0/user/auth",handler.PostUserAuth)
    //用户以发布房屋查看
    rou.GET("/api/v1.0/user/houses",handler.GetUserHouses)
    //发布房源
    rou.POST("/api/v1.0/houses",handler.PostHouses)
    //房源图片 /api/v1.0/houses/6/images
    rou.POST("/api/v1.0/houses/:id/images",handler.PostHouseImg)
    //请求房源详细信息
    rou.GET("/api/v1.0/houses/:id",handler.GetHouseInfo)
    //查找房源信息
    rou.GET("/api/v1.0/houses",handler.GetHouses)
    //发布 orders
    rou.POST("/api/v1.0/orders",handler.PostOrders)
    //get 查看房东/租客订单信息请求
    rou.GET("/api/v1.0/user/orders",handler.GetUserOrder)
    //put房东同意/拒绝订单
    //api/v1.0/orders/:id/status
    rou.PUT("/api/v1.0/orders/:id/status",handler.PutOrders)
    //PUT 用户评价订单信请求
    //api/v1.0/orders/:id/comment
    //api/v1.0/orders/1/comment
    rou.PUT("/api/v1.0/orders/:id/comment",handler.PutComment)


    // register html handle
	//service.Handle("/", http.FileServer(http.Dir("html")))
	service.Handle("/", rou)

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
