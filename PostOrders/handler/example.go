package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
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
		rsp.Errmsg="session err"
		return nil

	}
	userId,_:=redis.Int(valueId,nil)

	reqMap:=make(map[string]interface{})
	json.Unmarshal(req.Body,&reqMap)
	if reqMap["house_id"]==""||reqMap["start_date"]==""||reqMap["end_date"]==""{
		rsp.Errno=utils.RECODE_REQERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}

	beego.Info("data",reqMap)
	//获取订单起止时间
	startDateTime,_:=time.Parse("2006-01-02",reqMap["start_date"].(string))
	endDateTime,_:=time.Parse("2006-01-02 ",reqMap["end_date"].(string))
	beego.Info("startDate:",startDateTime)
	//计算订单天数
	days:=endDateTime.Sub(startDateTime).Hours()/24

	house_id, _ := strconv.Atoi(reqMap["house_id"].(string))
	//房屋对象
	house := models.House{Id: house_id}
	o := orm.NewOrm()
	if err := o.Read(&house); err != nil {
		rsp.Errno  =  utils.RECODE_NODATA
		rsp.Errmsg  = utils.RecodeText(rsp.Errno)
		return nil
	}
	o.LoadRelated(&house, "User")

	//确保当前的uers_id不是房源信息所关联的user_id
	if userId == house.User.Id {


		rsp.Errno  =  utils.RECODE_ROLEERR
		rsp.Errmsg  = utils.RecodeText(rsp.Errno)

		return nil
	}
	//确保用户选择的房屋未被预定,日期没有冲突
	if endDateTime.Before(startDateTime) {

		rsp.Errno  =  utils.RECODE_ROLEERR
		rsp.Errmsg  = "结束时间在开始时间之前"
		return nil
	}
	//添加征信步骤




	//封装order订单
	amount := days * float64(house.Price)
	order := models.OrderHouse{}
	order.House = &house
	user := models.User{Id: userId}
	order.User = &user
	order.Begin_date = startDateTime
	order.End_date = endDateTime
	order.Days = int(days)
	order.House_price = house.Price
	order.Amount = int(amount)
	order.Status = models.ORDER_STATUS_WAIT_ACCEPT
	//征信
	order.Credit = false

	//将订单信息入库表中
	if _, err := o.Insert(&order); err != nil {
		rsp.Errno  =  utils.RECODE_DBERR
		rsp.Errmsg  = err.Error()
		return nil
	}
	//返回order_id

	rsp.OrderId = int64(order.Id)
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	return nil
}
