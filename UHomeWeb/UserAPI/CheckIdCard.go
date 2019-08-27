package UserAPI
//Buy API Addr https://market.aliyun.com/products/57000002/cmapi029522.html?spm=5176.2020520132.101.3.3c7b7218z4f8ml#sku=yuncode2352200001
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
curl -i -k --get --include 'https://idcardcert.market.alicloudapi.com/idCardCert?idCard=510703198602170052&name=伍帅'
-H 'Authorization:APPCODE 你自己的AppCode'
*/
type CheckIDCardConf struct {
	APPCODE string
	RealName string
	CardID string
	//Url string //https://idcardcert.market.alicloudapi.com/idCardCert
}
func (this *CheckIDCardConf)CheckIDCard()(information map[string]string,err error){

	url:="https://idcardcert.market.alicloudapi.com/idCardCert"+"?idCard=" + this.CardID + "&name="+ this.RealName
	//url := "http://checkidc.market.alicloudapi.com/idcard/VerifyIdcardv2?cardNo=" + cardNo + "&realName=" + realName


	req,err:= http.NewRequest("", url,nil)
	if err!=nil{
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Authorization","APPCODE"+" "+this.APPCODE)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err!=nil{
		return
	}
	respBytes, err:= ioutil.ReadAll(resp.Body)
	if err!=nil{
		return
	}

	err=json.Unmarshal(respBytes, &information)
	if err!=nil{
		return
	}


return

}
func test()  {
	var conf CheckIDCardConf
	conf.RealName="name"
	conf.CardID="IDCard"
	conf.APPCODE=""
	res,_:=conf.CheckIDCard()
	fmt.Println(res)
}