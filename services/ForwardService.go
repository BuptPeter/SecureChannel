package services

import (
	"io"
	"net"
	"port-forward/models"
	"port-forward/utils"
	"strings"
	"sync"
	"time"
	"fmt"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
)

type ForwardService struct {
}

var (
	portMap            = make(map[string]net.Listener)
	portMapLock        = new(sync.Mutex)
	clientMap          = make(map[string]net.Conn)
	clientMapLock      = new(sync.Mutex)
	sessionId          = 0
	idLock             = new(sync.Mutex)

)

func init() {

}

func (_self *ForwardService) GetNewSessionId() int {
	idLock.Lock()
	defer idLock.Unlock()
	sessionId++

	return sessionId
}

func (_self *ForwardService) PortConflict(key string) bool {
	portMapLock.Lock()
	defer portMapLock.Unlock()

	if _, ok := portMap[key]; ok {
		return true
	} else {
		return false
	}

}

func (_self *ForwardService) RegistryPort(key string, listener net.Listener) {
	portMapLock.Lock()
	defer portMapLock.Unlock()

	portMap[key] = listener

}

func (_self *ForwardService) UnRegistryPort(key string) {
	portMapLock.Lock()
	defer portMapLock.Unlock()

	delete(portMap, key)
	logs.Debug("UnRegistryPort key: ", key)

}

func (_self *ForwardService) RegistryClient(sourcePort string, conn net.Conn) {
	clientMapLock.Lock()
	defer clientMapLock.Unlock()

	clientMap[sourcePort] = conn

}

func (_self *ForwardService) UnRegistryClient(sourcePort string) {
	clientMapLock.Lock()
	defer clientMapLock.Unlock()

	delete(clientMap, sourcePort)
	logs.Debug("UnRegistryClient sourcePort: ", sourcePort)

}


func (_self *ForwardService) GetKeyByEntity(entity *models.PortForward) string {

	fromAddr := fmt.Sprint(entity.Addr, ":", entity.Port)
	toAddr := fmt.Sprint(entity.TargetAddr, ":", entity.TargetPort)
	key := _self.GetKey(fromAddr, toAddr, entity.FType)

	return key
}

func (_self *ForwardService) GetKey(sourcePort, targetPort string, fType int) string {

	return fmt.Sprint(sourcePort, "_", fType, "_TCP_", targetPort)

}

func (_self *ForwardService) StartPortForward(portForward *models.PortForward, result chan models.ResultData) {

	_self.StartPortToPortForward(portForward, result)

}
func (_self *ForwardService) StartTlsPortForward(portForward *models.PortForward, result chan models.ResultData) {

	_self.StartTlsPortToPortForward(portForward, result)

}
func (_self *ForwardService) StartListener(portForward *models.PortForward) {

}

//开启端口转发
//建立普通socket连接
func (_self *ForwardService) StartPortToPortForward(portForward *models.PortForward, result chan models.ResultData) {

	sourcePort := fmt.Sprint(portForward.Addr, ":", portForward.Port)
	targetPort := fmt.Sprint(portForward.TargetAddr, ":", portForward.TargetPort)
	fType := portForward.FType
	End := portForward.End

	var localListener net.Listener

	resultData := &models.ResultData{Code: 0, Msg: ""}
	logs.Debug("StartTcpPortForward sourcePort: ", sourcePort, " targetPort:", targetPort)

	key := _self.GetKey(sourcePort, targetPort, fType)

	if _self.PortConflict(key) {
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("监听地址已被占用 ", sourcePort)
		result <- *resultData
		return
	}
	localListener, err := net.Listen("tcp", sourcePort)

	if err != nil {
		logs.Error("启动监听 ", sourcePort, " 出错：", err)
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("启动监听 ", sourcePort, " 出错：", err)
		result <- *resultData
		return
	}

	_self.RegistryPort(key, localListener)

	result <- *resultData

	go func() {
		for {
			logs.Debug("Ready to Accept ...")
			sourceConn, err := localListener.Accept()

			if err != nil {
				logs.Error("Accept err:", err)
				break
			}

			id := sourceConn.RemoteAddr().String()
			_self.RegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", id), sourceConn)

			logs.Debug("conn.RemoteAddr().String() ：", id)

			targetConn, err := net.DialTimeout("tcp", targetPort, 30*time.Second)
			if err != nil {
				log.Println(err)
				return
			}

			if fType == 0 { //透明转发
				go func() {
					_, err = _self.Copy(targetConn, sourceConn)
					if err != nil {
						logs.Error("客户端来源数据转发到目标端口异常：", err)
						_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
					}
				}()
				go func() {
					_, err = _self.Copy(sourceConn, targetConn)
					if err != nil {
						logs.Error("目标端口返回响应数据异常：", err)
						_self.UnRegistryPort(key)
					}
				}()
			}

			if fType == 1 { //AES加密通信
				if End == 0 { //AES加密通信（OVS端）
					go func() {
						_, err = _self.EncryptCopy(targetConn, sourceConn)
						if err != nil {
							logs.Error("客户端来源数据转发到目标端口异常：", err)
							_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
						}
					}()
					go func() {
						_, err = _self.DecryptCopy(sourceConn, targetConn)
						if err != nil {
							logs.Error("客户端来源数据转发到目标端口异常：", err)
							_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
						}
					}()
				}
				if End == 1 { //AES加密通信（RYU端）
					go func() {
						_, err = _self.DecryptCopy(targetConn, sourceConn)
						if err != nil {
							logs.Error("客户端来源数据转发到目标端口异常：", err)
							_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
						}
					}()
					go func() {
						_, err = _self.EncryptCopy(sourceConn, targetConn)
						if err != nil {
							logs.Error("客户端来源数据转发到目标端口异常：", err)
							_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
						}
					}()
				}
			}

		}
	}()

	logs.Debug("TcpPortForward sourcePort: ", sourcePort, " Close.")

}

