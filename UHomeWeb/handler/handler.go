package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"time"
	DELETESESSION "uhome/DeleteSession/proto/example"
	GETAREA "uhome/GetArea/proto/example"
	GETIMAGECD "uhome/GetImageCd/proto/example"
	GETSESSION "uhome/GetSession/proto/example"
	GETSMSCD "uhome/GetSmsCd/proto/example"
	GETUSERHOUSES "uhome/GetUserHouses/proto/example"
	GETUSERINFO "uhome/GetUserInfo/proto/example"
	POSTAVATAR "uhome/PostAvatar/proto/example"
	POSTHOUSEIMG "uhome/PostHouseImg/proto/example"
	POSTUSERHOUSES "uhome/PostHouses/proto/example"
	POSTLOGIN "uhome/PostLogin/proto/example"
	POSTRET "uhome/PostRet/proto/example"
	POSTUSERAUTH "uhome/PostUserAuth/proto/example"
	GETHOUSEINFO "uhome/GetHouseInfo/proto/example"
	GETINDEX "uhome/GetIndex/proto/example"
	GETHOUSES "uhome/GetHouses/proto/example"
	POSTORDERS"uhome/PostOrders/proto/example"
	GETUSERORDER"uhome/GetUserOrder/proto/example"
	PUTORDERS"uhome/PutOrders/proto/example"
	PUTCOMMENT"uhome/PutComment/proto/example"

	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
	example "micro/rpc/srv/proto/example"
)

//获取地区信息
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	beego.Info("请求地区信息 GetArea api/v1.0/areas")
	//创建服务获取句柄
	server := grpc.NewService()
	//服务初始化
	server.Init()

	//调用服务返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())

	//调用服务返回数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//接收数据
	//准备接收切片
	area_list :=[]models.Area{}
	//循环接收数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id:int(value.Aid),Name:value.Aname}
		area_list=append(area_list,tmp)
	}

	// 返回给前端的map
	response := map[string]interface{}{
		"errno": rsp.Error,
		"errmsg": rsp.Errmsg,
		"data":area_list,
	}

	//回传数据的时候是直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}


