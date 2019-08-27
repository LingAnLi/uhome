package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	example "uhome/DeleteSession/proto/example"
	"uhome/UHomeWeb/utils"
	"encoding/json"

	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) DeleteSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//获取数据
	sessionId:=req.SessionId
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
	//写入session
	bm.Delete(sessionId+"mobile")
	bm.Delete(sessionId+"name")
	bm.Delete(sessionId+"user_id")
	//返回错误码
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)


	return nil
}
