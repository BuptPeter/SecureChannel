package main

import (
	"time"
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
	"crypto/x509"
	"port-forward/utils"
)

func TestTime(){
	start_time := time.Now().Unix()
	fmt.Println(start_time)
	time.Sleep(2*time.Second)
	stop_time := time.Now().Unix()
	fmt.Println(stop_time)
	fmt.Println((stop_time-start_time))
}
func  StartTLSClient(targetPort string) (conn tls.Conn,err error) {

	cert, err := tls.LoadX509KeyPair("data/client.pem", "data/client.key")
	if err != nil {
		logs.Debug(err)
		return tls.Conn{},err
	}
	certBytes, err := ioutil.ReadFile("data/client.pem")
	if err != nil {
		logs.Debug("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		logs.Debug("failed to parse root certificate")
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	targetConn,err := tls.Dial("tcp", targetPort, conf)
	if err != nil {
		logs.Debug("failed to Dial up :",err)

		return tls.Conn{},err
	}
	defer targetConn.Close()
	return *targetConn,nil
	}

func main(){
	addr := "127.0.0.1"
	port := "8778"

	tslConn,_ := StartTLSClient((addr+":"+port))
	temp_b := make([]byte, 1000)
	utils.GetCipherText(temp_b,"astaxie12798akljzmknm.ahkjkljl;k")
	fmt.Println("连接建立，开始测试。")
	for i:=0;i>10;i++{
		n,err:=tslConn.Write(temp_b)
		if err !=nil{
			fmt.Println("写入数据失败：",err)
			break
		}
		fmt.Println("写入：",n,i)

	}

}