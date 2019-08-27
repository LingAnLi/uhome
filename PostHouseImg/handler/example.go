package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"uhome/DataManipulation"
	_"github.com/go-sql-driver/mysql"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	example "uhome/PostHouseImg/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouseImg(ctx context.Context, req *example.Request, rsp *example.Response) error {
	ext,size,sessionId,houseImg,houseId:=req.Ext,req.Size,req.SessionId,req.HouseImg,req.HouseId
	beego.Info("server 添加图片")
	//获取session
	_,err:=DataManipulation.GetSession(sessionId+"user_id")
	if err!=nil{

		return nil
	}
	//图片验证
	if int(size)!=len(houseImg){

		return nil
	}
	//图片储存
	f_id,err:=DataManipulation.UploadByBuffer(&houseImg,ext)
	if err!=nil{

		return nil
	}
	//加入数据库
	var house models.House
	house.Id,_= strconv.Atoi(houseId)
	o:=orm.NewOrm()
	o.Read(&house)
	if house.Index_image_url==""{
		house.Index_image_url=f_id
		o.Update(&house)
	}

		var HouseImage models.HouseImage
		HouseImage.House=&house
		HouseImage.Url=f_id

		o.Insert(&HouseImage)


	//返回数据
	rsp.Url=f_id
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)


	return nil
}
