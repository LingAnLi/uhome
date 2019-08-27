package handler

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	"github.com/micro/go-log"

	example "uhome/PostUserAuth/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostUserAuth(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//获取身份信息
	idCare,name,sessionId:=req.IdCard,req.RealName,req.SessionId
	//校验信息
	/*
	var CheckIDCard UserAPI.CheckIDCardConf
	CheckIDCard.RealName=name
	CheckIDCard.CardID=idCare
	CheckIDCard.APPCODE="3ed27c00c9a24cbca6e20652c982862c"
	data,err:=CheckIDCard.CheckIDCard()
	if err!=nil||data["status"]!="01"{
		rsp.Errno=utils.RECODE_IDCARDERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info(data)
	*/

	//获取用户id
	userId,err:=DataManipulation.GetSession(sessionId+"user_id")
	if err!=nil{
		rsp.Errno=utils.RECODE_USERERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	userId_int,err:=redis.Int(userId,nil)
	if err!=nil{
		rsp.Errno=utils.RECODE_USERERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//查询用户信息
	user:=models.User{Id:userId_int}
	err=DataManipulation.GetUserData(&user)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//一个身份只能注册一个
	tmp:=models.User{Id_card:idCare}
	err=DataManipulation.GetUserData(&tmp,"id_card")
	if err==nil{
		rsp.Errno=utils.RECODE_HASIDCARD
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//更新数据
	user.Real_name=name
	user.Id_card=idCare
	err=DataManipulation.UserUpData(&user,"")
	//返回ok
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)


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
