package client

import (
	"wechat-client/client/system"
	"strconv"
	"time"
	"encoding/json"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/xml"
)

type MsgFriendRequestXml struct {
	Fromusername     string `xml:"fromusername,attr"`
	Encryptusername  string `xml:"encryptusername,attr"`
	Fromnickname     string `xml:"fromnickname,attr"`
	Content          string `xml:"content,attr"`
	Fullpy           string `xml:"fullpy,attr"`
	Shortpy          string `xml:"shortpy,attr"`
	Imagestatus      string `xml:"imagestatus,attr"`
	Scene            string `xml:"scene,attr"`
	Country          string `xml:"country,attr"`
	Province         string `xml:"province,attr"`
	City             string `xml:"city,attr"`
	Sign             string `xml:"sign,attr"`
	Percard          string `xml:"percard,attr"`
	Sex              string `xml:"sex,attr"`
	Alias            string `xml:"alias,attr"`
	Weibo            string `xml:"weibo,attr"`
	Albumflag        string `xml:"albumflag,attr"`
	Albumstyle       string `xml:"albumstyle,attr"`
	Albumbgimgid     string `xml:"albumbgimgid,attr"`
	Snsflag          string `xml:"snsflag,attr"`
	Snsbgimgid       string `xml:"snsbgimgid,attr"`
	Snsbgobjectid    string `xml:"snsbgobjectid,attr"`
	Mhash            string `xml:"mhash,attr"`
	Mfullhash        string `xml:"mfullhash,attr"`
	Bigheadimgurl    string `xml:"bigheadimgurl,attr"`
	Smallheadimgurl  string `xml:"smallheadimgurl,attr"`
	Ticket           string `xml:"ticket,attr"`
	Opcode           string `xml:"opcode,attr"`
	Googlecontact    string `xml:"googlecontact,attr"`
	Qrticket         string `xml:"qrticket,attr"`
	Chatroomusername string `xml:"chatroomusername,attr"`
	Sourceusername   string `xml:"sourceusername,attr"`
	Sourcenickname   string `xml:"sourcenickname,attr"`
	Brandlist        struct {
		Count int    `xml:"count,attr"`
		Ver   string `xml:"ver,attr"`
	} `xml:"brandlist"`
}

type MsgFaceXml struct {
	Attr struct {
		Fromusername      string `xml:"fromusername,attr"`
		Tousername        string `xml:"tousername,attr"`
		Type              int    `xml:"type,attr"`
		Idbuffer          string `xml:"idbuffer,attr"`
		Md5               string `xml:"md5,attr"`
		Len               int    `xml:"len,attr"`
		Productid         int    `xml:"productid,attr"`
		Androidmd5        string `xml:"androidmd5,attr"`
		Androidlen        int    `xml:"androidlen,attr"`
		S60v3md5          string `xml:"s60v3md5,attr"`
		S60v3len          int    `xml:"s60v3len,attr"`
		S60v5md5          string `xml:"s60v5md5,attr"`
		S60v5len          int    `xml:"s60v5len,attr"`
		Cdnurl            string `xml:"cdnurl,attr"`
		Designerid        string `xml:"designerid,attr"`
		Thumburl          string `xml:"thumburl,attr"`
		Encrypturl        string `xml:"encrypturl,attr"`
		Aeskey            string `xml:"aeskey,attr"`
		Externurl         string `xml:"externurl,attr"`
		Externmd5         string `xml:"externmd5,attr"`
		Width             int    `xml:"width,attr"`
		Height            int    `xml:"height,attr"`
		Tpurl             string `xml:"tpurl,attr"`
		Tpauthkey         string `xml:"tpauthkey,attr"`
		Attachedtext      string `xml:"attachedtext,attr"`
		Attachedtextcolor string `xml:"attachedtextcolor,attr"`
		Lensid            string `xml:"lensid,attr"`
	} `xml:"emoji"`
	GameExt struct {
		Type    int    `xml:"type"`
		Content string `xml:"content"`
	} `xml:"gameext"`
}

