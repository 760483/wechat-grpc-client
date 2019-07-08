package client

import (
	"net"
	"wechat-client/client/system"
	pb "wechat-client/proto"
	"errors"
	"time"
	"fmt"
)

type LoginParam struct {
	LongHead   []byte
	MaxSyncKey []byte
	UUID       []byte
}

type WechatConn struct {
	WxUser        string
	ShortHost     string
	LongHost      string
	Client        net.Conn
	LoginParam    LoginParam
	VUser         *pb.User
	MsgSyncKey    []byte
	MsgMaxSyncKey []byte
	MsgLock       bool
}

//连接微信服务器，缓存tcp连接于wxUser中
func ConnectWx(wxUser string, host string) error {
	address, err := net.ResolveTCPAddr("tcp", host+":443")
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "微信tcp客户端创建失败，error: "+err.Error())
		return errors.New("微信tcp客户端创建失败")
	}
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "微信tcp客户端连接微信服务器失败，error: "+err.Error())
		return errors.New("微信tcp客户端连接微信服务器失败")
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.SetCache(wxUser, &WechatConn{Client: conn, WxUser: wxUser}, 30*24*time.Hour)
	} else {
		wxClient.Client = conn
		SetWechatConn(wxClient)
	}
	go client_system.AsyncRecvWxData(wxUser, conn)
	return nil
}

//获取微信客户端缓存
func GetWechatConn(wxUser string) (*WechatConn, error) {
	conn, has := client_system.GetCache(wxUser)
	if !has {
		return nil, errors.New("获取连接缓存失败")
	}
	return conn.(*WechatConn), nil
}

//重新设置微信客户端缓存
func SetWechatConn(conn *WechatConn) {
	client_system.SetCache(conn.WxUser, conn, 30*24*time.Hour)
}

//向微信服务器发送心跳
func SendHeartbeat(wxUser string) {
	requestProto := client_system.CreateWechatMsg(205, []byte{})
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return
	}
	//如果微信长连接断开，那么直接进行断线重连的操作
	defer ReConnect(wxUser)
	requestProto.BaseMsg.User.DeviceId = wxClient.VUser.DeviceId
	for {
		time.Sleep(1 * time.Minute)
		client_system.LogWrite(client_system.LOG_INFO, "开始发送心跳")
		_, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
			return nil, nil
		})
		if err != nil {
			client_system.LogWrite(client_system.LOG_ERROR, "微信服务器心跳发送失败了。退出发送心跳协程，error: "+err.Error())
			break
		}
	}
}

//断线自动重连，主要是用于发送心跳失败，sessionTimeout出现的时候来进行重连操作
func ReConnect(wxUser string) {
	client_system.LogWrite(client_system.LOG_ERROR, wxUser+" 微信连接客户端已断开，正在尝试重新连接...")
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return
	}
	ConnectWx(wxUser, wxClient.LongHost)
	_, err = AutoLogin(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "微信断线重连失败了", err.Error()))
		return
	}
	client_system.LogWrite(client_system.LOG_ERROR, "微信断线重连成功...")
	go SendHeartbeat(wxUser)
}
