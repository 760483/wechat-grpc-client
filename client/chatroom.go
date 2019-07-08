package client

const (
	CmdCreateRoom      = 119
	CmdAddRoomMembers  = 120
	CmdInviteMember    = 610
	CmdDelMember       = 179
	CmdSetRoom         = 681
	CmdGetRoomDetail   = 551
	CmdGetQrcode       = 168
	CmdSetAnnouncement = 993
	CmdTransferRoom    = 990
	CmdAgreeMember     = 774
	CmdAddRoomAdmin    = 889
	CmdDelRoomAdmin    = 259
)

//创建群聊
func CreateChatroom(wxUser string, members string) ([]byte, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Membernames": "` + members + `"
	}`)
	)
	FastLongRequest(wxUser, CmdCreateRoom, payload, "创建群聊失败", &result, &err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
 * 添加群成员，直接拉取好友
 */
func AddRoomMembers(wxUser string, roomId string, members string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Roomeid": "` + roomId + `",
		"Membernames": "` + members + `"
	}`)
	)
	FastLongRequest(wxUser, CmdAddRoomMembers, payload, "拉取好友进去失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 邀请好友加入群聊(链接)
 */
func InviteRoomMember(wxUser string, roomId string, user string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"ChatRoom": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastShortRequest(wxUser, CmdInviteMember, payload, "发送进群邀请链接失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 删除群成员
 */
func DelRoomMember(wxUser string, roomId string, user string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"ChatRoom": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastShortRequest(wxUser, CmdDelMember, payload, "删除群成员失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 获取群信息
 */
func GetRoomDetail(wxUser string, roomId string) ([]byte, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Chatroom": "` + roomId + `"
	}`)
	)
	FastShortRequest(wxUser, CmdGetRoomDetail, payload, "获取群信息失败", &result, &err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
 * 设置群名称
 */
func SetRoomName(wxUser string, roomId string, roomName string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid": 27,
		"Roomname": "` + roomName + `",
		"ChatRoom": "` + roomId + `"
	}`)
	)
	FastShortRequest(wxUser, CmdSetRoom, payload, "修改群名称失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 设置群聊保存通讯录
 */
func SetRoomSave(wxUser string, roomId string, isSave bool) error {
	bitVal := "2051"
	if !isSave {
		bitVal = "2"
	}
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid": 2,
		"CmdBuf":"` + roomId + `",
		"BitVal": ` + bitVal + `
	}`)
	)
	FastShortRequest(wxUser, CmdSetRoom, payload, "设置群聊通讯录保存失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 设置联系人置顶
 */
func SetRoomTop(wxUser string, roomId string, isTop bool) error {
	bitVal := "2055"
	if !isTop {
		bitVal = "7"
	}
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid": 2,
		"CmdBuf":"` + roomId + `",
		"BitVal": ` + bitVal + `,
		"Remark":""
	}`)
	)
	FastShortRequest(wxUser, CmdSetRoom, payload, "设置群聊通讯录置顶失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 获取群聊二维码
 */
func GetRoomQrcode(wxUser string, roomId string) ([]byte, error) {

	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Username": "` + roomId + `"
	}`)
	)
	FastLongRequest(wxUser, CmdGetQrcode, payload, "获取群二维码失败", &result, &err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
 * 设置进群邀请验证
 * 0关闭 2开启
 */
func SetRoomVerify(wxUser string, roomId string, status string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid": 66,
		"BitVal": ` + status + `,
		"ChatRoom": "` + roomId + `"
	}`)
	)
	FastShortRequest(wxUser, CmdSetRoom, payload, "设置进群邀请验证失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 设置群公告
 */
func SetRoomAnnouncement(wxUser string, roomId string, content string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Announcement": "` + content + `",
		"ChatRoomName": "` + roomId + `"
	}`)
	)
	FastShortRequest(wxUser, CmdSetAnnouncement, payload, "设置群公告失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 同意别的用户邀请进群
 */
func AgreeInviteRoom(wxUser string, roomId string, user string, ticket string, inviteUser string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Ticket": "` + ticket + `",
		"Inviterusername": "` + inviteUser + `",
		"Roomname": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastShortRequest(wxUser, CmdAgreeMember, payload, "同意进群操作失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 群组设置备注
 */
func ExitRoom(wxUser string, roomId string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Cmdid": 16,
		"ChatRoom": "` + roomId + `"
	}`)
	)
	FastShortRequest(wxUser, CmdSetRoom, payload, "退出群聊失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 转让群
 */
func TransferRoom(wxUser string, roomId string, user string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"ChatRoomName": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastLongRequest(wxUser, CmdTransferRoom, payload, "转让群失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 添加群管理员
 */
func AddRoomAdmin(wxUser string, roomId string, user string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"ChatRoomName": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastLongRequest(wxUser, CmdAddRoomAdmin, payload, "添加群管理员失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 删除群管理员
 */
func DelRoomAdmin(wxUser string, roomId string, user string) error {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Chatroom": "` + roomId + `",
		"Username": "` + user + `"
	}`)
	)
	FastLongRequest(wxUser, CmdDelRoomAdmin, payload, "删除群管理员失败", &result, &err)
	if err != nil {
		return err
	}
	return nil
}
