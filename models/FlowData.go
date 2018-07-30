package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type FlowData struct {
	Id  int `orm:"column(id);pk;null"`
	Flow_time time.Time `orm:"column(time);type(datetime)"`
	Flow string `orm:"column(flow);size(4096);null"`
	Flow_mac string `orm:"column(mac);size(256);null"`
	Flow_mod string `orm:"column(flowmod);size(512);null"`


}

func (t *FlowData) TableName() string {
	return "t_Flow"
}

func init() {
	orm.RegisterModel(new(FlowData))
}

