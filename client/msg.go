package client

import (
	"strconv"
	"time"
	"encoding/json"
	"wechat-client/client/system"
	"fmt"
	"bytes"
	"io/ioutil"
	"strings"
)

type SendImageRequest struct {
	ClientImgId string
	ToUserName  string
	StartPos    int
	TotalLen    int
	DataLen     int
	Data        []byte
}

type SendFaceRequest struct {
	ClientMsgId string
	ToUserName  string
	StartPos    int
	TotalLen    int
	Md5         string
	ExternXml   string
}

type SendVoiceRequest struct {
	ToUserName  string
	Offset      int32
	VoiceLength int32
	Length      int32
	EndFlag     int8
	VoiceFormat int8
	Data        []byte
}

type SendCdnImageRequest struct {
	ClientImgId       string
	ToUserName        string
	StartPos          int
	TotalLen          int
	DataLen           int
	CDNMidImgUrl      string
	AESKey            string
	CDNMidImgSize     int
	CDNThumbImgSize   int
	CDNThumbImgHeight int
	CDNThumbImgWidth  int
}

type SendCdnVideoRequest struct {
	ToUserName    string
	ThumbTotalLen int
	ThumbStartPos int
	VideoTotalLen int
	VideoStartPos int
	PlayLength    int
	AESKey        string
	CDNVideoUrl   string
}

type SendMsgRequest struct {
	ToUserName string
	Content    string
	MsgSource  string
	Type       int
}

type MassSendTextRequest struct {
	ToList       string
	ToListMd5    string
	DataStartPos int
	DataTotalLen int
	ClientId     string
	DataBuffer   []byte
	MsgType      int
	VoiceFormat  int
	//ToList = tb_ToUsername.Text,                      //多微信wxid 分号分割如 wxid1;wxid2;wxid3;
	//ToListMd5 = Utils.MD5Encrypt(tb_ToUsername.Text), //
	//DataStartPos = 0,
	//DataTotalLen = buf.Length, //音频大小在120KB以内
	//ClientId = ClientImgId,
	//DataBuffer = buf,
	//MsgType = 1,
	//VoiceFormat = 0 //0->Arm 无amr音频头 2->MP3 3->WAVE 4->SILK
}

type DownloadImg struct {
	MsgId        int64
	ToUsername   string
	StartPos     int
	TotalLen     int
	DataLen      int
	CompressType int //hdlength  1高清0缩略
	Md5          string `json:"md5"`
}

type DownloadVoice struct {
	MsgId       int64
	StartPos    int
	DataLen     int
	TotalLen    int
	ClientMsgId string
	Bufid       string `json:"bufid"`
}

const (
	CmdSyncMsg       = 138
	CmdParseMsg      = -318
	CmdSendMsg       = 522
	CmdSendAppMsg    = 222
	CmdSendVoice     = 127
	CmdSendImage     = 110
	CmdSendFace      = 175
	CmdSendVideo     = 149
	CmdMassSendText  = 193
	CmdDownloadImg   = 109
	CmdDownloadVoice = 128
)

//同步消息
func SyncMsg(wxUser string) []CallbackMsg {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return nil
	}
	if wxClient.MsgSyncKey != nil && wxClient.MsgLock {
		return nil
	}
	//快速设置锁
	wxClient.MsgLock = true
	SetWechatConn(wxClient)
	requestProto := client_system.CreateWechatMsg(CmdSyncMsg, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	requestProto.BaseMsg.User.MaxSyncKey = wxClient.MsgMaxSyncKey
	requestProto.BaseMsg.User.CurrentsyncKey = wxClient.MsgSyncKey
	ret, err := client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
	if err != nil {
		defer func() {
			ReConnect(wxUser)
			wxClient.MsgLock = false
			SetWechatConn(wxClient)
		}()
		client_system.LogWrite(client_system.LOG_INFO, "微信消息拉取失败,error："+err.Error())
		return []CallbackMsg{}
	}
	wxClient.MsgSyncKey = ret.GetBaseMsg().GetUser().GetCurrentsyncKey()
	wxClient.MsgMaxSyncKey = ret.GetBaseMsg().GetUser().GetMaxSyncKey()
	wxClient.MsgLock = false
	SetWechatConn(wxClient)
	var syncMsgList []CallbackMsg
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &syncMsgList)
	return syncMsgList
}

