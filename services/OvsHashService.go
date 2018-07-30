package services

import (
	"net"
	"github.com/astaxie/beego/logs"
	"os/exec"
	"fmt"
	"io/ioutil"
	"time"
	"port-forward/models"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"strings"
)


type OvsHashService struct {
}
var HashSeviceAlived bool =false
var CheckSeviceAlived bool =false


func (_self *OvsHashService)GetHashServiceState()bool{

	return HashSeviceAlived
}
func (_self *OvsHashService)GetCheckServiceState()bool{

	return CheckSeviceAlived
}

func ReleaseHashSeviceAlived(){
	HashSeviceAlived = false
	}
func ReleaseCheckSeviceAlived(){
	HashSeviceAlived = false
	}
func getHmacCode(s string) string {
    h := hmac.New(sha256.New, []byte("c3e032aee35128ed92fdcf70253f1d661e"))
    io.WriteString(h, s)
    return fmt.Sprintf("%x", h.Sum(nil))
}

func (_self *OvsHashService)StartCheckService(result chan models.ResultData){
	resultData := &models.ResultData{Code: 0, Msg: ""}
	if CheckSeviceAlived==true{
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("ovs-check服务已启动。 ")
		result <- *resultData
		return
	}

	CheckSeviceAlived=true
	defer ReleaseCheckSeviceAlived()
	logs.Info("StartCheckService success.")
	result <- *resultData
	for{//主循环定时检查

		logs.Info("开始定时校检...")
		flows:=GetFlows()
		mac:=getHmacCode(string(flows))
		flows_bingo,mac_bingo:= GetLastFlowMod()
		_,err:=_SaveFlowCheck(string(flows), mac, mac_bingo,flows_bingo)
		if err!=nil{
			logs.Error("_SaveFlowCheck error:",err)
		}
		logs.Info("_SaveFlowCheck success ..")
		time.Sleep(600*1000*1000*1000)//60s
	}

}
func  GetLastFlowMod() (string,string)  {
	//res, err := OrmerS.Raw("select * from t_flow_check order by id desc limit 0,1;").Exec()
	var entites []*models.FlowData
	entity := new(models.FlowData)
	qs := OrmerS.QueryTable(entity)
	_, err := qs.Limit(1).OrderBy("-Id").All(&entites)
	if err == nil {
		logs.Debug("select lasted flow mod success ", entites[0].Flow_mac)
		return entites[0].Flow,string(entites[0].Flow_mac)
		} else {
			logs.Error("select lasted flow mod failed ", err)
			return"",""
		}
}
func (_self *OvsHashService)StartHashService(result chan models.ResultData){
	resultData := &models.ResultData{Code: 0, Msg: ""}
	if HashSeviceAlived==true{
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("ovs-hash服务已启动。 ")
		result <- *resultData
		return
	}

	HashSeviceAlived=true
	defer ReleaseHashSeviceAlived()

	l, err := net.Listen("tcp", ":11777")
    if err != nil {
        logs.Info("listen error:", err)
        resultData.Code = 1
		resultData.Msg = fmt.Sprint("listen error:", err)
		result <- *resultData
        return
    }
    result <- *resultData
    for {
        context, err := l.Accept()
        if err != nil {
            logs.Info("accept error:", err)
            resultData.Code = 1
			resultData.Msg = fmt.Sprint("accept error:", err)
			result <- *resultData
            break
        }
        // start a new goroutine to handle
        // the new connection.
        go handleConn(context)
    }
}
func handleConn(con net.Conn){
	defer handleError(con)
	var buf = make([]byte, 1024)
        // read from the connection
        logs.Info("start to read flow_mod from ovs")
	for{
		n, err := con.Read(buf)
        if err != nil {
            logs.Info("conn read error:", err)
            return
        }
        logs.Info("read %d bytes, content : %s\n", n, string(buf[:n]))
        folwmod:=buf[:n]
        n,err=con.Write([]byte("ok!"))
        if err != nil {
            logs.Info("conn write error:", err)
            return
        }
		flows:=GetFlows()
		logs.Info("GetFlows success,flows:",string(flows))
		if flows==nil{
			logs.Info("GetFlows error.")
            return
		}
        logs.Info("send to ovs ok..")
        go SaveFlowMod(folwmod,flows)
	}

}
func GetFlows() []byte{
        cmd := exec.Command("/bin/bash", "-c", `ovs-ofctl dump-flows br0`)
        //cmd := exec.Command("cmd", "/c", `ipconfig`)
        //创建获取命令输出管道
        stdout, err := cmd.StdoutPipe()
        if err != nil {
            fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
            return nil
        }
        //执行命令
        if err := cmd.Start(); err != nil {
			fmt.Println("Error:The command is err,", err)
			return nil
		}
        //读取所有输出
        bytes, err := ioutil.ReadAll(stdout)
        if err != nil {
            fmt.Println("ReadAll Stdout:", err.Error())
            return nil
        }
        if err := cmd.Wait(); err != nil {
            fmt.Println("wait:", err.Error())
            return nil
        }
        flows:=strings.Split(string(bytes),"cookie=")
//      过滤非流表字段
//      n_packets、n_bytes，匹配到这条规则的网络包数、字节数。
//		idle_age：多久没有数据包经过这条规则，单位秒
//      hard_age：距这条规则创建、修改经过的时间，单位秒
		infoList:=[]string{"duration","n_packets","n_bytes","idle_age","hard_age"}
		for j:=1;j<len(flows) ;j++  {
			flows[j] = strings.Replace(flows[j]," actions",", actions",-1)
			for i:=0;i<len(infoList) ;i++  {
				flows[j] = removeInfo(flows[j],infoList[i])
			}
			flows[j] = "cookie="+flows[j]
		}
		_flows:=strings.Join(flows[1:],"")
		logs.Info("Get Flows:",_flows)


        return []byte(_flows)
    }

