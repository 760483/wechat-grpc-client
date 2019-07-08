package main

import (
	"wechat-client/client"
	"wechat-client/client/system"
	pb "wechat-client/proto"
)

func main() {
	go client.QuickUsernameLogin("mr.zhou", "libgjie123", "lingjie123", "", func(wxUser string, ret *pb.WechatMsg) {
		loginSuccess(wxUser, ret.BaseMsg.User.Nickname)
	})

	//go client.QuickQrcodeLogin("mr.zhou", "", func(wxUser string, ret *pb.WechatMsg) {
	//	loginSuccess(wxUser, ret.BaseMsg.User.Nickname)
	//})

	select {}
}

func loginSuccess(wxUser string, nickname []byte) {
	print("【" + string(nickname) + " 】登录成功，正在初始化消息系统...\n")
	if client.LoginInit(wxUser) == nil {
		print("正在拉取历史消息")
		for {
			if client.SyncMsg(wxUser) == nil {
				break
			}
			print(".")
		}
		client_system.AsyncCallFunc = client.CallFunc
		go client.SendHeartbeat(wxUser)
		print("\n初始化完成\n")
	}
	go client.SyncContactAllDetail(wxUser)
}
