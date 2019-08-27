package handler

import (
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/GetUserHouses/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	data,err:=DataManipulation.GetSession(req.SessionId+"user_id")
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	userId,_:=redis.Int(data,nil)
	//查询house
	 houses:= []models.House{}
	err=DataManipulation.GetUserHousesTable(userId,&houses)

	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}

	//数据编码
	house_data,_:=json.Marshal(houses)
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.AllData=house_data
	return nil
}