func removeInfo(str string ,substr string)string{
	start:=strings.Index(str,substr)
	if(start==-1){
		return str
	}
	end:=start+strings.Index(str[start:],",")
	//logs.Info("start:",start,"end:",end,"f[start:end]:",string(str[start:end]))
	f1:=strings.Replace(str,string(str[start-1:end+1]),"",-1)
	return f1
}
func handleError(c net.Conn){
    	c.Write([]byte("error!"))
    	c.Close()
	}

func SaveFlowMod(flowmod []byte ,flows []byte){

	//mac:=utils.GetHmac(flows,[]byte("astaxie12798akljzmknm.ahkjkljl;k"))
	mac:=getHmacCode(string(flows))
	logs.Info("Get mac :",string(mac))

	_,err:=_SaveFlowMod(string(flows),string(mac),string(flowmod))
	if err!=nil{
		logs.Error("_SaveFlowMod error:",err)
	}
	logs.Info("SaveFlowMod success ..")
}

func  _SaveFlowMod(flows string, mac string, flowmod string) (*models.FlowData, error) {


	entity := &models.FlowData{}
	//entity.Id =0
	entity.Flow = flows
	entity.Flow_mac = mac
	entity.Flow_mod = flowmod
	entity.Flow_time = time.Now()

	res, err := OrmerS.Raw("INSERT INTO t_flow(flow, mac, time, flowmod) values(?,?,?,?)",
			entity.Flow, entity.Flow_mac, entity.Flow_time, entity.Flow_mod).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		logs.Debug("Add flow mod success ", num)
		} else {
			logs.Error("Add flow mod failed ", err)

		}
	return entity, err
}

func  _SaveFlowCheck(flows string, mac string, mac_bingo string,flows_bingo string) (*models.FlowCheckData, error) {

	entity := &models.FlowCheckData{}
	entity.Flow = flows
	entity.Flow_bingo =flows_bingo
	entity.Flow_mac = mac
	entity.Flow_mac_bingo = mac_bingo
	entity.State=strings.Compare(mac,mac_bingo)
	entity.Time = time.Now()

	res, err := OrmerS.Raw("INSERT INTO t_flow_check(flows,flows_bingo, mac, time, mac_bingo,state) values(?,?,?,?,?,?)",
			entity.Flow,entity.Flow_bingo,entity.Flow_mac, entity.Time, entity.Flow_mac_bingo,entity.State).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		logs.Debug("Add flow check success ", num)
		} else {
			logs.Error("Add flow check failed ", err)

		}
	return entity, err
}