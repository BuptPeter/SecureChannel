package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["port-forward/controllers:DefaultCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:DefaultCtrl"],
		beego.ControllerComments{
			Method: "Default",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:DefaultCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:DefaultCtrl"],
		beego.ControllerComments{
			Method: "ApiAuthFail",
			Router: `/apiAuthFail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:FlowModCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:FlowModCtrl"],
		beego.ControllerComments{
			Method: "FlowCheckListJson",
			Router: `/u/FlowCheckList/json`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:FlowModCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:FlowModCtrl"],
		beego.ControllerComments{
			Method: "FlowModListJson",
			Router: `/u/FlowModList/json`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "AddForward",
			Router: `/u/AddForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "ClearNetAgentStatus",
			Router: `/u/ClearNetAgentStatus`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "ClientController",
			Router: `/u/ClientController`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "ClientKDC",
			Router: `/u/ClientKDC`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "CloseForward",
			Router: `/u/CloseForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "CloseMagicService",
			Router: `/u/CloseMagicService`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "DelForward",
			Router: `/u/DelForward`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "EditForward",
			Router: `/u/EditForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "ForwardList",
			Router: `/u/ForwardList`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "ForwardListJson",
			Router: `/u/ForwardList/json`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "GetMagicStatus",
			Router: `/u/GetMagicStatus`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "GetNetAgentStatus",
			Router: `/u/GetNetAgentStatus`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "Kerberos",
			Router: `/u/Kerberos`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "NetAgent",
			Router: `/u/NetAgent`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OpenFlowCheckService",
			Router: `/u/OpenFlowCheck`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OpenFlowModService",
			Router: `/u/OpenFlowMod`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OpenForward",
			Router: `/u/OpenForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OpenMagicService",
			Router: `/u/OpenMagicService`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OvsCheck",
			Router: `/u/OvsCheck`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "OvsHash",
			Router: `/u/OvsHash`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "SaveForward",
			Router: `/u/SaveForward`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:ForwardCtrl"],
		beego.ControllerComments{
			Method: "StartService",
			Router: `/u/StartService`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:HelpCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:HelpCtrl"],
		beego.ControllerComments{
			Method: "GetTcp",
			Router: `/GetTcp`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"],
		beego.ControllerComments{
			Method: "DoLogin",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:LoginCtrl"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"],
		beego.ControllerComments{
			Method: "CloseForward",
			Router: `/CloseForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"],
		beego.ControllerComments{
			Method: "OpenForward",
			Router: `/OpenForward`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:RestApiCtrl"],
		beego.ControllerComments{
			Method: "ServerSummary",
			Router: `/ServerSummary`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"],
		beego.ControllerComments{
			Method: "ChangePwd",
			Router: `/u/changePwd`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"],
		beego.ControllerComments{
			Method: "DoChangePwd",
			Router: `/u/doChangePwd`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"],
		beego.ControllerComments{
			Method: "GetServerTime",
			Router: `/u/getServerTime`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"] = append(beego.GlobalControllerRouter["port-forward/controllers:UCenterCtrl"],
		beego.ControllerComments{
			Method: "Main",
			Router: `/u/main`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
