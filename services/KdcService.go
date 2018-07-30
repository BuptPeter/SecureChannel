package services

import (
	"net"
	"encoding/base64"
	"crypto/aes"
	"errors"
	"crypto/cipher"
	"port-forward/models"
	"github.com/astaxie/beego/logs"
	"fmt"
	"time"
	mrand "math/rand"
	"io"
	"crypto/rand"
	"strconv"
)
func ClietnKDC(id int ,id_ovs string,id_con string,kdc string, result chan models.ResultData)  {
	resultData := &models.ResultData{Code: 0, Msg: ""}
	KDC_Port := "8848"
    //打开连接:
    conn, err := net.Dial("tcp", kdc+":"+KDC_Port)
    if err != nil {
        //由于目标计算机积极拒绝而无法创建连接
        logs.Error("Error dialing ", err)
        resultData.Code = 1
        logs.Error("Error dialing", err.Error())
        resultData.Msg = fmt.Sprint("Error dialing", err.Error())
    	result <- *resultData
        return // 终止程序
    }
    logs.Info("连接成功")
    logs.Info("发送ID：",id_ovs,id_con)
    _, err = conn.Write([]byte(id_ovs + "," + id_con))

    buf := make([]byte, 512)
    length, err := conn.Read(buf)
    if err != nil {
    	logs.Error("Error reading ", err)
    	resultData.Code = 1
    	resultData.Msg = fmt.Sprint("Error reading", err.Error())
    	result <- *resultData
            return //终止程序
        }
	//logs.Info("获取到票据")
    if length >0{
    	s:=string(buf[0:length])
    	logs.Info("收到密文票据：",len(s),s)
    	s_key,t_con,t_ovs:= GetSessionKey_OVS(s)

    	logs.Info(s_key,t_con)

		entity:=GetPortForwardById(id)
		entity.Key=s_key
		entity.Key_state=1
		entity.Key_time= time.Now().String()
		entity.Con_ticket=t_con
		entity.Ovs_ticket=t_ovs
		logs.Info("t_OVS :",t_ovs)
    	SavePortForward(entity)
	}
	conn.Close()
    result <- *resultData

}
func StartService(id int ,port string, result chan models.ResultData){
	resultData := &models.ResultData{Code: 0, Msg: ""}

    logs.Info("Starting the server ...")
    // 创建 listener
    listener, err := net.Listen("tcp", "0.0.0.0:"+port)
    if err != nil {
        logs.Error("Error listening", err.Error())
            	resultData.Code = 1
    	resultData.Msg = fmt.Sprint("Error listening", err.Error())
    	result <- *resultData
        return //终止程序
    }
    go KeyReciver(id,listener)

    //更新密钥分配状态和时间
    entity:=GetPortForwardById(id)
    entity.Key_state=1
    entity.Key_time= time.Now().String()
    SavePortForward(entity)

	result <- *resultData

}

func KeyReciver(id int ,listener net.Listener){
	    // 监听并接受来自OVS端的连接
		defer listener.Close()
        conn, err := listener.Accept()
        if err != nil {
            logs.Error("Error accepting.", err.Error())
            return // 终止程序
        }

        buf := make([]byte, 512)
        length, err := conn.Read(buf)
        if err != nil {
            logs.Error("Error reading.", err.Error())
            return //终止程序
        }
        s:=string(buf[0:length])
        logs.Info("Connected by:",conn.RemoteAddr().String(),"\r\n\tReceived data: ",len(s),s)

        s_key,TS,B,t_con:=GetSessionKey_Con(s)

        logs.Info("检验时间戳...",TS,s_key,B)
        if CheckTime(TS){
        	logs.Info("时间戳检查通过...",)
        	entity:=GetPortForwardById(id)
			entity.Key=s_key
			entity.Key_state=2
			entity.Key_time= time.Now().String()
			entity.Con_ticket=t_con

			SavePortForward(entity)

			logs.Info("发送认证因子B,关闭控制器端连接.")
			conn.Write([]byte(B))
			conn.Close()
			logs.Info("控制器端接受服务连接关闭...")

		}else {
       		entity:=GetPortForwardById(id)
			//entity.Key=s_key
			entity.Key_state=0
			entity.Key_time= time.Now().String()
			SavePortForward(entity)
			logs.Info("Error TS..")
			conn.Close()
            return //终止程序
		}
}
func CheckTime(TS string )bool{
	//没有判断是否第一次出现
	//string->int64
	t0,_:=strconv.ParseInt(TS, 10, 64)
	if (time.Now().Unix()-t0>5*60){//超过5min,失效。
		return false
	}
	return true
}
func CreatA(n int,skey string)(string,string){
		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		//b[i] = letters[rand.Intn(len(letters))]
		b[i] = letters[mrand.Intn(len(letters))]
	}
	TS:=strconv.FormatInt(time.Now().Unix(),10)
	A:=TS+string(b)
	logs.Info("A: ",len(A),A)

	A_,_:=Encrypt([]byte(skey),[]byte(A))

	logs.Info("A_: ",len(A_),A_)
	return A_,TS

}
func ClientCon(id int,con_ip string,port string,T string,skey string, result chan models.ResultData){
	resultData := &models.ResultData{Code: 0, Msg: ""}

    //打开连接:
    conn, err := net.Dial("tcp",con_ip+ ":"+port)
    if err != nil {
        //由于目标计算机积极拒绝而无法创建连接
        logs.Error("Error dialing ", err)
        resultData.Code = 1
        logs.Error("Error dialing", err.Error())
        resultData.Msg = fmt.Sprint("Error dialing", err.Error())
    	result <- *resultData
        return // 终止程序
    }
    logs.Info("连接成功")
    A,TS:=CreatA(20,skey)
    temp:=T+A
    //发送票据和认证因子A
    _, err = conn.Write([]byte(temp))
    if err != nil {
        //发送数据失败
        logs.Error("Error sending ", err)
        resultData.Code = 1
        logs.Error("Error sending", err.Error())
        resultData.Msg = fmt.Sprint("Error sending", err.Error())
    	result <- *resultData
        return // 终止程序
    }
    logs.Info("票据+A发送成功",len(temp),temp)
//---------------------------------
    buf := make([]byte, 512)
    length, err := conn.Read(buf)
    if err != nil {
    	logs.Error("Error reading ", err)
    	resultData.Code = 1
    	resultData.Msg = fmt.Sprint("Error reading", err.Error())
    	result <- *resultData
            return //终止程序
        }
    if length >0{
    	s:=string(buf[0:length])
    	logs.Info("收到认证因子B：",len(s),s)

    	TS_recive,_:=Decrypt([]byte(skey),s)

    	if (TS==string(TS_recive)){

    	logs.Info("时间戳正确,完成身份认证及密钥分配",TS,string(TS_recive))

		entity:=GetPortForwardById(id)
		entity.Key_state=2
		entity.Key_time= time.Now().String()
    	SavePortForward(entity)
		}
	}
	conn.Close()
    logs.Info("OVS端连接关闭...")
    result <- *resultData

}

