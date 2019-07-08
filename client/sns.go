package client

import (
	pb "wechat-client/proto"
	"encoding/json"
	"strconv"
	"time"
	"wechat-client/client/system"
)

const (
	CmdSyncSns     = 214
	CmdSnsUpload   = 207
	CmdSnsSend     = 209
	CmdUserSns     = 212
	CmdSnsTimeLine = 211
	CmdSnsComment  = 213
	CmdSnsOp       = 218
)

type UploadSnsImage struct {
	StartPos  int
	ClientId  string
	TotalLen  int
	Uploadbuf []byte
}

type UploadSnsResponse struct {
	Url string
}

type SendSnsRequest struct {
	Content []byte
}

//同步朋友圈最新动态
func SyncSns(wxUser string) (interface{}, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return nil, err
	}
	requestProto := client_system.CreateWechatMsg(CmdSyncSns, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "同步朋友圈信息失败了,error："+err.Error())
		return nil, err
	}
	wxClient.VUser = ret.GetBaseMsg().GetUser()
	SetWechatConn(wxClient)
	client_system.LogWriteData(string(ret.GetBaseMsg().GetPayloads()))
	var responseData interface{}
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &responseData)
	return responseData, nil
}

/**
 * 上传朋友圈图片
 */
func UploadSns(wxUser string, image []byte) (map[string]string, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return nil, err
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsUpload, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	clientImgId := wxClient.VUser.Userame + `_` + strconv.Itoa(int(time.Now().Unix()))
	maxLength := 102400 //65535
	fileLength := len(image)
	startPos := 0
	var ret *pb.WechatMsg
	for startPos < fileLength {
		count := 0
		if fileLength-startPos > maxLength {
			count = maxLength
		} else {
			count = fileLength - startPos
		}
		sendBuffer := image[startPos:(startPos + count)]
		requestProto.BaseMsg.Payloads, _ = json.Marshal(UploadSnsImage{startPos, clientImgId, fileLength, sendBuffer})
		startPos += count
		if startPos < fileLength {
			_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
			if err != nil {
				return nil, err
			}
		} else {
			ret, err = client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
				reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
				reply.BaseMsg.Payloads = recv
				return client_system.HelloWechat(reply)
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return map[string]string{"url": string(ret.GetBaseMsg().GetPayloads())}, nil
}

/**
 * 发送朋友圈（异常）
 */
func SendSns(wxUser string, content string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	payloadJson, _ := json.Marshal(SendSnsRequest{[]byte(content)})
	requestProto := client_system.CreateWechatMsg(CmdSnsSend, payloadJson)
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "同步朋友圈信息失败了,error："+err.Error())
		return err
	}
	return nil
}

/**
 * 获取用户朋友圈
 * md5Page 首页为空 第二页请附带md5
 * user    访问好友朋友圈的wxid
 * maxId   首页为0 次页朋友圈数据id 的最小值
 */
func GetUserSns(wxUser string, user string, md5Page string, maxId int) interface{} {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdUserSns, []byte(`{
		"Username": "`+user+`",
		"FirstPageMd5": "`+md5Page+`",
		"MaxId": `+strconv.Itoa(maxId)+`
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取好友朋友圈信息失败了,error："+err.Error())
		return err
	}
	client_system.LogWriteData(string(ret.GetBaseMsg().GetPayloads()))
	var responseData interface{}
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &responseData)
	return responseData
}

/**
 * 获取自己朋友圈的动态
 * lastId 为朋友圈最小数据ID，主要是为了翻页吧
 */
func SnsTimeLine(wxUser string, lastId string) interface{} {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsTimeLine, []byte(`{
		"ClientLatestId": `+lastId+`
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -CmdUserSns
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取朋友圈动态失败了,error："+err.Error())
		return err
	}
	client_system.LogWriteData(string(ret.GetBaseMsg().GetPayloads()))
	var responseData interface{}
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &responseData)
	return responseData
}

/**
 * 删除朋友圈
 * ids 动态id
 */
func DelSns(wxUser string, ids string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsOp, []byte(`{
		"Ids": "`+ids+`",
        "CommentId": 0,
		"Type": 1
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "删除动态失败了,error："+err.Error())
		return err
	}
	return nil
}

/**
 * 设置隐私
 * ids 动态id
 * isPrivate true隐藏|false公开
 */
func SetSnsPrivacy(wxUser string, ids string, isPrivate bool) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	opType := "3"
	if isPrivate {
		opType = "2"
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsOp, []byte(`{
		"Ids":"`+ids+`",
        "CommentId": 0,
		"Type": `+opType+`
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "动态设置隐私失败了失败了,error："+err.Error())
		return err
	}
	return nil
}

/**
 * 点赞朋友圈
 * id 动态id
 * user 指定用户
 * isPraise true点赞|false取消点赞
 */
func SetSnsPraise(wxUser string, id string, user string, isPraise bool) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	var cmd int32 = CmdSnsOp
	payloadStr := `{
		"Ids": "` + id + `",
		"CommentId": 0,
		"Type": 5
	}`
	if isPraise {
		cmd = CmdSnsComment
		payloadStr = `{
			"ID": "` + id + `",
			"ToUsername": "` + user + `",
			"Type": 1,
			"Content": ""
		}`
	}
	requestProto := client_system.CreateWechatMsg(cmd, []byte(payloadStr))
	requestProto.BaseMsg.User = wxClient.VUser
	if isPraise {
		_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
			return nil, nil
		})
	} else {
		_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
	}
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "点赞操作失败了,error："+err.Error())
		return err
	}
	return nil
}

/**
 * 评论朋友圈
 * id 指定的朋友圈的动态id
 * user 指定用户(用户回复其他回复者)
 * content 评论内容
 */
func CommentSns(wxUser string, id string, user string, content string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsComment, []byte(`{
		"ID": "`+id+`",
		"ToUsername": "`+user+`",
		"Type": 2,
		"Content": "`+content+`"
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		return nil, nil
	})
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "评论动态失败了,error："+err.Error())
		return err
	}
	return nil
}

/**
 * 删除评论
 * ids 朋友圈的动态id
 * commentId 评论id
 */
func DelCommentSns(wxUser string, ids string, commentId string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdSnsOp, []byte(`{
		"Ids": "`+ids+`",
        "CommentId": "`+commentId+`",
		"Type": 4
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "删除评论失败了,error："+err.Error())
		return err
	}
	return nil
}
