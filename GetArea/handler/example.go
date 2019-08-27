package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	"github.com/micro/go-log"
	"github.com/astaxie/beego/cache"
	_"github.com/garyburd/redigo/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/astaxie/beego/cache/redis"
	example "uhome/GetArea/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("请求地域信息 GetArea" )
	//初始化错误码
	rsp.Error=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Error)
	//1.获取数据——>缓存 （如果有数据直接发送给前端）
	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	beego.Info(redis_conf)
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		beego.Info("redis连接失败",err)
		rsp.Error=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
	}
	//获取数据------ area_info
	area_value:=bm.Get("area_info")
	//如果有数据直接发送给前端
	if area_value!=nil{
		beego.Info("从缓存获取数据")
		area_map:= []map[string]interface{}{}
		//
		json.Unmarshal(area_value.([]byte),&area_map)
		for _, value := range area_map{
			tmp:=example.Response_Areas{Aid:int32(value["aid"].(float64)),Aname:value["aname"].(string)}
			rsp.Data=append(rsp.Data,&tmp)
		}

		return nil
	}
	beego.Info("从mysql获取数据")
	//2.如没有数据从mysql获取数据
	var area []models.Area
	o:=orm.NewOrm()
	num,err:=o.QueryTable("area").All(&area)
	if err!=nil{
		beego.Info("数据查询失败")
		rsp.Error=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
		return nil
	}
	if num==0 {
		beego.Info("数据库没有数据")
		rsp.Error=utils.RECODE_NODATA
		rsp.Errmsg=utils.RecodeText(rsp.Error)
		return nil
	}
	//3.存到缓存中
	//将获取的数据转化为json
	area_js,_:=json.Marshal(area)
	err=bm.Put("area_info",area_js,time.Second*3600)
	if err!=nil{
		beego.Info("数据缓存失败",err)
		rsp.Error=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
	}

	//4.发送给前端
	for _, value := range area {
		tmp:=example.Response_Areas{Aid:int32(value.Id),Aname:value.Name}
		rsp.Data=append(rsp.Data,&tmp)
	}

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
