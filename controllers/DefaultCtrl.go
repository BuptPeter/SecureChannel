package controllers

import (
	"port-forward/models"

	"github.com/astaxie/beego"
)

type DefaultCtrl struct {
	beego.Controller
}

// @router / [get]
func (c *DefaultCtrl) Default() {

	c.Ctx.Redirect(302, "/login")
}

// @router /apiAuthFail [get]
func (c *DefaultCtrl) ApiAuthFail() {

	c.Data["json"] = models.ResultData{Code: 1, Msg: "ApiAuth鉴权失败"}

	c.ServeJSON()

}
