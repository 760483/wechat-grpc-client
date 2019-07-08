package client

import (
	"strconv"
	"encoding/json"
	"wechat-client/client/system"
	"strings"
)

//同步联系人username数组结构
type ContactBaseUsername struct {
	Username string
}

//同步联系人列表结构
type ContactBase struct {
	ContinueFlag              int8
	CurrentChatRoomContactSeq int32
	CurrentWxcontactSeq       int32
	UsernameLists             []ContactBaseUsername
}

//同步联系人列表结构
type ContactBaseDetailList struct {
	ContinueFlag              int8
	CurrentChatRoomContactSeq int32
	CurrentWxcontactSeq       int32
	ContactList               []ContactDetail
}

//用户详情数据结构
type ContactDetail struct {
	Alias           string
	BigHeadImgUrl   string
	ChatRoomOwner   string
	ChatroomVersion int32
	City            string
	ContactType     int32
	EncryptUsername string
	ExtInfo         string
	ExtInfoExt      string
	LabelLists      string
	MsgType         int8
	NickName        string
	Province        string
	Remark          string
	Sex             int8
	Signature       string
	SmallHeadImgUrl string
	Ticket          string
	UserName        string
	VerifyFlag      int8
}

type ContactsCache struct {
	User     []ContactDetail
	ChatRoom []ContactDetail
	GH       []ContactDetail
	System   []ContactDetail
}

const (
	CmdSyncContactBase   = 851
	CmdSyncContactDetail = 945
	CmdSearchContact     = 106
	CmdGetContact        = 182
	CmdAddUser           = 1137
	CmdAcceptUser        = 137
	CmdContactOp         = 681
)

/**
 * 同步全部详情通讯录
 */
func SyncContactAllDetail(wxUser string) {
	var currentWxcontactSeq, currentChatRoomContactSeq int32
	var continueFlag int8
	continueFlag = 1
	var contact ContactsCache
	for continueFlag != 0 {
		baseDetailList, _ := SyncContactDetail(wxUser, currentWxcontactSeq, currentChatRoomContactSeq)
		currentWxcontactSeq, currentChatRoomContactSeq, continueFlag = baseDetailList.CurrentWxcontactSeq, baseDetailList.CurrentChatRoomContactSeq, baseDetailList.ContinueFlag
		for _, contactInfo := range baseDetailList.ContactList {
			if strings.Contains(contactInfo.UserName, "gh_") || strings.Contains("medianote,qqsafe,filehelper,newsapp,fmessage,weibo,qqmail,tmessage,qmessage,qqsync,weixin,floatbottle", contactInfo.UserName) {
				contact.GH = append(contact.GH, contactInfo)
			} else if strings.Contains(contactInfo.UserName, "@chatroom") {
				contact.ChatRoom = append(contact.ChatRoom, contactInfo)
			} else if strings.Contains(contactInfo.UserName, "system") {
				contact.System = append(contact.System, contactInfo)
			} else {
				contact.User = append(contact.User, contactInfo)
			}

		}
	}
	contactJson, _ := json.Marshal(contact)
	client_system.LogWriteData(string(contactJson))
	client_system.SetCache(wxUser+"_contact", &contact, 0)
}

/**
 * 同步全部基础通讯录
 */
func SyncContactAllBase(wxUser string) map[string]interface{} {
	var currentWxcontactSeq, currentChatRoomContactSeq int32
	var continueFlag int8
	continueFlag = 1
	var gh, user, chatroom, system []string
	for continueFlag != 0 {
		baseList, _ := SyncContact(wxUser, currentWxcontactSeq, currentChatRoomContactSeq)
		currentWxcontactSeq, currentChatRoomContactSeq, continueFlag = baseList.CurrentWxcontactSeq, baseList.CurrentChatRoomContactSeq, baseList.ContinueFlag
		for _, contactInfo := range baseList.UsernameLists {
			if strings.Contains(contactInfo.Username, "gh_") {
				gh = append(gh, contactInfo.Username)
			} else if strings.Contains(contactInfo.Username, "@chatroom") {
				chatroom = append(chatroom, contactInfo.Username)
			} else if strings.Contains("medianote,qqsafe,filehelper,newsapp,fmessage,weibo,qqmail,tmessage,qmessage,qqsync,weixin,floatbottle", contactInfo.Username) {
				system = append(system, contactInfo.Username)
			} else {
				user = append(user, contactInfo.Username)
			}
		}
	}
	contacts := map[string]interface{}{"gh": gh, "user": user, "chatroom": chatroom, "system": system}
	jsonStr, _ := json.Marshal(contacts)
	client_system.LogWriteData(string(jsonStr))
	return contacts
}

/**
 * 同步基础通讯录
 */
func SyncContact(wxUser string, currentUser int32, currentRoom int32) (ContactBase, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"CurrentWxcontactSeq": ` + strconv.Itoa(int(currentUser)) + `,
		"CurrentChatRoomContactSeq": ` + strconv.Itoa(int(currentRoom)) + `
	}`)
	)
	FastShortRequest(wxUser, CmdSyncContactBase, payload, "获取通讯录基础信息失败", &result, &err)
	if err != nil {
		return ContactBase{}, err
	}
	var contactBaseList ContactBase
	json.Unmarshal(result, &contactBaseList)
	return contactBaseList, nil
}

/**
 * 同步详情通讯录
 */
func SyncContactDetail(wxUser string, currentUser int32, currentRoom int32) (ContactBaseDetailList, error) {
	var (
		result []byte
		err    error
	)
	baseList, err := SyncContact(wxUser, currentUser, currentRoom)
	if err != nil {
		return ContactBaseDetailList{}, err
	}
	payload, _ := json.Marshal(baseList)
	FastShortRequest(wxUser, CmdSyncContactDetail, payload, "获取通讯录详情失败", &result, &err)
	if err != nil {
		return ContactBaseDetailList{}, err
	}
	var contactDetailList []ContactDetail
	json.Unmarshal(result, &contactDetailList)
	return ContactBaseDetailList{
		CurrentWxcontactSeq:       baseList.CurrentWxcontactSeq,
		CurrentChatRoomContactSeq: baseList.CurrentChatRoomContactSeq,
		ContinueFlag:              baseList.ContinueFlag,
		ContactList:               contactDetailList,
	}, nil
}

/**
 * 搜索微信用户
 */
func SearchContact(wxUser string, wxid string) (ContactDetail, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Username": "啊` + wxid + `"
	}`)
	)
	FastLongRequest(wxUser, CmdSearchContact, payload, "点击菜单操作失败", &result, &err)
	if err != nil {
		return ContactDetail{}, err
	}
	var contactInfo ContactDetail
	json.Unmarshal(result, &contactInfo)
	return contactInfo, nil
}

/**
 * 获取微信好友(群成员)信息列表
 */
func GetContact(wxUser string, userList string, chatRoom string) ([]ContactDetail, error) {
	var (
		result  []byte
		err     error
		payload = []byte(`{
		"Chatroomid": "` + chatRoom + `",
		"UserNameList": "` + userList + `"
	}`)
	)
	FastShortRequest(wxUser, CmdGetContact, payload, "获取用户信息失败", &result, &err)
	if err != nil {
		return []ContactDetail{}, err
	}
	var contactInfoList []ContactDetail
	json.Unmarshal(result, &contactInfoList)
	return contactInfoList, nil
}