//
// sourcePort 源地址和端口，例如：0.0.0.0:8700，本程序会新建立监听
// targetPort 数据转发给哪个端口，例如：192.168.1.100:3306
// 传输过程使用TLS方式
func (_self *ForwardService) StartTlsPortToPortForward(portForward *models.PortForward, result chan models.ResultData) {

	sourcePort := fmt.Sprint(portForward.Addr, ":", portForward.Port)
	targetPort := fmt.Sprint(portForward.TargetAddr, ":", portForward.TargetPort)
	fType := portForward.FType
	End := portForward.End

	var localListener net.Listener
	var targetConn net.Conn

	resultData := &models.ResultData{Code: 0, Msg: ""}
	logs.Debug("StartTcpTLSPortForward sourcePort: ", sourcePort, " targetPort:", targetPort)

	key := _self.GetKey(sourcePort, targetPort, fType)

	if _self.PortConflict(key) {
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("监听地址已被占用 ", sourcePort)
		result <- *resultData
		return
	}
	if End == 1 { //如果是控制器端，开启TLS监听
		var error error
		cert, err := tls.LoadX509KeyPair("data/server.pem", "data/server.key")
		if err != nil {
			resultData.Code = 1
			resultData.Msg = fmt.Sprint("配置TLS出错:", err)
			result <- *resultData
			return
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
		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    clientCertPool,
		}

		localListener, error = tls.Listen("tcp", sourcePort, config)
		if error != nil {
			logs.Error("启动监听 ", sourcePort, " 出错：", error)
			resultData.Code = 1
			resultData.Msg = fmt.Sprint("启动监听 ", sourcePort, " 出错：", error)
			result <- *resultData
			return
		}
		//defer localListener.Close()

	} else { //OVS端则TCP监听
		var error error
		localListener, error = net.Listen("tcp", sourcePort)

		if error != nil {
			logs.Error("启动监听 ", sourcePort, " 出错：", error)
			resultData.Code = 1
			resultData.Msg = fmt.Sprint("启动监听 ", sourcePort, " 出错：", error)
			result <- *resultData
			return
		}
		//defer localListener.Close()
	}
	_self.RegistryPort(key, localListener)

	result <- *resultData

	go func() {
		for {
			logs.Debug("Ready to Accept ...")
			sourceConn, err := localListener.Accept()

			if err != nil {
				logs.Error("Accept err:", err)
				break
			}

			id := sourceConn.RemoteAddr().String()
			_self.RegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", id), sourceConn)

			logs.Debug("conn.RemoteAddr().String() ：", id)

			if End == 0 { //OVS端请求TLS连接（加密通信链路）
				var err_dial error

				cert, err := tls.LoadX509KeyPair("data/client.pem", "data/client.key")
				if err != nil {
					log.Println(err)
					return
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
				targetConn, err_dial = tls.Dial("tcp", targetPort, conf)
				if err_dial != nil {
					log.Println(err_dial)
					return
				}

			} else { //RYU端则请求TCP连接（本地与RYU连接）
				var err_dial error
				targetConn, err_dial = net.DialTimeout("tcp", targetPort, 30*time.Second)
				if err_dial != nil {
					log.Println(err_dial)
					return
				}
			}
				if fType == 0 { //透明转发
					go func() {
						_, err = _self.Copy(targetConn, sourceConn)
						if err != nil {
							logs.Error("客户端来源数据转发到目标端口异常：", err)
							_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
						}
					}()
					go func() {
						_, err = _self.Copy(sourceConn, targetConn)
						if err != nil {
							logs.Error("目标端口返回响应数据异常：", err)
							_self.UnRegistryPort(key)
						}
					}()
				}

				if fType == 1 { //加密通信
					if End == 0 { //加密通信（OVS端）
						go func() {
							_, err = _self.EncryptCopy(targetConn, sourceConn)
							if err != nil {
								logs.Error("客户端来源数据转发到目标端口异常：", err)
								_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
							}
						}()
						go func() {
							_, err = _self.DecryptCopy(sourceConn, targetConn)
							if err != nil {
								logs.Error("客户端来源数据转发到目标端口异常：", err)
								_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
							}
						}()
					}
					if End == 1 { //加密通信（RYU端）
						go func() {
							_, err = _self.DecryptCopy(targetConn, sourceConn)
							if err != nil {
								logs.Error("客户端来源数据转发到目标端口异常：", err)
								_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
							}
						}()

						go func() {
							_, err = _self.EncryptCopy(sourceConn, targetConn)
							if err != nil {
								logs.Error("客户端来源数据转发到目标端口异常：", err)
								_self.UnRegistryClient(fmt.Sprint(sourcePort, "_", fType, "_", sourceConn.RemoteAddr().String()))
							}
						}()
					}
				}
		}
	}()

	logs.Debug("TcpPortForward sourcePort: ", sourcePort, " Close.")

}

func (_self *ForwardService) ClosePortForward(sourcePort string, targetPort string, fType int, result chan models.ResultData) {
	resultData := &models.ResultData{Code: 0, Msg: ""}

	logs.Debug("CloseTcpPortForward:", sourcePort)
	//先关闭客户端连接
	for cId, conn := range clientMap {
		//logs.Debug("clientMap id：", cId)
		if strings.HasPrefix(cId, fmt.Sprint(sourcePort, "_", fType)) {
			logs.Debug("close clientMap id：", cId)
			if conn != nil {
				conn.Close()
			}
			_self.UnRegistryClient(cId)
		}

	}

	//关闭本地监听
	key := _self.GetKey(sourcePort, targetPort, fType)
	if localListener, ok := portMap[key]; ok {
		if localListener != nil {
			localListener.Close()
			logs.Debug("listener close:", key)
		}

		_self.UnRegistryPort(key)
	} else {
		resultData.Code = 1
		resultData.Msg = fmt.Sprint("未启用监听 ", key)

	}

	result <- *resultData

	logs.Debug("CloseTcpPortForward finished.")

}

func (_self *ForwardService) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
func (_self *ForwardService) EncryptCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	temp_reader := make([]byte, 32*1024)

	//var start_time,stop_time time.Time
	for {
		nr, er := src.Read(temp_reader)
		if nr > 0 {
			//logs.Info("接收明文消息：", temp_reader[:n])
			//start_time := time.Now()
			nw, ew := dst.Write(utils.GetCipherText(temp_reader[0:nr], "astaxie12798akljzmknm.ahkjkljl;k"))
			//stop_time := time.Now()
			//logs.Info("转发密文消息 OVS -> SecChan ：", temp_reader[:n])
			//logs.Info("消息长度变化：", nr, "->", nw, " ，加密耗时：", stop_time.UnixNano()-start_time.UnixNano())

			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("short write")
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func (_self *ForwardService) DecryptCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	temp_reader := make([]byte, 32*1024)

	//var start_time,stop_time time.Time
	for {
		nr, er := src.Read(temp_reader)
		if nr > 0 {
			//logs.Info("接收明文消息：", temp_reader[:n])
			//start_time := time.Now()
			nw, ew := dst.Write(utils.GetPlainText(temp_reader[0:nr], "astaxie12798akljzmknm.ahkjkljl;k"))
			//stop_time := time.Now()
			//logs.Info("转发密文消息 OVS -> SecChan ：", temp_reader[:n])
			//logs.Info("消息长度变化：", nr, "->", nw, " ，解密耗时：", stop_time.UnixNano()-start_time.UnixNano())

			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("short write")
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
