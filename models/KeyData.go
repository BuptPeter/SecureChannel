
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type KeyData struct {
	Key_state int `orm:"column(state);null"`
	Key_id  int `orm:"column(id);pk"`
	Key_time time.Time `orm:"column(time);type(datetime)"`
	//Key string `orm:"column(key);size(256);null"`

	//Id   int    `orm:"column(id);pk"`

}

func (t *KeyData) TableName() string {
	return "t_key"
}

func init() {
	orm.RegisterModel(new(KeyData))
}

