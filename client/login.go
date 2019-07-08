package client

import (
	"encoding/json"
	"wechat-client/client/system"
	pb "wechat-client/proto"
	"fmt"
	"encoding/base64"
	"errors"
	"strings"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
	"crypto/md5"
	"io/ioutil"
)

type QrcodeInfo struct {
	CheckTime   int32
	ExpiredTime int32
	HeadImgUrl  string
	ImgBuf      string
	Nickname    string
	NotifyKey   string
	Password    string
	RandomKey   string
	Status      int32
	Username    string
	Uuid        string
	LongHead    []byte
}

const (
	CmdGetLoginQrcode   = 502
	CmdCheckLoginQrcode = 503
	CmdAutoLogin        = 702
	CmdLoginQrcode      = 1111
	CmdParseLogin       = -1001
	CmdUsernameLogin    = 2222
	CmdLogoutOther      = 281
	CmdLogoutSelf       = 282
	CmdLoginInit        = 1002
)

//获取登录二维码
func GetLoginQrcode(wxUser string, deviceId string) (QrcodeInfo, error) {
	requestProto := client_system.CreateWechatMsg(CmdGetLoginQrcode, []byte{})
	if deviceId != "" {
		requestProto.BaseMsg.User.DeviceId = deviceId
	}
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return QrcodeInfo{}, errors.New("获取微信实例缓存信息失败了,error：" + err.Error())
	}
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -CmdGetLoginQrcode
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		return QrcodeInfo{}, errors.New("获取二维码失败了,error：" + err.Error())
	}
	var qrcodeInfo QrcodeInfo
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &qrcodeInfo)
	wxClient.LoginParam.LongHead = ret.GetBaseMsg().GetLongHead()
	wxClient.LoginParam.UUID = []byte(qrcodeInfo.Uuid)
	wxClient.LoginParam.MaxSyncKey, _ = base64.StdEncoding.DecodeString(qrcodeInfo.NotifyKey)
	wxClient.VUser = ret.BaseMsg.User
	SetWechatConn(wxClient)
	qrcodeInfo.LongHead = ret.GetBaseMsg().GetLongHead()
	return qrcodeInfo, nil
}

//验证二维码状态
func CheckQrcode(wxUser string) (QrcodeInfo, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return QrcodeInfo{}, errors.New("获取微信实例缓存信息失败了,error：" + err.Error())
	}
	requestProto := client_system.CreateWechatMsg(CmdCheckLoginQrcode, wxClient.LoginParam.UUID)
	requestProto.BaseMsg.User.MaxSyncKey = wxClient.LoginParam.MaxSyncKey
	requestProto.BaseMsg.LongHead = wxClient.LoginParam.LongHead
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		if recv[16] != 191 {
			return nil, errors.New("微信数据包异常")
		}
		reply.BaseMsg.Cmd = -CmdCheckLoginQrcode
		reply.BaseMsg.Payloads = recv
		reply.BaseMsg.User.MaxSyncKey = requestProto.BaseMsg.User.MaxSyncKey
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		return QrcodeInfo{}, errors.New("获取二维码状态失败了,error：" + err.Error())
	}
	var checkQrcodeInfo QrcodeInfo
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &checkQrcodeInfo)
	checkQrcodeInfo.LongHead = ret.GetBaseMsg().GetLongHead()
	return checkQrcodeInfo, nil
}

/**
 * 二维码登录操作
 */