//发送文字消息
func SendTextMsg(wxUser string, wxid string, content string, atList string) error {
	if atList != "" && !strings.Contains(content, "@") {
		atUserStr := ""
		users, _ := GetContact(wxUser, atList, wxid)
		for _, member := range users {
			atUserStr += "@" + member.NickName
		}
		content = atUserStr + " " + content
	}
	var err error
	payload, _ := json.Marshal(SendMsgRequest{wxid, content, atList, 0})
	FastLongRequestOperate(wxUser, CmdSendMsg, payload, "发送文字消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

//发送APP消息
func SendAppMsg(wxUser string, wxid string, content string) error {
	var err error
	payload, _ := json.Marshal(SendMsgRequest{wxid, content, "", 5})
	FastLongRequestOperate(wxUser, CmdSendAppMsg, payload, "发送卡片消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 发送图片消息
 */
func SendImageMsg(wxUser string, wxid string, file []byte) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	clientImgId := wxClient.VUser.Userame + `_` + strconv.Itoa(int(time.Now().Unix()))
	requestProto := client_system.CreateWechatMsg(CmdSendImage, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	//maxLength := 65535
	maxLength := 102400
	fileLength := len(file)
	kb := 1024
	startPos := 0
	for startPos < fileLength {
		count := 0
		if fileLength-startPos > maxLength {
			count = maxLength
		} else {
			count = fileLength - startPos
		}
		sendBuffer := file[startPos:(startPos + count)]
		requestProto.BaseMsg.Payloads, _ = json.Marshal(SendImageRequest{clientImgId, wxid, startPos, fileLength, len(sendBuffer), sendBuffer})
		startPos += count
		//发送请求
		_, err = client_system.LongRequest(wxClient.Client, requestProto, nil)
		if err != nil {
			return err
		}
		print(fmt.Sprintf("\n正在发送图片，已发送 %d kb,剩余 %d kb,总计 %d kb\n", startPos/kb, (fileLength-startPos)/kb, fileLength/kb))
	}
	return nil
}

/**
 * 发送语音消息
 * //0->Arm 无amr音频头 2->MP3 3->WAVE 4->SILK
 */
func SendVoiceMsg(wxUser string, wxid string, file []byte, voiceLength int32) error {
	var err error
	payload, _ := json.Marshal(SendVoiceRequest{wxid, 0, voiceLength, int32(len(file)), 1, 4, file})
	FastLongRequestOperate(wxUser, CmdSendVoice, payload, "发送语音消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 发送名片消息
 * `<msg username=\"wxid_k9jdv2j4n8cf12\" nickname=\"洛小兮\" fullpy=\"luoxiaoxi\" shortpy=\"\" alias=\"wyw521-0312\" imagestatus=\"3\" scene=\"17\" province=\"河南\" city=\"中国\" sign=\"……\" sex=\"1\" certflag=\"0\" certinfo=\"\" brandIconUrl=\"\" brandHomeUrl=\"\" brandSubscriptConfigUrl=\"\" brandFlags=\"0\" regionCode=\"CN_Henan_Zhengzhou\"></msg>`
 */
func SendContactMsg(wxUser string, wxid string, content string) error {
	var err error
	payload, _ := json.Marshal(SendMsgRequest{wxid, content, "", 42})
	FastLongRequestOperate(wxUser, CmdSendMsg, payload, "发送名片消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 发送表情消息消息
 * <<msg><emoji fromusername = "wxid_k9jdv2j4n8cf12" tousername = "5154395343@chatroom" type="2" idbuffer="media:0_0" md5="41688f98bc32ce19ac73534c956ceba9" len = "121264" productid="" androidmd5="41688f98bc32ce19ac73534c956ceba9" androidlen="121264" s60v3md5 = "41688f98bc32ce19ac73534c956ceba9" s60v3len="121264" s60v5md5 = "41688f98bc32ce19ac73534c956ceba9" s60v5len="121264" cdnurl = "http://emoji.qpic.cn/wx_emoji/MXol8KjyEwLTlm5zqy3U9XdowwFREzzxWsypEyuGKPoM7oNq2RK0jucSTXCR2JE9/" designerid = "" thumburl = "" encrypturl = "http://emoji.qpic.cn/wx_emoji/MXol8KjyEwLTlm5zqy3U9XdowwFREzzxWsypEyuGKPrHre70Pric2NEhxZq6wLw5O/" aeskey= "61ea40b631a79c241bd2643fcd9c1f67" externurl = "http://emoji.qpic.cn/wx_emoji/MXol8KjyEwLTlm5zqy3U9XdowwFREzzxWsypEyuGKPo3QqgJK7pN8icNV46FZibtPE/" externmd5 = "eb7aa8d7f61b2f9d43efaa84db623a15" width= "1000" height= "1000" tpurl= "" tpauthkey= "" attachedtext= "" ></emoji> <gameext type="0" content="0" ></gameext></msg>
 */
func SendFaceMsg(wxUser string, faceData SendFaceRequest) error {
	var retsult []byte
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return err
	}
	faceData.ClientMsgId = wxClient.VUser.GetUserame() + `_` + strconv.Itoa(int(time.Now().Unix()))
	payload, _ := json.Marshal(faceData)
	FastShortRequest(wxUser, CmdSendFace, payload, "发送表情消息失败", &retsult, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 发送CDN视频消息
 */
func SendCdnVideoMsg(wxUser string, requestData SendCdnVideoRequest) error {
	var err error
	payload, _ := json.Marshal(requestData)
	FastLongRequestOperate(wxUser, CmdSendVideo, payload, "发送名片消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 发送CDN图片消息
 */
func SendCdnImageMsg(wxUser string, requestData SendCdnImageRequest) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestData.ClientImgId = wxClient.VUser.GetUserame() + `_` + strconv.Itoa(int(time.Now().Unix()))
	payload, _ := json.Marshal(requestData)
	FastLongRequestOperate(wxUser, -CmdSendImage, payload, "发送名片消息失败", &err)
	if err != nil {
		return err
	}
	return nil
}

//下载图片
func DownloadMsgImage(wxUser string, img DownloadImg) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdDownloadImg, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	imgBuffer := new(bytes.Buffer)
	totalLength, pos := img.TotalLen, 0
	for pos < totalLength {
		var length int
		if pos+65536 > totalLength {
			length = totalLength - pos
		} else {
			length = 65536
		}
		img.StartPos = pos
		img.DataLen = length
		img.TotalLen = totalLength
		requestProto.BaseMsg.Payloads, _ = json.Marshal(img)
		ret, err := client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
		if err != nil {
			return err
		}
		pos += length
		imgBuffer.Write(ret.GetBaseMsg().GetPayloads())
	}
	ioutil.WriteFile("./runtime/download/img/img_"+fmt.Sprintf("%d", img.MsgId)+".jpg", imgBuffer.Bytes(), 0666)
	return nil
}

//下载语音消息
func DownloadMsgVoice(wxUser string, voice DownloadVoice) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, "获取微信实例缓存信息失败了,error："+err.Error())
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdDownloadVoice, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	voiceBuffer := new(bytes.Buffer)
	totalLength, pos := voice.TotalLen, 0
	for pos < totalLength {
		var length int
		if pos+65536 > totalLength {
			length = totalLength - pos
		} else {
			length = 65536
		}
		voice.StartPos = pos
		voice.DataLen = length
		requestProto.BaseMsg.Payloads, _ = json.Marshal(voice)
		ret, err := client_system.ShortRequest(requestProto, wxClient.ShortHost, nil)
		if err != nil {
			return err
		}
		pos += length
		voiceBuffer.Write(ret.GetBaseMsg().GetPayloads())
	}
	ioutil.WriteFile("./runtime/download/voice/voice_"+fmt.Sprintf("%d", voice.MsgId)+".silk", voiceBuffer.Bytes(), 0666)
	return nil
}
