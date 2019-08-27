package handler

import (
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/GetHouseInfo/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouseInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//获取信息
	houseId,sessionId:=req.HousesId,req.SessionId
	id,err:=strconv.Atoi(houseId)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	userId,err:=DataManipulation.GetSession(sessionId+"user_id")
	userId_int,_:=redis.Int( userId,nil)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//查询
	house:= models.House{}
	err=DataManipulation.GetHouseInfo(&house,id)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	house_desc:=house.To_one_house_desc().(map[string]interface{})
	data,_:=json.Marshal(&house_desc)
	//返回数据
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.Data=data
	rsp.UserId=int64(userId_int)
	return nil
}

