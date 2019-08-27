package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/GetIndex/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetIndex(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//从缓存获取

	beego.Info("redis")
	var data interface{}
	err:=DataManipulation.GetIndexByRedis("home_page_data",&data)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	if data!=nil{
		rsp.Errno=utils.RECODE_OK
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		data_bytes,_:=redis.Bytes(data,nil)
		rsp.Data=data_bytes
		return nil
	}
	//从mysql获取
	beego.Info("mysql")
	var house []models.House
	err=DataManipulation.GetIndexByMysql(2,&house)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//处理数据

	houses:=make([]interface{},0)
	for _, value := range house {
		house_info:=value.To_house_info()
		houses=append(houses,house_info)
	}
	houses_json,_:=json.Marshal(houses)
	//
	DataManipulation.PutIndexByRedis("home_page_data",&houses_json)


	rsp.Errno=utils.RECODE_DATAERR
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.Data=houses_json

	return nil
}



