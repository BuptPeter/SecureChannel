package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type PortForward struct {
	Key_state int `orm:"column(key_state);null"`
	Key_id  int `orm:"column(key_id);null"`
	Key_time string `orm:"column(key_time);size(256);null"`
	Kdc string `orm:"column(kdc);size(256);null"`
	Con_ticket string `orm:"column(con_ticket);size(512);null"`
	Ovs_ticket string `orm:"column(ovs_ticket);size(512);null"`
	Key string `orm:"column(key);size(256);null"`

	Id   int    `orm:"column(id);pk"`
	Name string `orm:"column(name);size(256);null"`
	// 0:禁用,1:启用
	Status int    `orm:"column(status);null"`
	Addr   string `orm:"column(addr);size(256);null"`
	// 端口号
	Port int `orm:"column(port);null"`

	TargetAddr string `orm:"column(targetAddr);size(256);null"`
	// 端口号
	TargetPort int       `orm:"column(targetPort);null"`
	CreateTime time.Time `orm:"column(createTime);type(datetime)"`
	// 0:普通映射,1:加密通信（OVS端）,2:加密通信（控制器端）
	FType int `orm:"column(fType);null"`
	//吞吐测试标志位
	Tls int `orm:"column(tls);null"`
	//工作在哪一端
	End int `orm:"column(end);null"`
}

func (t *PortForward) TableName() string {
	return "t_port_forward"
}

func init() {
	orm.RegisterModel(new(PortForward))
}
