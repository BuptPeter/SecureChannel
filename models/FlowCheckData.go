package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type FlowCheckData struct {
	Id  int `orm:"column(id);pk;null"`
	Time time.Time `orm:"column(time);type(datetime)"`
	Flow string `orm:"column(flows);size(4096);null"`
	Flow_bingo string `orm:"column(flows_bingo);size(4096);null"`
	Flow_mac_bingo string `orm:"column(mac_bingo);size(256);null"`
	Flow_mac string `orm:"column(mac);size(256);null"`
	State int `orm:"column(state);null"`


}

func (t *FlowCheckData) TableName() string {
	return "t_flow_check"
}

func init() {
	orm.RegisterModel(new(FlowCheckData))
}