type MsgContactCardXml struct {
	XMLName         xml.Name `xml:"msg"`
	Bigheadimgurl   string   `xml:"bigheadimgurl,attr"`
	Smallheadimgurl string   `xml:"smallheadimgurl,attr"`
	Username        string   `xml:"username,attr"`
	Nickname        string   `xml:"nickname,attr"`
	Fullpy          string   `xml:"fullpy,attr"`
	Shortpy         string   `xml:"shortpy,attr"`
	Alias           string   `xml:"alias,attr"`
	Imagestatus     string   `xml:"imagestatus,attr"`
	Scene           string   `xml:"scene,attr"`
	Province        string   `xml:"province,attr"`
	City            string   `xml:"city,attr"`
	Sign            string   `xml:"sign,attr"`
	Certflag        string   `xml:"certflag,attr"`
	Certinfo        string   `xml:"certinfo,attr"`
	BrandIconUrl    string   `xml:"brandIconUrl,attr"`
	BrandHomeUrl    string   `xml:"brandHomeUrl,attr"`
	BrandFlags      string   `xml:"brandFlags,attr"`
	RegionCode      string   `xml:"regionCode,attr"`
	Antispamticket  string   `xml:"antispamticket,attr"`
}

type MsgVoiceXml struct {
	Attr struct {
		Endflag      int    `xml:"endflag,attr"`
		Cancelflag   int    `xml:"cancelflag,attr"`
		Forwardflag  int    `xml:"forwardflag,attr"`
		Voiceformat  int    `xml:"voiceformat,attr"`
		Voicelength  int32  `xml:"voicelength,attr"`
		Length       int    `xml:"length,attr"`
		Bufid        string `xml:"bufid,attr"`
		Clientmsgid  string `xml:"clientmsgid,attr"`
		Fromusername string `xml:"fromusername,attr"`
	} `xml:"voicemsg"`
}

type MsgVideoXml struct {
	Attr struct {
		Aeskey         string `xml:"aeskey,attr"`
		Cdnthumbaeskey string `xml:"cdnthumbaeskey,attr"`
		Cdnvideourl    string `xml:"cdnvideourl,attr"`
		Cdnthumburl    string `xml:"cdnthumburl,attr"`
		Length         int    `xml:"length,attr"`
		Playlength     int    `xml:"playlength,attr"`
		Cdnthumblength int    `xml:"cdnthumblength,attr"`
		Cdnthumbwidth  int    `xml:"cdnthumbwidth,attr"`
		Cdnthumbheight int    `xml:"cdnthumbheight,attr"`
		Fromusername   string `xml:"fromusername,attr"`
		Md5            string `xml:"md5,attr"`
		Newmd5         string `xml:"newmd5,attr"`
		Isad           int    `xml:"isad,attr"`
	} `xml:"videomsg"`
}

type MsgImgXml struct {
	Attr struct {
		AesKey         string `xml:"aeskey,attr"`
		Encryver       string `xml:"encryver,attr"`
		Cdnthumbaeskey string `xml:"cdnthumbaeskey,attr"`
		Cdnthumburl    string `xml:"cdnthumburl,attr"`
		Cdnthumblength int    `xml:"cdnthumblength,attr"`
		Cdnthumbheight int    `xml:"cdnthumbheight,attr"`
		Cdnthumbwidth  int    `xml:"cdnthumbwidth,attr"`
		Cdnmidheight   int    `xml:"cdnmidheight,attr"`
		Cdnmidwidth    int    `xml:"cdnmidwidth,attr"`
		Cdnhdheight    int    `xml:"cdnhdheight,attr"`
		Cdnhdwidth     int    `xml:"cdnhdwidth,attr"`
		Cdnmidimgurl   string `xml:"cdnmidimgurl,attr"`
		Length         int    `xml:"length,attr"`
		Md5            string `xml:"md5,attr"`
	} `xml:"img"`
}

type MsgLocationXml struct {
	Attr struct {
		X            string `xml:"x,attr"`
		Y            string `xml:"y,attr"`
		Scale        string `xml:"scale,attr"`
		Label        string `xml:"label,attr"`
		Maptype      int    `xml:"maptype,attr"`
		Poiname      int    `xml:"poiname,attr"`
		Poiid        int    `xml:"poiid,attr"`
		Fromusername int    `xml:"fromusername,attr"`
	} `xml:"location"`
}

type MsgAtXml struct {
	XMLName     xml.Name `xml:"msgsource"`
	AtUserList  string   `xml:"atuserlist"`
	Silence     string   `xml:"silence"`
	MemberCount string   `xml:"membercount"`
}

type CallbackMsg struct {
	Content      string
	CreateTime   int64
	FromUserName string
	ImgBuf       string
	ImgStatus    int8
	MsgId        int64
	MsgSource    string
	MsgType      int8
	NewMsgId     int64
	PushContent  string
	Status       int8
	ToUserName   string
}

