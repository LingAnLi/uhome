// https://www.e253.com/
package UserAPI

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)
type Mydata struct {

	Code   int
	Msg string
	Fee int
	Mobile string
	Sid string
	Uid string
}
type ResJsonlice struct {
	Total_fee int
	Data []Mydata

}
//
type SmsConf struct {
	UserId string //用户Id
	Clientid string//帐号，6位
	Password string//密码，8－12位，bu加密
	//支持多号码，号码之间用英文逗号隔开，最多100个。
	Mobile string//发送手机号码，国内短信不要加前缀，国际短信号码前须带相应的国家区号，如日本：0081
	Smstype string//短信类型， "0"：通知短信，"4"：验证码短信，"5"：营销短信
	Content string//【签名】+ 短信内容 ex 【云通讯】您的验证码为：1234"
	SendTime string//为空表示立即发送，定时发送格式2016-11-11
}
func (this *SmsConf)SendSms() ResJsonlice {
	var s ResJsonlice
	var tmp ResJsonlice
	var data Mydata
	//获取信息
	url:="https://u.smsyun.cc/sms-partner/access/"+this.UserId+"/sendsms"
	Sms:=make(map[string]string)
	pwdMd5:=md5.Sum([]byte(this.Password))
	Sms["clientid"]=this.Clientid
	Sms["password"]=hex.EncodeToString(pwdMd5[:])
	Sms["mobile"]=this.Mobile
	Sms["smstype"]=this.Smstype
	Sms["content"]=this.Content
	//
	sms_js,_:=json.Marshal(Sms)
	req,err:= http.NewRequest("POST", url,bytes.NewReader(sms_js))
	if err!=nil{
		fmt.Println("NewRequest err",err)
		data.Code=-1
		data.Msg=err.Error()
		tmp.Data=append(tmp.Data,data)
		return tmp
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err!=nil{
		fmt.Println("client.Do err",err)
		data.Code=-1
		data.Msg=err.Error()
		tmp.Data=append(tmp.Data,data)
		return tmp
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		data.Code=-1
		data.Msg=err.Error()
		tmp.Data=append(tmp.Data,data)
		return tmp
	}

	json.Unmarshal(respBytes, &s)

	for _,data:=range s.Data{
		if data.Code!=0{
			tmp.Total_fee=-1
			tmp.Data=append(tmp.Data,data)
		}
	}
		fmt.Println(tmp)
		return tmp

}
func Exp() {
	var conf SmsConf
	conf.UserId="2019076290"
	conf.Clientid="b04t18"
	conf.Password="12345678"
	conf.Mobile="15179379271"
	conf.Smstype="0"
	conf.Content="【IHome注册】您的验证码是1234，如非本人操作，请忽略此条短信。"
	s:=conf.SendSms()
	if s.Total_fee!=0{
		fmt.Println(s)
	}
}