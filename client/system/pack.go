package client_system

import (
	"bytes"
	"encoding/binary"
	"net"
	"fmt"
	"strconv"
)

var AsyncCallFunc func(wxUser string, buffer []byte)

//封包协程
func AsyncRecvWxData(wxUser string, wxConn net.Conn) {

	newBuff := new(bytes.Buffer)
	buf := make([]byte, 1024)
	var totalLength, size int32
	unpack := func(unpackBuf []byte) (int, []byte) {
		if totalLength == 0 {
			totalLength = ReadInt(unpackBuf, 0)
			LogWrite(LOG_INFO, "处理粘包：获取粘包的包长度->"+strconv.Itoa(int(totalLength)))
		}
		n := len(unpackBuf)
		size += int32(n)
		LogWrite(LOG_INFO, "处理粘包：已读取的包总长->"+strconv.Itoa(int(size)))
		if totalLength == size {
			LogWrite(LOG_INFO, "处理粘包：刚好完成粘包的封包操作")
			newBuff.Write(unpackBuf)
			CompleteBuffer(wxUser, newBuff.Bytes())
			totalLength = 0
			size = 0
			newBuff.Reset()
			return 0, nil
		} else if totalLength < size {
			LogWrite(LOG_INFO, "处理粘包：处理完成一个包，发现还有多的粘包...继续循环处理")
			endPost := totalLength - (size - int32(n))
			newBuff.Write(unpackBuf[:endPost])
			CompleteBuffer(wxUser, newBuff.Bytes())
			newBuff.Reset()
			totalLength = 0
			size = 0
			return 1, unpackBuf[endPost:]
		} else {
			newBuff.Write(unpackBuf)
			LogWrite(LOG_INFO, "处理粘包：未处理完一个包，跳出循环进行下一波的缓冲读取处理")
			return -1, nil
		}
	}
	for {
		LogWrite(LOG_INFO, "开始封包：↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
		n, err := wxConn.Read(buf)
		if err != nil {
			LogWrite(LOG_ERROR, "封包服务器已被微信断开，error: "+err.Error())
			break
		}
		if totalLength == 0 {
			totalLength = ReadInt(buf, 0)
			if totalLength < 0 {
				continue
			}
			LogWrite(LOG_INFO, "获取到包头大小: "+fmt.Sprintf("%d", totalLength))
		}
		size += int32(n)
		if totalLength == size {
			LogWrite(LOG_INFO, "结束封包：↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑")
			newBuff.Write(buf[:n])
			CompleteBuffer(wxUser, newBuff.Bytes())
			totalLength = 0
			size = 0
			newBuff.Reset()
			continue
		} else if size > totalLength {
			LogWrite(LOG_INFO, "粘包了："+fmt.Sprintf("\n封包数据过多，需要拆解。\n总计长度：%d\n上次长度：%d\n本次取到数据：%d, 多出了长度：%d, 需要截取保存长度：%d", totalLength, size-int32(n), n, size-totalLength, int32(n)-(size-totalLength)))
			LogWrite(LOG_INFO, "结束封包：↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑")
			if totalLength < 0 {
				continue
			}
			endPost := totalLength - (size - int32(n))
			newBuff.Write(buf[:endPost])
			CompleteBuffer(wxUser, newBuff.Bytes())
			newBuff.Reset()
			totalLength = 0
			size = 0
			continue
			unCode, unBuf := unpack(buf[endPost:])
			if unCode == 1 {
				for unCode == 1 {
					unCode, unBuf = unpack(unBuf)
				}
			} else {
				continue
			}
		} else {
			newBuff.Write(buf[:n])
		}
	}
}

//完成封包，回调通知事件
func CompleteBuffer(wxUser string, buffer []byte) {
	cmd := ReadInt(buffer, 8)
	SetWechatBufferCache(buffer)
	if (cmd == 24 && len(buffer) == 20) || cmd == 318 {
		if AsyncCallFunc != nil {
			//这里一定要异步处理，不然会造成上层封包阻塞
			go AsyncCallFunc(wxUser, buffer)
		}
	}
}

//读取包头
func ReadInt(buf []byte, index int) int32 {
	if len(buf) >= 16 {
		buffer := bytes.NewBuffer(buf[index : index+4])
		var size int32
		binary.Read(buffer, binary.BigEndian, &size)
		return size
	}
	return 1
}
