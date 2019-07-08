package client

import (
	"wechat-client/client/system"
	pb "wechat-client/proto"
)

//添加好友
func AddUser(wxUser string, v1 string, v2 string, sence string, content string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdAddUser, []byte(`{
		"Encryptusername": "`+v1+`",
		"Ticket": "`+v2+`",
		"Type": 2,
		"Sence": `+sence+`,
		"Content": "`+content+`"
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "发送好友申请失败,error："+err.Error())
		return err
	}
	print(string(ret.BaseMsg.Payloads))
	return nil
}

//接受好友申请
func AcceptUser(wxUser string, v1 string, v2 string, sence string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdAcceptUser, []byte(`{
		"Encryptusername": "`+v1+`",
		"Ticket": "`+v2+`",
		"Type": 3,
		"Sence": `+sence+`,
		"Content": ""
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "接受好友申请失败,error："+err.Error())
		return err
	}
	print(string(ret.BaseMsg.Payloads))
	return nil
}

//拉黑
func SetUserBlack(wxUser string, wxid string, isBlack bool) error {
	bitVal := "15"
	if !isBlack {
		bitVal = "7"
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 2,
		"CmdBuf":"`+wxid+`",
		"BitVal": `+bitVal+`,
		"Remark":""
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "设置联系人黑名单失败,error："+err.Error())
		return err
	}
	return nil
}

//设置联系人标星
func SetUserStar(wxUser string, user string, isStar bool) error {
	bitVal := "71"
	if !isStar {
		bitVal = "7"
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 2,
		"CmdBuf":"`+user+`",
		"BitVal": `+bitVal+`,
		"Remark":""
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "设置联系人星标失败,error："+err.Error())
		return err
	}
	return nil
}

//设置联系人置顶
func SetUserTop(wxUser string, wxid string, isTop bool) error {
	bitVal := "2055"
	if !isTop {
		bitVal = "7"
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 2,
		"CmdBuf":"`+wxid+`",
		"BitVal": `+bitVal+`,
		"Remark":""
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "设置联系人置顶失败,error："+err.Error())
		return err
	}
	return nil
}

//设置备注
func SetUserRemark(wxUser string, user string, mark string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 2,
		"BitVal": 7,
		"Remark": "`+mark+`",
		"CmdBuf": "`+user+`"
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "设置备注失败,error："+err.Error())
		return err
	}
	return nil
}

//删除好友
func DelUser(wxUser string, user string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 7,
		"CmdBuf": "`+user+`"
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "删除好友失败,error："+err.Error())
		return err
	}
	return nil
}

//设置自动通过好友请求（异常）
func SetAutoAcceptUser(wxUser string, isAuto bool) error {
	addVerity := "1"
	if isAuto {
		addVerity = "2"
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdContactOp, []byte(`{
		"Cmdid": 23,
		"AddFromMobile":2,
		"AddFromWechat":2,
		"AddFromChatroom":2,
		"AddFromQrcode"2:,
		"AddFromCard":2,
		"SnsOpenFlag":0,
		"AddVerity": `+addVerity+`
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_INFO, "设置自动通过好友申请失败,error："+err.Error())
		return err
	}
	return nil
}
