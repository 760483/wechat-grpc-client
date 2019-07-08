package client

import (
	"github.com/fastgoo/WechatApi/base"
	"encoding/json"
)

const (
	CmdGetLabelList = 639
	CmdAddLabel     = 635
	CmdSetUserLabel = 638
)

/**
 * 获取标签列表
 */
func GetLabels() interface{} {
	ret := base.FastCallUploadData(CmdGetLabelList, []byte{})
	print(string(ret.GetBaseMsg().GetPayloads()))
	var responseData interface{}
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &responseData)
	return responseData
}

/**
 * 新增标签
 */
func AddLabel(name string) interface{} {
	ret := base.FastCallUploadData(CmdAddLabel, []byte(`{
		"LabelName": "`+name+`"
	}`))
	print(string(ret.GetBaseMsg().GetPayloads()))
	var responseData interface{}
	json.Unmarshal(ret.GetBaseMsg().GetPayloads(), &responseData)
	return responseData
}

//不可用
func SetUserLabel(user string, labelIds string) {
	ret := base.FastCallUploadData(CmdSetUserLabel, []byte(`{
		"Username": "`+user+`",
		"Labelids": "`+labelIds+`"
	}`))
	print(string(ret.GetBaseMsg().GetPayloads()))
}