func GetSessionKey_OVS(tickets string) (string,string,string){
	length:=len(tickets)
	ticket_ovs:=tickets[0:length/2]
	ticket_con:=tickets[length/2:length]
	logs.Info(len(ticket_ovs),ticket_ovs)
	logs.Info(len(ticket_con),ticket_con)
	s_key_byte,_ := Decrypt(GetPreKey(),ticket_ovs)
	s_key:=string(s_key_byte)[0:32]
	logs.Info("解密票据,得到会话密钥：",len(s_key),string(s_key))
	return s_key,ticket_con,ticket_ovs
}

func GetSessionKey_Con(tickets string) (string,string,string,string){
	length:=len(tickets)
	ticket_con:=tickets[0:95]
	A:=tickets[95:length]
	logs.Info("Ticekt_Con :",len(ticket_con),ticket_con)
	logs.Info("A :",len(A),A)
	s_key_byte,_ := Decrypt(GetPreKey(),ticket_con)
	s_key:=string(s_key_byte)[0:32]
	logs.Info("解密票据,得到会话密钥：",len(s_key),string(s_key))

	TimeS_byte,_:=Decrypt([]byte(s_key),string(A))
	TimeS:=string(TimeS_byte)[0:len(TimeS_byte)-20]
	logs.Info("得到时间戳：",TimeS)

	B,_:=Encrypt([]byte(s_key),[]byte(TimeS))
	logs.Info("得到双向认证因子B:",B)

	return s_key,TimeS,B,ticket_con

}
//Decrypt 解密算法
func Decrypt(key []byte, encrypted string) ([]byte, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
//Encrypt 加密算法
func Encrypt(key []byte, data []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}
func SavePortForward(entity *models.PortForward) error {

	if entity.Id > 0 {
		update := &models.PortForward{}

		qs := OrmerS.QueryTable(new(models.PortForward))
		qs = qs.Filter("Id", entity.Id)
		err := qs.One(update)
		if err != nil {
			//如果没查到数据，会抛出 no row found
			return err
		}

		update.Name = entity.Name
		update.Addr = entity.Addr
		update.Port = entity.Port
		//update.Protocol = entity.Protocol
		update.TargetAddr = entity.TargetAddr
		update.TargetPort = entity.TargetPort
		//update.Others = entity.Others
		update.FType = entity.FType
		update.Tls = entity.Tls
		update.End = entity.End


		update.Key  = entity.Key
		update.Key_state  = entity.Key_state
		update.Key_id   = entity.Key_id
		update.Key_time  = entity.Key_time
		update.Kdc  = entity.Kdc
		update.Con_ticket  = entity.Con_ticket
		update.Ovs_ticket  = entity.Ovs_ticket
		_, err1 := OrmerS.Update(update)
		return err1
	} else {
		entity.CreateTime = time.Now()
		res, err := OrmerS.Raw("INSERT INTO t_port_forward(name, status, addr, port, targetAddr, targetPort, createTime, fType,tls,end) values(?,?,?,?,?,?,?,?,?,?)",
			entity.Name, entity.Status, entity.Addr, entity.Port, entity.TargetAddr, entity.TargetPort, entity.CreateTime, entity.FType, entity.Tls, entity.End).Exec()
		if err == nil {
			num, _ := res.RowsAffected()
			logs.Debug("AddPortForward", num)

		} else {
			logs.Error("AddPortForward", err)
		}
		return err
	}
}

func  GetPortForwardById(id int) *models.PortForward {

	entity := new(models.PortForward)
	qs := OrmerS.QueryTable(entity)

	qs = qs.Filter("Id", id)

	err := qs.One(entity)

	if err != nil {
		logs.Error("GetPortForwardById ", err)
		return nil
	}

	return entity

}
func GetPreKey()[]byte{
		prekey:="astaxie12798akljzmknm.ahkjkljl;k"
    return []byte(prekey)

}