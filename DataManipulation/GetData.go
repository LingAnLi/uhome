package DataManipulation

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gomodule/redigo/redis"
	"time"
	"uhome/UHomeWeb/models"
	"uhome/UHomeWeb/utils"
)

func GetSession(key string) (data interface{},err error) {
	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		return
	}
	data=bm.Get(key)

	return
}
func GetIndexByRedis(key string,data *interface{})(err error){

	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		return
	}
	*data=bm.Get(key)
	return
}
func PutIndexByRedis(key string,data *[]byte)(err error){

	redis_conf:=map[string]string{
		"key":utils.G_server_name,
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
		//"password": this is password
	}
	//将map转换成json
	redis_conf_js,_:=json.Marshal(redis_conf)
	//创建redis句炳
	bm,err:=cache.NewCache("redis",string(redis_conf_js))
	if err!=nil{
		return
	}
	err=bm.Put(key,*data,time.Second*3600)
	return
}


func GetUserData(user *models.User,key ...string)(err error)  {
	o:=orm.NewOrm()
	if len(key)==0{
		err=o.Read(user)
		return
	}

	for _,v:=range key{
		err=o.Read(user,v)
		if err!=nil{
			return
		}
	}
	return
}
func GetUserHousesTable(userId int,houses_list *[]models.House)(err error)   {
	o:=orm.NewOrm()
	qs:=o.QueryTable("house")
	n,err:=qs.RelatedSel("Area").Filter("user_id",userId).All(houses_list)
	if n==0{
		err=nil
	}

	return
}
func GetHouseInfo(house *models.House,id int)(err error){

	o:=orm.NewOrm()
	err=o.QueryTable("House").Filter("Id",id).One(house)
	o.LoadRelated(house,"Area")
	o.LoadRelated(house,"User")
	o.LoadRelated(house,"Images")
	o.LoadRelated(house,"Facilities")
	beego.Info(house.Images)
	beego.Info(house.User)

	return
}
func GetIndexByMysql(Max int,houses *[]models.House)(err error){
	o:=orm.NewOrm()
	_,err=o.QueryTable("House").RelatedSel("Area","User").OrderBy("-Id").Limit(Max).All(houses)


	return
}