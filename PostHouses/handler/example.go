package handler

import (
	"context"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/PostHouses/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//获取userID
	//req.SessionId
	userId,err:=DataManipulation.GetSession(req.SessionId+"user_id")
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	id,_:=redis.Int(userId,nil)
	var request = make(map[string]interface{})
	json.Unmarshal(req.Body,&request)
	//获取房屋信息
	house:=models.House{}
	house.Title=request["title"].(string)
	house.Price,_=strconv.Atoi(request["price"].(string))
	house.Address=request["address"].(string)
	house.Room_count,_=strconv.Atoi(request["room_count"].(string))
	house.Acreage,_=strconv.Atoi(request["acreage"].(string))
	house.Unit=request["unit"].(string)
	house.Capacity,_=strconv.Atoi(request["capacity"].(string))
	house.Beds=request["beds"].(string)
	deposit,_:=strconv.Atoi(request["deposit"].(string))
	house.Deposit=deposit*100
	house.Min_days,_=strconv.Atoi(request["min_days"].(string))
	house.Max_days,_=strconv.Atoi(request["max_days"].(string))
	area_id,_:=strconv.Atoi(request["area_id"].(string))
	var area = models.Area{Id:area_id}
	house.Area=&area
	//获取用户信息
	var user models.User
	user.Id=id
	house.User=&user

	houseId,err:=DataManipulation.InsertHouse(&house)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//
	var facility []*models.Facility
	for _, value := range request["facility"].([]interface{}) {
		fid,_:=strconv.Atoi(value.(string))
		ftmp:=&models.Facility{Id:fid}
		facility=append(facility,ftmp)
	}
	err=DataManipulation.M2MInsertHouse(&house,&facility)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//返回数据
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.HouseId=strconv.Itoa(houseId)
	return nil
}
