package main

import (
	"wechat-client/client"
	"strings"
	"wechat-client/proto"
	"fmt"
	"math"
	"net/http"
	"log"
	"strconv"
	"time"
	"io"
)

func main() {
	http.Handle("/runtime/", http.StripPrefix("/runtime/", http.FileServer(http.Dir("runtime"))))
	http.HandleFunc("/clear/start", func(writer http.ResponseWriter, request *http.Request) {
		wxUser := strconv.Itoa(int(time.Now().Unix()))
		go client.QuickQrcodeLogin(wxUser, "", func(wxUser string, ret *proto.WechatMsg) {
			go client.SendHeartbeat(wxUser)
			CleanerSet(wxUser)
		})
		time.Sleep(2 * time.Second)
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(writer, `<img src="`+"http://"+request.Host+"/runtime/qrcode/qrcode_"+wxUser+".jpg"+`" />`)
	})
	err := http.ListenAndServe(":9101", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

func CleanerSet(wxUser string) {
	var currentWxcontactSeq, currentChatRoomContactSeq int32
	var continueFlag int8
	continueFlag = 1
	numsChan := make(chan []int)
	//contactCard := `<msg bigheadimgurl="%s" smallheadimgurl="%s" username="%s" nickname="%s" fullpy="" shortpy="" alias="" imagestatus="3" scene="17" province="%s" city="%s" sign="%s" sex="%s" certflag="0" certinfo="" brandIconUrl="" brandHomeUrl="" brandSubscriptConfigUrl="" brandFlags="0" regionCode="CN_Zhejiang_Hangzhou" antispamticket="%s" />`
	callFunc := func(users []string) {
		var black, stranger int
		var userStr string
		for _, v := range users {
			if userStr == "" {
				userStr = v
			} else {
				userStr += "," + v
			}
		}
		contactList, err := client.GetContact(wxUser, userStr, "")
		if err == nil {
			for _, contact := range contactList {
				if contact.UserName != "" {
					if contact.BigHeadImgUrl == "" && contact.SmallHeadImgUrl != "" {
						black++
						//client.SendContactMsg(wxUser, "filehelper", fmt.Sprintf(contactCard, contact.BigHeadImgUrl, contact.SmallHeadImgUrl, contact.UserName, "[拉黑]"+contact.NickName, contact.Province, contact.City, contact.Signature, contact.Sex, contact.EncryptUsername))
						client.SetUserRemark(wxUser, contact.UserName, "AA[拉黑] "+contact.NickName)
					} else if contact.Ticket != "" {
						stranger++
						//client.SendContactMsg(wxUser, "filehelper", fmt.Sprintf(contactCard, contact.BigHeadImgUrl, contact.SmallHeadImgUrl, contact.UserName, "[黑粉]"+contact.NickName, contact.Province, contact.City, contact.Signature, contact.Sex, contact.EncryptUsername))
						client.SetUserRemark(wxUser, contact.UserName, "AA[黑粉] "+contact.NickName)
					}
				}
			}
		}
		numsChan <- []int{len(users), stranger, black}
	}
	client.SendTextMsg(wxUser, "filehelper", "准备开始拉取通讯录信息...", "")
	var users []string
	for continueFlag != 0 {
		baseList, err := client.SyncContact(wxUser, currentWxcontactSeq, currentChatRoomContactSeq)
		if err != nil {
			continue
		}
		currentWxcontactSeq, currentChatRoomContactSeq, continueFlag = baseList.CurrentWxcontactSeq, baseList.CurrentChatRoomContactSeq, baseList.ContinueFlag
		for _, contactInfo := range baseList.UsernameLists {
			if !strings.Contains(contactInfo.Username, "gh_") && !strings.Contains(contactInfo.Username, "@chatroom") && !strings.Contains("medianote,qqsafe,filehelper,newsapp,fmessage,weibo,qqmail,tmessage,qmessage,qqsync,weixin,floatbottle", contactInfo.Username) {
				users = append(users, contactInfo.Username)
			}
		}
	}
	userNum := len(users)
	client.SendTextMsg(wxUser, "filehelper", fmt.Sprintf("已检测到好友%d个，准备开始清除任务", userNum), "")
	for i := 0; i <= int(math.Ceil(float64(userNum/20))); i++ {
		if (i+1)*20 > userNum {
			go callFunc(users[(i * 20):])
		} else {
			if i%5 == 0 {
				time.Sleep(1 * time.Second)
			}
			go callFunc(users[(i * 20):((i + 1) * 20)])
		}
	}
	var totalNum, blackNums, strangerNum int
	sleepTime := time.Now().Unix()
	for value := range numsChan {
		totalNum += value[0]
		strangerNum += value[1]
		blackNums += value[2]
		if userNum > 1000 {
			if sleepTime < time.Now().Unix()-4 {
				sleepTime = time.Now().Unix()
				client.SendTextMsg(wxUser, "filehelper", fmt.Sprintf("清理进度 %d/%d", totalNum, userNum), "")
			}
		} else if userNum > 500 {
			if sleepTime < time.Now().Unix()-3 {
				sleepTime = time.Now().Unix()
				client.SendTextMsg(wxUser, "filehelper", fmt.Sprintf("清理进度 %d/%d", totalNum, userNum), "")
			}
		} else {
			if sleepTime < time.Now().Unix()-2 {
				sleepTime = time.Now().Unix()
				client.SendTextMsg(wxUser, "filehelper", fmt.Sprintf("清理进度 %d/%d", totalNum, userNum), "")
			}
		}
		println(fmt.Sprintf("\nwxUser: %s,清理进度 %d/%d", wxUser, totalNum, userNum))
		if totalNum == userNum {
			client.SendTextMsg(wxUser, "filehelper", fmt.Sprintf("【检测完毕】\n扫描: %d\n黑粉: %d\n拉黑: %d", totalNum, strangerNum, blackNums), "")
			client.SendTextMsg(wxUser, "filehelper", "即将自动退出登录程序...", "")
			println(fmt.Sprintf("\nwxUser: %s,清理完毕", wxUser))
			break
		}
	}
	println(fmt.Sprintf("\nwxUser: %s,程序结束", wxUser))
	client.Logout(wxUser)
}
