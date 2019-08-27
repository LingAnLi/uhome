package DataManipulation

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"uhome/UHomeWeb/models"
)

func UserUpData(user*models.User,key string)  (err error){
	o:=orm.NewOrm()
	if key==""{
		_,err=o.Update(user)
		return
	}
	_,err=o.Update(user,key)
	return
}
