package client_system

import (
	pb "wechat-client/proto"
	"net/http"
	"bytes"
	"io/ioutil"
	"strings"
	"time"
	"net"
	"errors"
)

type Callback func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error)

//长连接请求
func LongRequest(wxClient net.Conn, msg *pb.WechatMsg, callFunction Callback) (*pb.WechatMsg, error) {
	reply, err := HelloWechat(msg)
	if err != nil {
		LogWrite(LOG_ERROR, "请求封包异常，error: "+err.Error())
		return nil, err
	}
	newBuffer := new(bytes.Buffer)
	newBuffer.Write(reply.GetBaseMsg().GetLongHead())
	newBuffer.Write(reply.GetBaseMsg().GetPayloads())
	_, writeError := (wxClient).Write(newBuffer.Bytes())
	if writeError != nil {
		if strings.Contains(writeError.Error(), "connection was aborted by the software") {
			LogWrite(LOG_ERROR, "连接已被微信服务器断开，error: "+writeError.Error())
		}
		if strings.Contains(writeError.Error(), "connection timed out") {
			LogWrite(LOG_ERROR, "微信服务器已经连接超时，error: "+writeError.Error())
		} else {
			LogWrite(LOG_ERROR, "微信服务器传输数据失败，error: "+writeError.Error())
		}
		return nil, err
	}
	if callFunction == nil {
		return nil, nil
	}
	//获取缓存数据
	recv := GetWechatBufferCache(reply.GetBaseMsg().GetLongHead())
	if recv == nil {
		LogWrite(LOG_ERROR, "未取到微信缓存包，即将循环5秒去等待获取")
		currentTime := time.Now().Unix()
		nowTime := time.Now().Unix()
		for nowTime-currentTime <= 5 && recv == nil {
			recv = GetWechatBufferCache(reply.GetBaseMsg().GetLongHead())
			nowTime = time.Now().Unix()
		}
		if recv == nil {
			LogWrite(LOG_ERROR, "5秒后仍未获取到，退出循环抛出异常")
			return nil, errors.New("获取缓存失败")
		}
	}
	return callFunction(reply, recv)
}

//短连接请求
func ShortRequest(msg *pb.WechatMsg, shortHost string, callFunction Callback) (*pb.WechatMsg, error) {
	reply, err := HelloWechat(msg)
	if err != nil {
		LogWrite(LOG_ERROR, "请求封包异常，error: "+err.Error())
		return nil, err
	}
	res, err := http.Post("http://"+shortHost+reply.GetBaseMsg().GetCmdUrl(), "", bytes.NewReader(reply.GetBaseMsg().GetPayloads()))
	if err != nil {
		LogWrite(LOG_ERROR, "请求微信短连接失败，error: "+err.Error())
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogWrite(LOG_ERROR, "微信短连接响应结果获取失败，error: "+err.Error())
		return nil, err
	}
	if callFunction != nil {
		return callFunction(reply, body)
	}
	reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
	reply.BaseMsg.Payloads = body
	return HelloWechat(reply)
}
