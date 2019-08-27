package handler

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"uhome/DataManipulation"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"

	"github.com/micro/go-log"

	example "uhome/PostAvatar/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//获取数据
	size:=req.Size
	ext:=req.Ext
	avatar:=req.Avatar
	sessionId:=req.SessionId
	//校验数据
	if int(size)!=len(avatar){
		rsp.Errno=utils.RECODE_GETDATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//查询用户Id
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
	//上传数据
	f_id,err:=DataManipulation.UploadByBuffer(&avatar,ext)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	//更新数据库
	user.Avatar_url=f_id
	err=DataManipulation.UserUpData(&user,"avatar_url")
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Errno)
	rsp.AvatarUrl=f_id





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