//获取首页轮播图
func GetIndex(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	server:=grpc.NewService()
	server.Init()
	exampleClient:=GETINDEX.NewExampleService("go.micro.srv.GetIndex",server.Client())
	rsp,err:=exampleClient.GetIndex(context.TODO(),&GETINDEX.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data:=[]interface{}{}
	json.Unmarshal(rsp.Data,&data)
	// we want to augment the response

	test:=make(map[string]string)
	test["house_id"]="8"
	test["img_url"]= "http://127.0.0.1:8889/group1/M00/00/00/CtM3BV1EYHSAE3ToAACw-tYiDKA252.png"
	tpm:=make(map[string]interface{})
	tpm["houses"]=data

	response := map[string]interface{}{
		"errno": utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":data,
	}
	//设置发送格式

	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//获取session
 func GetSession(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	//cookie获取
	cookie,err:=r.Cookie("userlogin")
	if err!=nil{
		//返回未登入
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}



	server:=grpc.NewService()
	server.Init()
	exampleClient:=GETSESSION.NewExampleService("go.micro.srv.GetSession",server.Client())
	rsp,err:=exampleClient.GetSession(context.TODO(),&GETSESSION.Request{
		SessionId:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data:=make(map[string]string)
	data["name"]=rsp.UserName
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}

	//设置发送格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//获取验证码
func GetImageCd(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {

	beego.Info("获取验证码  GetImageCd  /api/v1.0/imagecode/:uuid")
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", server.Client())

	uuid:=ps.ByName("uuid")

	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid:uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var img image.RGBA
	img.Stride=int(rsp.Stride)
	img.Pix=[]uint8(rsp.Pix)
	img.Rect.Max.X=int(rsp.Max.X)
	img.Rect.Max.Y=int(rsp.Max.Y)
	img.Rect.Min.X=int(rsp.Min.X)
	img.Rect.Min.Y=int(rsp.Min.Y)
	var image captcha.Image
	image.RGBA=&img

	png.Encode(w,image)
}
//获取短信验证码
func GetSmsCd(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	//Get UUid & Text  mobile
	test:=r.URL.Query()["text"][0]
	id:=r.URL.Query()["id"][0]
	mobile:=ps.ByName("mobile")
	beego.Info("获取验证码  GETSMSCD")
	//校验手机号码
	//mobile_reg:=regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)

	//bl:=mobile_reg.MatchString(mobile)
	//if !bl{
	//	response:=map[string]interface{}{
	//		"error":utils.RECODE_MOBILEERR,
	//		"errmsg":utils.RecodeText(utils.RECODE_MOBILEERR),
	//	}
	//	w.Header().Set("Content-Type","application/json")
	//	if err:= json.NewEncoder(w).Encode(response);err!=nil{
	//		http.Error(w,err.Error(),500)
	//		return
	//	}


	//}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmsCd", server.Client())



	rsp, err := exampleClient.GetSmsCd(context.TODO(), &GETSMSCD.Request{
		Mobile:mobile,
		Imagestr:test,
		Uuid:id,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response:=map[string]interface{}{
		"erron":rsp.Error,
		"errmsg":rsp.Errmsg,
	}
	w.Header().Set("Content-Type","application/json")
	if err:= json.NewEncoder(w).Encode(response);err!=nil{
		http.Error(w,err.Error(),500)
		return
	}
}
//注册
func PostRet(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	beego.Info("PostRet  注册")
	server:= grpc.NewService()
	server.Init()
	//获取数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	//校验数据
	if request["mobile"].(string)==""||request["password"].(string)==""||request["sms_code"].(string)==""{
		// we want to augment the response
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//回传数据的数据格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		return


	}
	// call the backend service
	exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", server.Client())
	rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
			Mobile:request["mobile"].(string),
			Passwd:request["password"].(string),
			SmsCode:request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	//read cookie
	cookie,err:=r.Cookie("userlogin")
	if err!=nil || ""==cookie.Value{
		//creat cookie
		cookie:=http.Cookie{Name:"userlogin",Value:rsp.SessionId,Path:"/",MaxAge:3600}
		http.SetCookie(w,&cookie)
	}




	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Erron,
		"errmsg": rsp.Errmsg,
	}

	//回传数据的时候是直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//登陆
func PostLogin(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//request["mobile"]==nil request["password"]==nil
	if request["mobile"]==nil ||request["password"]==nil{
		response := map[string]interface{}{
			"errno": utils.RECODE_MOBILEERR,
			"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", server.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Mobile:request["mobile"].(string),
		Passwd:request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//写入cookie
	cookie:=http.Cookie{Name:"userlogin",Value:rsp.SessionId,MaxAge:60*60*12*7}
	http.SetCookie(w,&cookie)

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//退出
func DeleteSession(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	beego.Info("EXIT USER")
	//获取cookie
	rcookie,err:=r.Cookie("userlogin")
	if err!=nil||rcookie.Value==""{
		//返回
		response := map[string]interface{}{
			"errno": utils.RECODE_OK,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	cookie:=http.Cookie{Name:"userlogin",Value:rcookie.Value,MaxAge:-1}
	http.SetCookie(w,&cookie)
	// call the backend service
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", client.DefaultClient)
	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		SessionId:rcookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	//获取cookie
	cookie,err:=r.Cookie("userlogin")
	if err!=nil||cookie.Value==""{
		//返回
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		SessionId:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=make(map[string]string)
	data["user_id"]=rsp.UserId
	data["name"]=rsp.Name
	data["mobile"]=rsp.Mobile
	data["real_name"]=rsp.RealName
	data["id_card"]=rsp.IdCard
	data["avatar_url"]=rsp.AvatarUrl
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	beego.Info(response)
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//上传头像
func PostAvatar(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// 获取图片数据
	f,f_h,err:=r.FormFile("avatar")
	if err!=nil{
		response := map[string]interface{}{
			"errno": utils.RECODE_GETDATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_GETDATAERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	//获取cookie
	cookie,err:=r.Cookie("userlogin")
	if err!=nil||cookie.Value==""{
		//返回
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	//获取文件数据字节流
	fbuff:=make([]byte,f_h.Size)
	_,err=f.Read(fbuff)
	if err!=nil{
		//返回
		response := map[string]interface{}{
			"errno": utils.RECODE_GETDATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_GETDATAERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	f.Close()
	//文件名
	fname:=f_h.Filename

	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTAVATAR.NewExampleService("go.micro.srv.PostAvatar", server.Client())
	rsp, err := exampleClient.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		SessionId:cookie.Value,
		Ext:path.Ext(fname)[1:],
		Size:f_h.Size,
		Avatar:fbuff,


	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data:=make(map[string]string)
	data["avatar_url"]=utils.AddDomain2Url(rsp.AvatarUrl)
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//获取用户信息
func GetUserAuth(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	//获取cookie
	cookie,err:=r.Cookie("userlogin")
	if err!=nil||cookie.Value==""{
		//返回
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置发送格式
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		SessionId:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	data:=make(map[string]string)
	data["user_id"]=rsp.UserId
	data["name"]=rsp.Name
	data["mobile"]=rsp.Mobile
	data["real_name"]=rsp.RealName
	data["id_card"]=rsp.IdCard
	data["avatar_url"]=rsp.AvatarUrl
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	beego.Info(response)
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//用户实名认证
func PostUserAuth(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	idCard:=request["id_card"].(string)

	//正则匹配idCard
	re,_:=regexp.Compile(`^(\d{6})(\d{4})(\d{2})(\d{2})(\d{3})([0-9]|X)$`)
	bl:=re.MatchString(idCard)
	beego.Info(bl,idCard)
	if !bl{
		beego.Info(bl,idCard)
		response := map[string]interface{}{
			"errno": utils.RECODE_IDCARDERR,
			"errmsg": utils.RecodeText(utils.RECODE_IDCARDERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	cookie,err:=r.Cookie("userlogin")
	if err!=nil||cookie.Value=="" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTUSERAUTH.NewExampleService("go.micro.srv.PostUserAuth", server.Client())
	rsp, err := exampleClient.PostUserAuth(context.TODO(), &POSTUSERAUTH.Request{
		RealName: request["real_name"].(string),
		IdCard:request["id_card"].(string),
		SessionId:cookie.Value,

	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//用户以发布房屋查看
func GetUserHouses(w http.ResponseWriter, r *http.Request,_ httprouter.Params ) {
	//获取sessionId
	sessionId,err:=ToObtainSessionId(r)
	if err!=nil{
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETUSERHOUSES.NewExampleService("go.micro.srv.GetUserHouses", server.Client())
	rsp, err := exampleClient.GetUserHouses(context.TODO(), &GETUSERHOUSES.Request{
		SessionId:sessionId,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	houses_list:=[]models.House{}
	json.Unmarshal(rsp.AllData,&houses_list)
	var house[] interface{}
	for _, value := range houses_list {
		//获取需要返回的数据
		house=append(house,value.To_house_info())
	}
	data :=make(map[string]interface{})
	data["houses"]=house


	// we want to augment the response
	response := map[string]interface{}{
		"errno":rsp.Errno,
		"errmsg":rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//发布房源
func PostHouses(w http.ResponseWriter, r *http.Request,_ httprouter.Params ) {
	beego.Info("发布房源")
	body,_:=ioutil.ReadAll(r.Body)
	sessionId,err:=ToObtainSessionId(r)
	if err!=nil{
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTUSERHOUSES.NewExampleService("go.micro.srv.PostHouses", server.Client())
	rsp, err := exampleClient.PostHouses(context.TODO(), &POSTUSERHOUSES.Request{
		SessionId:sessionId,
		Body:body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data:=make(map[string]interface{})
	data["house_id"]=rsp.HouseId
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//上传房源图片
func PostHouseImg(w http.ResponseWriter, r *http.Request,ps httprouter.Params ) {

	houseId:=ps.ByName("id")
	f,f_h,err:=r.FormFile("house_image")
	if err!=nil{
		response := map[string]interface{}{
			"errno": utils.RECODE_GETDATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_GETDATAERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	houseImg:=make([]byte,f_h.Size)
	fname:=f_h.Filename
	ext:=path.Ext(fname)
	f.Read(houseImg)
	f.Close()
	sessionId,err:=ToObtainSessionId(r)
	if err!=nil{
		response := map[string]interface{}{
			"errno": utils.RECODE_GETDATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_GETDATAERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTHOUSEIMG.NewExampleService("go.micro.srv.PostHouseImg", server.Client())
	rsp, err := exampleClient.PostHouseImg(context.TODO(), &POSTHOUSEIMG.Request{
		HouseId:houseId,
		HouseImg:houseImg,
		SessionId:sessionId,
		Size:f_h.Size,
		Ext:ext[1:],
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data:=make(map[string]string)

	data["url"]=utils.AddDomain2Url(rsp.Url)
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data": data,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//GetHouseInfo 获取房源详细信息
func GetHouseInfo(w http.ResponseWriter, r *http.Request,ps httprouter.Params ) {
	//获取cookie和房屋id
	sessionId,err:=ToObtainSessionId(r)
	houseID:=ps.ByName("id")
	if err!=nil{
		response := map[string]interface{}{
			"errno": utils.RECODE_GETDATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_GETDATAERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETHOUSEINFO.NewExampleService("go.micro.srv.GetHouseInfo", server.Client())
	rsp, err := exampleClient.GetHouseInfo(context.TODO(), &GETHOUSEINFO.Request{
		SessionId:sessionId,
		HousesId:houseID,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var house_desc map[string]interface{}
	json.Unmarshal(rsp.Data,&house_desc)
	//urls:=house_desc["img_urls"].([]interface{})
	//for i, value := range urls {
	//	urls[i]=utils.AddDomain2Url(value.(string))
	//}
	data := make(map[string]interface{})
	data["house"]=house_desc
	data["user_id"]=rsp.UserId
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetHouses(w http.ResponseWriter, r *http.Request,_ httprouter.Params ) {
	server:= grpc.NewService()
	server.Init()
	// call the backend service
	var Page string
	if len(r.URL.Query()["p"])==0{
		Page="1"
	}else {
		Page=r.URL.Query()["p"][0]
	}
	exampleClient := GETHOUSES.NewExampleService("go.micro.srv.GetHouses", server.Client())
	rsp, err := exampleClient.GetHouses(context.TODO(), &GETHOUSES.Request{
		Aid:r.URL.Query()["aid"][0],
		Sd:r.URL.Query()["sd"][0],
		Ed:r.URL.Query()["ed"][0],
		P:Page,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	houses:=[]interface{}{}
	json.Unmarshal(rsp.Data,&houses)
	data:=map[string]interface{}{}
	data["current_page"]=rsp.CurrentPage
	data["houses"]=houses
	data["total_page"]=rsp.TotalPage
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostOrders(w http.ResponseWriter, r *http.Request,_ httprouter.Params ) {

	body,_:=ioutil.ReadAll(r.Body)
	sessionId,err:=ToObtainSessionId(r)
	if err!=nil{
		response := map[string]interface{}{
			"errno":utils.RECODE_SESSIONERR,
			"errmsg":utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type","application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	service:=grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := POSTORDERS.NewExampleService("go.micro.srv.PostOrders", service.Client())
	rsp, err := exampleClient.PostOrders(context.TODO(), &POSTORDERS.Request{
		SessionId:sessionId,
		Body:body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	order_id:=map[string]interface{}{"order_id":int(rsp.OrderId)}


	// we want to augment the response
	response := map[string]interface{}{
		"errno":rsp.Errno,
		"errmsg":rsp.Errmsg,
		"data":order_id,
	}
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
//获取订单
func GetUserOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	beego.Info("/api/v1.0/user/orders   GetUserOrder 获取订单 ")
	server :=grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := GETUSERORDER.NewExampleService("go.micro.srv.GetUserOrder", server.Client())

	//获取cookie
	userlogin,err:=r.Cookie("userlogin")
	if err != nil{
		resp := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	//获取role
	role := r.URL.Query()["role"][0] //role


	rsp, err := exampleClient.GetUserOrder(context.TODO(), &GETUSERORDER.Request{
		Sessionid:userlogin.Value,
		Role:role,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	order_list := []interface{}{}
	json.Unmarshal(rsp.Orders,&order_list)

	data := map[string]interface{}{}
	data["orders"] = order_list



	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}


	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//房东同意/拒绝订单
func PutOrders(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	// decode the incoming request as json
	//接收请求携带的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
	//获取cookie
	userlogin,err:=r.Cookie("userlogin")
	if err != nil{
		resp := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 502)
			beego.Info(err)
			return
		}
		return
	}
	server:=grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := PUTORDERS.NewExampleService("go.micro.srv.PutOrders", server.Client())

	rsp, err := exampleClient.PutOrders(context.TODO(), &PUTORDERS.Request{
		Sessionid:userlogin.Value,
		Action:request["action"].(string),
		Orderid:ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 504)
		return
	}
}
//用户评价订单
func PutComment(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	beego.Info("PutComment  用户评价 /api/v1.0/orders/:id/comment")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := grpc.NewService()
	service.Init()
	// call the backend service
	exampleClient := PUTCOMMENT.NewExampleService("go.micro.srv.PutComment", service.Client())

	//获取cookie
	userlogin,err:=r.Cookie("userlogin")
	if err != nil{
		resp := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.PutComment(context.TODO(), &PUTCOMMENT.Request{

		Sessionid:userlogin.Value,
		Comment:request["comment"].(string),
		OrderId:ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}
//mo ren
func ExampleCall(w http.ResponseWriter, r *http.Request ) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}