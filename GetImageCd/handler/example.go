package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"image/color"
	"time"
	"uhome/UHomeWeb/utils"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/gomodule/redigo/redis"
	_"github.com/garyburd/redigo/redis"
	"github.com/micro/go-log"

	example "uhome/GetImageCd/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImageCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取验证码  GetImageCd  /api/v1.0/imagecode/:uuid")
	//验证码生成
	cap:=captcha.New()
	//字体
	if err:=cap.SetFont("comic.ttf");err!=nil{
		beego.Info("字体设置失败",err)
	}
	//图片大小
	cap.SetSize(90,41)
	//干扰强度
	cap.SetDisturbance(captcha.NORMAL)
	//前景色和背景色
	cap.SetFrontColor(color.RGBA{255,255,255,255})
	cap.SetBkgColor(color.RGBA{255,0,0,255},color.RGBA{0,0,255,255})
	//生成图片
	img,str:=cap.Create(4,captcha.NUM)
	//缓存到redis
	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	beego.Info(redis_conf,str)
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		beego.Info("redis连接失败",err)
		rsp.Error=utils.RECODE_DBERR
		rsp.Errmsg=utils.RecodeText(rsp.Error)
	}

	bm.Put(req.Uuid,str,time.Second*300)
	//图片解引用
	img1:=*img
	img2:=*img1.RGBA
	rsp.Error=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(rsp.Error)
	rsp.Pix=img2.Pix
	rsp.Stride=int64(img2.Stride)
	rsp.Max=&example.Response_Point{X:int64(img2.Rect.Max.X),Y:int64(img2.Rect.Max.Y)}
	rsp.Min=&example.Response_Point{X:int64(img2.Rect.Min.X),Y:int64(img2.Rect.Min.Y)}
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
