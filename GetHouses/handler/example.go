package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/GetHouses/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("GetHousesService")
    if req.Aid ==""{
    	rsp.Errno=utils.RECODE_DATAERR
    	rsp.Errmsg=utils.RecodeText(rsp.Errno)
    	return nil
	}
	//huo qu shu ju
	startData:=req.Sd
	endData:=req.Ed
	aid,_:=strconv.Atoi(req.Aid)//addr id
	sk:=req.Sk
	page,_:=strconv.Atoi(req.P)
	beego.Info(startData,endData,aid,sk,page)

	//huo qu fang wu shu ju
	houses:=[]models.House{}

	o:=orm.NewOrm()
	qs:=o.QueryTable("house")
	num,err:=qs.Filter("area_id",aid).Count()
	if err!=nil{
		rsp.Errno=utils.RECODE_PARAMERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	total_page:=int(math.Ceil(float64(num)/float64(models.HOUSE_LIST_PAGE_CAPACITY)))
	if page>total_page{
		page=total_page
	}
	qs.Filter("area_id",aid).Limit(models.HOUSE_LIST_PAGE_CAPACITY,(page-1)*models.HOUSE_LIST_PAGE_CAPACITY).All(&houses)
	houses_list:=[]interface{}{}
	for _, house := range houses {
		o.LoadRelated(&house,"Area")
		o.LoadRelated(&house,"User")
		o.LoadRelated(&house,"Images")
		o.LoadRelated(&house,"Facilities")
		houses_list=append(houses_list,house.To_house_info())
	}
	rsp.TotalPage=int64(total_page)
	rsp.CurrentPage=int64(page)
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.Data,_=json.Marshal(houses_list)


	return nil
}

