package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"time"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	"encoding/json"
	"github.com/micro/go-log"
	example "uhome/PostRet/proto/example"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostRet(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("PostRet")
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
		rsp.Erron=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Erron)
		return nil
	}
	//Get SmsCode
	sms_code:=bm.Get(req.Mobile)
	if sms_code==nil{
		beego.Info("get sms code err")
		rsp.Erron=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Erron)

		return nil
	}
	sms_str,_:=redis.String(sms_code,nil)
	//校验数据
	if req.SmsCode!=sms_str{
		beego.Info("sms code err")
		rsp.Erron=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Erron)

		return nil
	}
    //将注册信息加入 mysql
    o:=orm.NewOrm()

    user:=models.User{Mobile:req.Mobile,Password_hash:Md5String(req.Passwd),Name:req.Mobile}
    _,err=o.Insert(&user)
    if err!=nil{
    	beego.Info("zhucheshibai")
		rsp.Erron=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Erron)
		return nil
	}
    //设置session
    sessionId:=Md5String(req.Passwd+req.Mobile)
    bm.Put(sessionId+"mobile",user.Mobile,time.Second*3600)
	bm.Put(sessionId+"user_id",user.Id,time.Second*3600)
	bm.Put(sessionId+"name",user.Name,time.Second*3600)
	//返回数据
    rsp.SessionId=sessionId
	rsp.Erron=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Erron)
	return nil
}

func Md5String(s string) string {
	hash:=md5.Sum([]byte(s))
    return hex.EncodeToString(hash[:])
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
