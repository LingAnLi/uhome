package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strconv"
	"time"
	example "uhome/GetSmsCd/proto/example"
	"uhome/UHomeWeb/UserAPI"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSmsCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("/api/v1.0/smscode/:mobile  GetSmsCd")
	rsp.Error=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Error)
	//mysql has this mobile?
	o:=orm.NewOrm()
	user:=models.User{Mobile:req.Mobile}
	err:=o.Read(&user,"Mobile")
	//
	if err==nil{
		beego.Info("yong hu yi cun zai")
		rsp.Error=utils.RECODE_MOBILEERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
	}

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
		rsp.Error=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
		return nil
	}
	v:=bm.Get(req.Uuid)
	if v ==nil{
		beego.Info("redis Get Err")
		rsp.Error=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
		return nil
	}

	v_str,_:=redis.String(v,nil)
	if v_str!=req.Imagestr{
		beego.Info("imgstr err")
		rsp.Error=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)

		return nil
	}
	//Send SMS
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	size:=r.Intn(9999)+1001
		var conf UserAPI.SmsConf
		conf.UserId="2019076290"
		conf.Clientid="b04t18"
		conf.Password="12345678"
		conf.Mobile=req.Mobile
		conf.Smstype="0"
		conf.Content="【IHome注册】您的验证码是"+strconv.Itoa(size )+"，如非本人操作，请忽略此条短信。"
		//m:=conf.SendSms()
		//if m.Total_fee!=0{
		//	rsp.Error=utils.RECODE_MOBILEERR
		//	rsp.Errmsg=utils.RecodeText(rsp.Error)
		//	return nil
		//}
		beego.Info(size)
		bm.Put(req.Mobile,size,time.Second*300)
	return nil
}
