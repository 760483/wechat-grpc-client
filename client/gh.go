package client

const (
	CmdAddGh      = 137
	CmdClickGh    = 359
	CmdSearchGh   = 719
	CmdRequestUrl = 233
	CmdDel        = 681
)

/**
 * 关注公众号
 */
func ConcernGh(wxUser string, encryptUsername string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Encryptusername": "` + encryptUsername + `@qr",
		"Ticket": "",
		"Type": 1,
		"Sence": 3,
		"Content": "",
		"ProtocolVer": 1
	}`)
	)
	FastLongRequest(wxUser, CmdAddGh, payload, "关注公众号失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 取消关注
 */
func UnConcernGh(wxUser string, wxid string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid":7,
		"CmdBuf": "` + wxid + `"
	}`)
	)
	FastShortRequest(wxUser, CmdDel, payload, "取消关注公众号失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 点击菜单
 */
func ClickMenu(wxUser string, wxid string, menuInfo string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Clickinfo": "` + menuInfo + `",
		"Username": "` + wxid + `",
		"ClickCommandType": 1
	}`)
	)
	FastLongRequest(wxUser, CmdClickGh, payload, "点击菜单操作失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 搜索文章公众号
 */
func SearchContent(wxUser string, content string) ([]byte, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Searchinfo": "` + content + `"
	}`)
	)
	FastShortRequest(wxUser, CmdSearchGh, payload, "搜索公众号文章失败", &result, &err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
 * 请求url
 * ReqUrl //要获取key的连接 授权登陆时的链接即为转跳链接https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx53c677b55caa45fd&redirect_uri=http%3A%2F%2Fmeidang.cimiworld.com%2Fh5%2Fchourenpin%3Fs%3D77e881961fee12eb65f5497bbff02fac%26from%3Dsinglemessage%26isappinstalled%3D0&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect
 * Scene //2来源好友或群 必须设置来源的id  3 历史阅读 4 二维码连接 7 来源公众号 必须设置公众号的id
 * Username //来源 来源设置wxid 来源群id@chatroom 来源公众号gh_e09c57858a0c原始id
 */
func RequestUrl(wxUser string, url string, username string, scene string) ([]byte, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"ReqUrl": "` + url + `",
		"Username": "` + username + `",
		"Scene": ` + scene + `,
		"ProtocolVer": 1
	}`)
	)
	FastLongRequest(wxUser, CmdRequestUrl, payload, "请求url操作失败", &result, &err)
	if err != nil {
		return nil, err
	}
	return result, nil
}