//消息回调
func CallFunc(wxUser string, recv []byte) {
	cmd := client_system.ReadInt(recv, 8)
	if cmd != 24 && cmd != 318 {
		return
	}
	if cmd == 318 {
		print("\n收到新解析消息：")
	} else if cmd == 24 && len(recv) == 20 {
		//同步消息
		selector := client_system.ReadInt(recv, 16)
		if selector == -1 {
			//用户主动下线
		}
		if selector > 0 {
			syncMsgList := SyncMsg(wxUser)
			if syncMsgList == nil {
				return
			}
			for _, msgInfo := range syncMsgList {
				if msgInfo.FromUserName == "weixin" {
					continue
				}
				//设置缓存，5秒去重
				msgIdKey := strconv.Itoa(int(msgInfo.MsgId))
				_, found := client_system.GetCache(msgIdKey)
				if found {
					continue
				} else {
					client_system.SetCache(msgIdKey, 1, 5*time.Second)
				}
				//判断推送人条件开始推送
				if msgInfo.FromUserName != msgInfo.ToUserName && msgInfo.NewMsgId != 0 && msgInfo.MsgId != 0 {
					//用户自己操作的消息记录
					if strings.Contains(msgInfo.Content, "op=") {
						continue
					}
					msgJson, _ := json.Marshal(msgInfo)
					client_system.LogWriteData(string(msgJson))
					//处理群消息的前面的发送人字符串
					content := msgInfo.Content
					isChatroom := false
					fromUser := msgInfo.FromUserName
					fromChatroom := ""
					if strings.Contains(msgInfo.FromUserName, "@chatroom") {
						index := strings.Index(msgInfo.Content, ":\n")
						content = string([]byte(msgInfo.Content)[index+2:])
						fromUser = string([]byte(msgInfo.Content)[:index])
						fromChatroom = msgInfo.FromUserName
						isChatroom = true
					}
					msgContactInfo, _ := GetContact(wxUser, fromUser, fromChatroom)
					nicknameStr := msgInfo.FromUserName
					if len(msgContactInfo) > 0 {
						nicknameStr = msgContactInfo[0].NickName
					}
					if isChatroom {
						var atXml MsgAtXml
						if xml.Unmarshal([]byte(msgInfo.MsgSource), &atXml) != nil {
							continue
						}
						wxClient, err := GetWechatConn(wxUser)
						if err != nil {
							continue
						}
						if atXml.AtUserList != wxClient.VUser.GetUserame() {
							continue
						}
						contactInfo, err := GetContact(wxUser, wxClient.VUser.GetUserame(), "")
						if len(contactInfo) > 0 {
							atUserStr := "@" + contactInfo[0].NickName
							content = strings.Replace(content, atUserStr, "", 1)
						}
					} else {
						switch msgInfo.MsgType {
						case 1: //文字
							println(fmt.Sprintf("\n【%s】%s: %s", "文字", nicknameStr, content))
							switch content {
							case "删除好友":
								DelUser(wxUser, msgInfo.FromUserName)
								fmt.Println(SearchContact(wxUser, msgInfo.FromUserName))
								break
							case "设置黑名单":
								SetUserBlack(wxUser, msgInfo.FromUserName, true)
								break
							case "取消黑名单":
								SetUserBlack(wxUser, msgInfo.FromUserName, false)
								break
							case "设置标星":
								SetUserStar(wxUser, msgInfo.FromUserName, true)
								break
							case "取消标星":
								SetUserStar(wxUser, msgInfo.FromUserName, false)
								break
							case "设置置顶":
								SetUserTop(wxUser, msgInfo.FromUserName, true)
								break
							case "取消置顶":
								SetUserTop(wxUser, msgInfo.FromUserName, false)
								break
							case "设置备注":
								SetUserRemark(wxUser, msgInfo.FromUserName, "哈哈哈哈")
								break
							case "开启自动接受":
								SetAutoAcceptUser(wxUser, true)
								break
							case "关闭自动接受":
								SetAutoAcceptUser(wxUser, false)
								break
							case "同步朋友圈":
								SyncSns(wxUser)
								break
							case "获取我的朋友圈":
								SnsTimeLine(wxUser, "0")
								break
							case "获取你的朋友圈":
								GetUserSns(wxUser, msgInfo.FromUserName, "", 0)
								break
							case "上传图片":
								UploadSns(wxUser, []byte{})
								break
							case "发送朋友圈":
								SendSns(
									wxUser,
									fmt.Sprintf(`<TimelineObject><id>13078019682292478014</id><username>%s</username><createTime>1559021435</createTime><contentDesc>哈哈</contentDesc><contentDescShowType>0</contentDescShowType><contentDescScene>3</contentDescScene><private>0</private><sightFolded>0</sightFolded><showFlag>0</showFlag><appInfo><id></id><version></version><appName></appName><installUrl></installUrl><fromUrl></fromUrl><isForceUpdate>0</isForceUpdate></appInfo><sourceUserName></sourceUserName><sourceNickName></sourceNickName><statisticsData></statisticsData><statExtStr></statExtStr><ContentObject><contentStyle>2</contentStyle><title></title><description></description><mediaList></mediaList><contentUrl></contentUrl></ContentObject><actionInfo><appMsg><messageAction></messageAction></appMsg></actionInfo><location poiClassifyId="" poiName="" poiAddress="" poiClassifyType="0" city=""></location><publicUserName></publicUserName><streamvideo><streamvideourl></streamvideourl><streamvideothumburl></streamvideothumburl><streamvideoweburl></streamvideoweburl></streamvideo></TimelineObject>`, msgInfo.ToUserName),
								)
								break
							case "删除朋友圈":
								DelSns(wxUser, "13078019682292478014")
								break
							case "朋友圈设置隐私":
								SetSnsPrivacy(wxUser, "13078019682292478014", true)
								break
							case "朋友圈设置公开":
								SetSnsPrivacy(wxUser, "13078019682292478014", false)
								break
							case "点赞":
								SetSnsPraise(wxUser, "13078019682292478014", "wxid_t7p01dw592qt12", true)
								break
							case "取消点赞":
								SetSnsPraise(wxUser, "13078019682292478014", "wxid_t7p01dw592qt12", false)
								break
							case "评论朋友圈":
								CommentSns(wxUser, "13078019682292478014", "wxid_t7p01dw592qt12", "哈哈")
								break
							case "删除评论":
								DelCommentSns(wxUser, "13078019682292478014", "")
								break
							}
							go SendTextMsg(wxUser, msgInfo.FromUserName, content, "")
							break
						case 3: //图片
							println(fmt.Sprintf("\n【%s】%s: %s", "图片", nicknameStr, "[图片]"))
							var imgXml MsgImgXml
							xml.Unmarshal([]byte(content), &imgXml)
							//起个协程把这个图片存起来
							go DownloadMsgImage(wxUser, DownloadImg{
								MsgId:        msgInfo.MsgId,
								ToUsername:   msgInfo.ToUserName,
								StartPos:     0,
								TotalLen:     imgXml.Attr.Length,
								DataLen:      0,
								CompressType: 1,
								Md5:          imgXml.Attr.Md5,
							})
							//发送cdn图片
							SendCdnImageMsg(wxUser, SendCdnImageRequest{
								"",
								msgInfo.FromUserName,
								0,
								imgXml.Attr.Length,
								imgXml.Attr.Length,
								imgXml.Attr.Cdnmidimgurl,
								imgXml.Attr.AesKey,
								imgXml.Attr.Length,
								imgXml.Attr.Cdnthumblength,
								imgXml.Attr.Cdnthumbheight,
								imgXml.Attr.Cdnthumbwidth,
							})
							//content, _ := ioutil.ReadFile("/Users/Mr.Zhou/Project/golang/wechat-grpc-golang/test.jpg")
							//SendImageMsg(wxUser, msgInfo.FromUserName, content)
						case 34: //语音
							println(fmt.Sprintf("\n【%s】%s: %s", "语音", nicknameStr, "[语音]"))
							var voiceXml MsgVoiceXml
							xml.Unmarshal([]byte(content), &voiceXml)
							DownloadMsgVoice(wxUser, DownloadVoice{
								msgInfo.MsgId,
								0,
								0,
								voiceXml.Attr.Length,
								voiceXml.Attr.Clientmsgid,
								voiceXml.Attr.Bufid,
							})
							silkContent, _ := ioutil.ReadFile("./runtime/download/voice/voice_" + fmt.Sprintf("%d", msgInfo.MsgId) + ".silk")
							SendVoiceMsg(wxUser, msgInfo.FromUserName, silkContent, voiceXml.Attr.Voicelength)
							break
						case 37: //好友请求
							var friendRequestXml MsgFriendRequestXml
							xml.Unmarshal([]byte(content), &friendRequestXml)
							print(friendRequestXml.Ticket)
							err := AcceptUser(wxUser, friendRequestXml.Encryptusername, friendRequestXml.Ticket, friendRequestXml.Scene)
							if err != nil {
								print(err.Error())
							}
							err = SendTextMsg(wxUser, friendRequestXml.Fromusername, "我已经通过了你的好友申请，哈哈哈哈", "")
							println(fmt.Sprintf("\n【%s】%s: %s", "好友请求", nicknameStr, friendRequestXml.Fromnickname+"请求添加你为好友"))
							break
						case 42: //名片消息
							var contactCardXml MsgContactCardXml
							xml.Unmarshal([]byte(content), &contactCardXml)
							contactCardXml.Username = "huoniaojugege"
							contactStr, _ := xml.Marshal(contactCardXml)
							SendContactMsg(wxUser, msgInfo.FromUserName, string(contactStr))
							println(fmt.Sprintf("\n【%s】%s: %s", "好友请求", nicknameStr, "[名片]"+contactCardXml.Nickname))
							break
						case 43: //视频消息
							println(fmt.Sprintf("\n【%s】%s: %s", "视频", nicknameStr, "[视频]"))
							var videoXml MsgVideoXml
							xml.Unmarshal([]byte(content), &videoXml)
							SendCdnVideoMsg(wxUser, SendCdnVideoRequest{
								msgInfo.FromUserName,
								videoXml.Attr.Cdnthumblength,
								0,
								videoXml.Attr.Length,
								0,
								videoXml.Attr.Playlength,
								videoXml.Attr.Aeskey,
								videoXml.Attr.Cdnvideourl,
							})

							break
						case 47: //表情消息
							println(fmt.Sprintf("\n【%s】%s: %s", "表情", nicknameStr, "[表情]"))
							var faceXml MsgFaceXml
							xml.Unmarshal([]byte(content), &faceXml)
							//发送表情包
							SendFaceMsg(wxUser, SendFaceRequest{
								"",
								msgInfo.FromUserName,
								0,
								faceXml.Attr.Len,
								faceXml.Attr.Md5,
								`<gameext type=\"0\" content=\"0\" ></gameext>`,
							})
							break
						case 48: //定位消息
							println(fmt.Sprintf("\n【%s】%s: %s", "定位", nicknameStr, "[定位]"))
							var locationXml MsgLocationXml
							xml.Unmarshal([]byte(content), &locationXml)
							break
						case 49: //appmsg消息
							println(fmt.Sprintf("\n【%s】%s: %s", "卡片", nicknameStr, "[卡片]"))
							SendAppMsg(wxUser, msgInfo.FromUserName, fmt.Sprintf("<appmsg appid='' sdkver=''><title>%s</title><des>%s</des><action>view</action><type>5</type><showtype>0</showtype><content></content><url>%s</url><thumburl>%s</thumburl></appmsg>", "测试标题", "测试描述", "http://www.baidu.com", "https://ss2.bdstatic.com/70cFvnSh_Q1YnxGkpoWK1HF6hhy/it/u=466110605,1135815309&fm=200&gp=0.jpg"))
							//发送音乐链接
							//msg.SendAppMsg(msgInfo.FromUserName, `<appmsg appid='wx873a91b8917c375b' sdkver='0'><title>歌曲名字</title><des>歌手名字</des><type>3</type><showtype>0</showtype><soundtype>0</soundtype><contentattr>0</contentattr><url>http://i.y.qq.com/v8/playsong.html?songid=歌曲ID</url><lowurl>http://i.y.qq.com/v8/playsong.html?songid=歌曲ID</lowurl><dataurl>音乐地址</dataurl><lowdataurl>音乐地址</lowdataurl> <thumburl>歌曲图片</thumburl></appmsg>`)
							break
						case 50: //发送语音请求
							print(fmt.Sprintf("\n【收到语音请求消息】： %s", string(msgJson)))
							break
						case 51: //系统事件
							print(fmt.Sprintf("\n【收到系统事件消息】： %s", string(msgJson)))
							break
						case 52: //语音请求发送
							print(fmt.Sprintf("\n【收到发送语音请求消息】： %s", string(msgJson)))
							break
						case 53: //语音邀请
							break
						case 62: //小视频
							print(fmt.Sprintf("\n【收到小视频消息】： %s", string(msgJson)))
							break
						}
					}
				}
			}
		}
	}
}
