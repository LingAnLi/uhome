package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"github.com/micro/go-log"
	"time"
	"uhome/PostRet/handler"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	example "uhome/PostLogin/proto/example"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"

)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostLogin(ctx context.Context, req *example.Request, rsp *example.Response) error {
	o:=orm.NewOrm()
	var user models.User
	user.Mobile=req.Mobile
	err:=o.Read(&user,"Mobile")
	pwd_hash:=handler.Md5String(req.Passwd)
	//校验帐号密码
	if err!=nil || pwd_hash!=user.Password_hash{
		rsp.Errno=utils.RECODE_PWDERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//sessionid
	sessionId:=handler.Md5String(req.Passwd+req.Mobile)
	//redis
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
	//写入session
	bm.Put(sessionId+"mobile",user.Mobile,time.Hour*12)
	bm.Put(sessionId+"user_id",user.Id,time.Hour*12)
	bm.Put(sessionId+"name",user.Name,time.Hour*12)
	//返回sessionId和错误码
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.SessionId=sessionId
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
