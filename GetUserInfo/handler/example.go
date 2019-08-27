package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	"encoding/json"
	"github.com/micro/go-log"

	example "uhome/GetUserInfo/proto/example"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example)GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	sessionid:=req.SessionId
	//连接redis
	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		beego.Info("redis连接失败",err)
		rsp.Errno=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	userId:=bm.Get(sessionid+"user_id")
	if userId==nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	userId_int,err:=redis.Int(userId,nil)
	//查询数据
	var user models.User
	user.Id=userId_int
	o:=orm.NewOrm()
	err=o.Read(&user)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.AvatarUrl=utils.AddDomain2Url(user.Avatar_url)
	rsp.IdCard=user.Id_card
	rsp.RealName=user.Real_name
	rsp.Mobile=user.Mobile
	rsp.UserId=strconv.Itoa(user.Id)
	rsp.Name=user.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
