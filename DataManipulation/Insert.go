package DataManipulation

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"uhome/UHomeWeb/models"
)

func InsertHouse(house *models.House)(houseId int,err error) {
	o:=orm.NewOrm()
	_,err=o.Insert(house)
	houseId=house.Id
	return
}
func M2MInsertHouse(house *models.House,facility *[]*models.Facility)(err error)  {
	o:=orm.NewOrm()
	m2m:=o.QueryM2M(house,"Facilities")
	_,err=m2m.Add(*facility)
	return
}