func QrcodeLogin(wxUser, username string, password string) (*pb.WechatMsg, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return nil, errors.New("获取微信实例缓存信息失败了,error：" + err.Error())
	}
	requestProto := client_system.CreateWechatMsg(CmdLoginQrcode, []byte(`{
		"Username":   "`+username+`",
		"Password":   "`+password+`",
		"UUid":       "`+strings.ToUpper(fmt.Sprintf("%s", uuid.Must(uuid.NewV4())))+`",
		"DeviceType": "`+createDeviceType()+`",
		"DeviceName": "`+createDeivceName()+`",
		"ProtocolVer": 1
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		if len(recv) == 509 {
			return nil, errors.New("微信包异常")
		}
		if recv[16] != 191 {
			return nil, errors.New("微信包异常")
		}
		reply.BaseMsg.Cmd = CmdParseLogin
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		return nil, errors.New("二维码登录失败了,error：" + err.Error())
	}
	//微信重定向连接
	if ret.GetBaseMsg().GetRet() == -301 {
		ConnectWx(wxUser, ret.GetBaseMsg().GetLongHost())
		return QrcodeLogin(wxUser, username, password)
	}
	if ret.GetBaseMsg().GetRet() != 0 {
		return nil, errors.New(fmt.Sprintf("登录异常 /n 状态码: %d \n 异常信息: %s", ret.GetBaseMsg().GetRet(), string(ret.GetBaseMsg().GetPayloads())))
	}
	wxClient.VUser = ret.GetBaseMsg().GetUser()
	wxClient.LongHost = ret.GetBaseMsg().GetLongHost()
	wxClient.ShortHost = ret.GetBaseMsg().GetShortHost()
	SetWechatConn(wxClient)
	return ret, nil
}

/**
 * 账号密码登录
 */
func UsernameLogin(wxUser string, username string, password string, deviceId string) (*pb.WechatMsg, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return nil, errors.New("获取微信实例缓存信息失败了,error：" + err.Error())
	}
	requestProto := client_system.CreateWechatMsg(CmdUsernameLogin, []byte(`{
		"Username":   "`+username+`",
		"Password":   "`+password+`",
		"UUid":       "`+strings.ToUpper(fmt.Sprintf("%s", uuid.Must(uuid.NewV4())))+`",
		"DeviceType": "`+createDeviceType()+`",
		"DeviceName": "`+createDeivceName()+`",
		"ProtocolVer": 1
	}`))
	//如果设置过设备号，使用设置设备号。如果未传，那么使用用户名md5设备号
	if deviceId != "" {
		requestProto.BaseMsg.User.DeviceId = deviceId
	} else {
		requestProto.BaseMsg.User.DeviceId = GetSixTwoDataMd5(username)
	}
	if wxClient.ShortHost == "" {
		wxClient.ShortHost = "short.weixin.qq.com"
	}
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = CmdParseLogin
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		return nil, errors.New("账号密码登录失败了,error：" + err.Error())
	}
	//微信重定向连接
	if ret.GetBaseMsg().GetRet() == -301 {
		wxClient.ShortHost = ret.GetBaseMsg().GetShortHost()
		wxClient.LongHost = ret.GetBaseMsg().GetLongHost()
		SetWechatConn(wxClient)
		ConnectWx(wxUser, ret.GetBaseMsg().GetLongHost())
		return UsernameLogin(wxUser, username, password, deviceId)
	}
	if ret.GetBaseMsg().GetRet() != 0 {
		return nil, errors.New(fmt.Sprintf("登录异常 /n 状态码: %d \n 异常信息: %s", ret.GetBaseMsg().GetRet(), string(ret.GetBaseMsg().GetPayloads())))
	}
	wxClient.VUser = ret.GetBaseMsg().GetUser()
	wxClient.LongHost = ret.GetBaseMsg().GetLongHost()
	wxClient.ShortHost = ret.GetBaseMsg().GetShortHost()
	SetWechatConn(wxClient)
	return ret, nil
}

/**
 * token登录操作
 */
func AutoLogin(wxUser string) (*pb.WechatMsg, error) {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return nil, errors.New("获取微信实例缓存信息失败了,error：" + err.Error())
	}
	requestProto := client_system.CreateWechatMsg(CmdAutoLogin, []byte(`{
		"UUid":       "`+strings.ToUpper(fmt.Sprintf("%s", uuid.Must(uuid.NewV4())))+`",
		"DeviceType": "`+createDeviceType()+`",
		"DeviceName": "`+createDeivceName()+`",
		"ProtocolVer": 1
	}`))
	requestProto.BaseMsg.User = wxClient.VUser
	ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
		reply.BaseMsg.Payloads = recv
		return client_system.HelloWechat(reply)
	})
	if err != nil {
		return nil, errors.New("自动登录失败了,error：" + err.Error())
	}
	if ret.GetBaseMsg().GetRet() == -301 {
		wxClient.ShortHost = ret.GetBaseMsg().GetShortHost()
		wxClient.LongHost = ret.GetBaseMsg().GetLongHost()
		SetWechatConn(wxClient)
		ConnectWx(wxUser, ret.GetBaseMsg().GetLongHost())
		return AutoLogin(wxUser)
	}
	if ret.GetBaseMsg().GetRet() != 0 {
		return nil, errors.New(fmt.Sprintf("登录异常 /n 状态码: %d \n 异常信息: %s", ret.GetBaseMsg().GetRet(), string(ret.GetBaseMsg().GetPayloads())))
	}
	if ret.GetBaseMsg().GetLongHost() != "" {
		wxClient.LongHost = ret.GetBaseMsg().GetLongHost()
	}
	if ret.GetBaseMsg().GetShortHost() != "" {
		wxClient.ShortHost = ret.GetBaseMsg().GetShortHost()
	}
	wxClient.VUser = ret.GetBaseMsg().GetUser()
	SetWechatConn(wxClient)
	return ret, nil
}

//退出当前登录,清除缓存
func Logout(wxUser string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdLogoutSelf, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	_, err = client_system.ShortRequest(requestProto, wxClient.ShortHost, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
		return nil, nil
	})
	if err != nil {
		return err
	}
	client_system.DeleteCache(wxUser)
	return nil
}

//登录成功初始化操作
func LoginInit(wxUser string) error {
	wxClient, err := GetWechatConn(wxUser)
	if err != nil {
		return err
	}
	requestProto := client_system.CreateWechatMsg(CmdLoginInit, []byte{})
	requestProto.BaseMsg.User = wxClient.VUser
	var initCode int32
	for initCode != 8888 {
		ret, err := client_system.LongRequest(wxClient.Client, requestProto, func(reply *pb.WechatMsg, recv []byte) (*pb.WechatMsg, error) {
			reply.BaseMsg.Cmd = -reply.BaseMsg.Cmd
			reply.BaseMsg.Payloads = recv
			return client_system.HelloWechat(reply)
		})
		if err != nil {
			return err
		}
		initCode = ret.GetBaseMsg().GetRet()
		wxClient.VUser = ret.GetBaseMsg().GetUser()
		SetWechatConn(wxClient)
		requestProto.BaseMsg.User = ret.GetBaseMsg().GetUser()
	}
	return nil
}

//二维码状态验证协程
func QuickQrcodeLogin(wxUser string, deviceId string, callback func(wxUser string, ret *pb.WechatMsg)) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := ConnectWx(wxUser, "long.weixin.qq.com"); err != nil {
		print(err.Error())
		return ""
	}
	qrcodeInfo, err := GetLoginQrcode(wxUser, deviceId)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "获取二维码失败", err.Error()))
		return ""
	}
	if qrcodeInfo.Uuid == "" {
		client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "获取二维码失败", "uuid空的"))
		return ""
	}
	qrcodeFile, _ := base64.StdEncoding.DecodeString(qrcodeInfo.ImgBuf)
	ioutil.WriteFile("./runtime/qrcode/qrcode_"+wxUser+".jpg", qrcodeFile, 0666)
	checkQrcodeFunc := func() {
		loginStatus := false
		isFirstCheck1 := false
		isFirstCheck2 := false
		var user *pb.WechatMsg
		for {
			time.Sleep(1 * time.Second)
			checkQrcodeInfo, err := CheckQrcode(wxUser)
			if err != nil {
				client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "获取二维码登录状态失败", err.Error()))
				continue
			}
			switch checkQrcodeInfo.Status {
			case 0:
				if !isFirstCheck1 {
					println("等待扫码中...")
					isFirstCheck1 = true
				}
				break
			case 1:
				if !isFirstCheck2 {
					println("【" + checkQrcodeInfo.Nickname + "】已扫码，等待确认...")
					isFirstCheck2 = true
				}
				break
			case 2:
				println("【" + checkQrcodeInfo.Nickname + "】已确认，准备登陆...")
				user, err = QrcodeLogin(wxUser, checkQrcodeInfo.Username, checkQrcodeInfo.Password)
				if err != nil {
					client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "二维码登录操作失败了", err.Error()))
					break
				}
				if user.GetBaseMsg().GetUser().GetUserame() != "" {
					loginStatus = true
					println("【" + string(user.GetBaseMsg().GetUser().GetNickname()) + "】登陆成功")
				}
				break
			case 4:
				println("【" + checkQrcodeInfo.Nickname + "】已取消扫码")
				break
			case -2007:
				println("二维码已过期")
				break
			}
			if checkQrcodeInfo.Status == 2 || checkQrcodeInfo.Status == 4 || checkQrcodeInfo.Status == -2007 {
				break
			}
		}
		if loginStatus {
			callback(wxUser, user)
		} else {
			wxClient, _ := GetWechatConn(wxUser)
			if wxClient != nil {
				wxClient.Client.Close()
			}
			client_system.DeleteCache(wxUser)
		}
	}
	go checkQrcodeFunc()
	return "http://weixin.qq.com/x/" + qrcodeInfo.Uuid
}

func QuickUsernameLogin(wxUser, username, password, deviceId string, callback func(wxUser string, ret *pb.WechatMsg)) {
	if err := ConnectWx(wxUser, "long.weixin.qq.com"); err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "创建tcp微信客户端连接失败了", err.Error()))
		return
	}
	ret, err := UsernameLogin(wxUser, username, password, deviceId)
	if err != nil {
		client_system.LogWrite(client_system.LOG_ERROR, fmt.Sprintf("【%s】%s,error: %s", wxUser, "账号密码登录失败失败了", err.Error()))
		return
	}
	callback(wxUser, ret)
}

func createMac() string {
	var m [6]byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		macByte := rand.Intn(256)
		m[i] = byte(macByte)
		rand.Seed(int64(macByte))
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", m[0], m[1], m[2], m[3], m[4], m[5])
}

func createDeivceName() string {
	var deviceName string
	switch rand.Intn(13) {
	case 0:
		deviceName = "iPad"
		break
	case 1:
		deviceName = "iPad"
		break
	case 2:
		deviceName = "iPad 2"
		break
	case 3:
		deviceName = "iPad 3"
		break
	case 4:
		deviceName = "iPad 4"
		break
	case 5:
		deviceName = "iPad 5"
		break
	case 6:
		deviceName = "iPad Air"
		break
	case 7:
		deviceName = "iPad Air 2"
		break
	case 8:
		deviceName = "iPad mini"
		break
	case 9:
		deviceName = "iPad mini 2"
		break
	case 10:
		deviceName = "iPad mini 3"
		break
	case 11:
		deviceName = "iPad mini 4"
		break
	case 12:
		deviceName = "iPad Pro"
		break
	case 13:
		deviceName = "iPad Pro 2"
		break
	}
	return deviceName
}

func createDeviceType() string {
	var wifiName, operators, macAddress string
	macAddress = createMac()
	switch rand.Intn(5) {
	case 0:
		wifiName = "TP-Link_"
		break
	case 1:
		wifiName = "MI-WIFI-"
		break
	case 2:
		wifiName = "Tenda_"
		break
	case 3:
		wifiName = "HUAWEI-"
		break
	case 4:
		wifiName = "360WIFI_"
		break
	}
	switch rand.Intn(3) {
	case 0:
		operators = "中国移动"
		break
	case 1:
		operators = "中国联通"
		break
	case 2:
		operators = "中国电信"
		break
	}
	wifiName = wifiName + strings.ToUpper(string([]byte(macAddress)[0:2])) + strings.ToUpper(string([]byte(macAddress)[3:5])) + strings.ToUpper(string([]byte(macAddress)[6:8]))
	return fmt.Sprintf("<k21>%s</k21><k22>%s</k22><k24>%s</k24>", wifiName, operators, macAddress)
}

func GetSixTwoDataMd5(str string) string {
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
}
