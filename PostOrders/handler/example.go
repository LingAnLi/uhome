package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/utils"

	example "uhome/PostOrders/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostOrders(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("api/v1.0/orders")
	sessionid_userid:=req.SessionId+"user_id"
	valueId,err:=DataManipulation.GetSession(sessionid_userid)
	if err!=nil{
		rsp.Errno=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil

	}
	userId:=int(valueId.([]uint8)[0])

	reqMap:=make(map[string]interface{})
	json.Unmarshal(req.Body,&reqMap)
	if reqMap["house_id"]==""||reqMap["start_date"]==""||reqMap["end_date"]==""{
		rsp.Errno=utils.RECODE_REQERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}

	//
	startDateTime,_:=time.Parse("2006-01-02 15:04:05",reqMap["start_date"].(string)+"00:00:00")
	endDateTime,_:=time.Parse("2006-01-02 15:04:05",reqMap["end_date"].(string)+"00:00:00")
	//
	days:=endDateTime.Sub(startDateTime).Hours()/24
	
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	return nil
}
