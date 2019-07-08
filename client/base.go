package client

import (
	"wechat-client/client/system"
	pb "wechat-client/proto"
	"errors"
)

//快速请求短连接，会解析值
func FastShortRequest(wxUser string, cmd int32, payload []byte, logMsg string, out *[]byte, err *error) {
	wxClient, err1 := GetWechatConn(wxUser)
	errStr := `【` + wxUser + `】` + logMsg + `，error: `
	if err1 != nil {
		errStr += "获取微信实例缓存信息失败," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	requestProto := client_system.CreateWechatMsg(cmd, payload)
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err1 := client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err1 != nil {
		errStr += "grpc通讯异常," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	if ret.GetBaseMsg().GetRet() != 0 {
		errStr += "微信响应异常," + string(ret.GetBaseMsg().GetPayloads())
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	*out = ret.GetBaseMsg().GetPayloads()
	return
}

//快速请求长连接，快速解析值
func FastLongRequest(wxUser string, cmd int32, payload []byte, logMsg string, out *[]byte, err *error) {
	wxClient, err1 := GetWechatConn(wxUser)
	errStr := `【` + wxUser + `】` + logMsg + `，error: `
	if err1 != nil {
		errStr += "获取微信实例缓存信息失败," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	requestProto := client_system.CreateWechatMsg(cmd, payload)
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err1 := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err1 != nil {
		errStr += "grpc通讯异常," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	if ret.GetBaseMsg().GetRet() != 0 {
		errStr += "微信响应异常," + string(ret.GetBaseMsg().GetPayloads())
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	*out = ret.GetBaseMsg().GetPayloads()
	return
}

//快速请求长连接操作，不会解析组包的值
func FastLongRequestOperate(wxUser string, cmd int32, payload []byte, logMsg string, err *error) {
	wxClient, err1 := GetWechatConn(wxUser)
	errStr := `【` + wxUser + `】` + logMsg + `，error: `
	if err1 != nil {
		errStr += "获取微信实例缓存信息失败," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	requestProto := client_system.CreateWechatMsg(cmd, payload)
	requestProto.BaseMsg.User = wxClient.VUser
	_, err1 = client_system.LongRequest(wxClient.Client, requestProto, nil)
	if err1 != nil {
		errStr += "grpc通讯异常," + err1.Error()
		client_system.LogWrite(client_system.LOG_ERROR, errStr)
		*err = errors.New(errStr)
		return
	}
	return
}
