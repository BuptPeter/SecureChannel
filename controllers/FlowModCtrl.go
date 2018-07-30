package controllers

import (
	"port-forward/controllers/base"
	"port-forward/models"
	"port-forward/services"
	"github.com/astaxie/beego/logs"
)

type FlowModCtrl struct {
	BaseController.ConsoleController
}

// @router /u/FlowModList/json [post]
func (c *FlowModCtrl) FlowModListJson() {
	logs.Info("start FlowModListJson...")
	pageParam := new(models.PageParam)
	pageParam.PIndex, _ = c.GetInt64("pIndex")
	pageParam.PSize, _ = c.GetInt64("pSize")
	//port, _ := c.GetInt("port")

	query := &models.FlowData{}
	//query.Port = port
	//query.FType = -1

	pageData := services.SysDataS.GetFlowModList(query, pageParam.PIndex, pageParam.PSize)

	//for _, entity := range pageData.Data.([]*models.FlowData) {
	//	key := services.ForwardS.GetKeyByEntity(entity)
	//	entity.Status = utils.If(services.ForwardS.PortConflict(key), 1, 0).(int)
	//}
	logs.Info(pageData.TotalRows)
	state := services.FlowModS.GetHashServiceState()
	c.Data["json"] = models.ResultData{Code: 0, Msg: "success", State:state,Data: pageData}

	c.ServeJSON()

}
// @router /u/FlowCheckList/json [post]
func (c *FlowModCtrl) FlowCheckListJson() {
	logs.Info("start FlowCheckListJson...")
	pageParam := new(models.PageParam)
	pageParam.PIndex, _ = c.GetInt64("pIndex")
	pageParam.PSize, _ = c.GetInt64("pSize")

	query := &models.FlowCheckData{}


	pageData := services.SysDataS.GetFlowCheckList(query, pageParam.PIndex, pageParam.PSize)

	logs.Info(pageData.TotalRows)
	state := services.FlowModS.GetCheckServiceState()
	c.Data["json"] = models.ResultData{Code: 0, Msg: "success", State:state,Data: pageData}

	c.ServeJSON()

}

// @router /u/OpenFlowMod [get,post]
func (c *ForwardCtrl) OpenFlowModService() {
	resultChan := make(chan models.ResultData)

	logs.Debug("Start OpenFlowModService...  ")
	go services.FlowModS.StartHashService(resultChan)

	c.Data["json"] = <-resultChan

	c.ServeJSON()

}

// @router /u/OpenFlowCheck [get,post]
func (c *ForwardCtrl) OpenFlowCheckService() {
	resultChan := make(chan models.ResultData)

	logs.Debug("Start OpenFlowCheckService...  ")
	go services.FlowModS.StartCheckService(resultChan)

	c.Data["json"] = <-resultChan

	c.ServeJSON()

}