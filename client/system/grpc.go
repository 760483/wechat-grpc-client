package client_system

import (
	pb "wechat-client/proto"
	"strings"
	"context"
	"google.golang.org/grpc"
	"time"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
)

const (
	connTLS = true
	address     = ""
	appid       = ""
	appkey      = ""
	machineCode = ""
	version     = "6.7.4"
	clientIp    = "127.0.0.1"
	sixData     = ""
	longServer  = "long.weixin.qq.com"
	shortServer = "short.weixin.qq.com"
	serverName  = "wechat@root"
	serverCrt   = `-----BEGIN CERTIFICATE-----
MIIDqjCCApKgAwIBAgIJAOf7+/avi9foMA0GCSqGSIb3DQEBCwUAMGoxCzAJBgNV
BAYTAkNOMQswCQYDVQQIDAJ4eDELMAkGA1UEBwwCeHgxCzAJBgNVBAoMAnh4MQsw
CQYDVQQLDAJ4eDEUMBIGA1UEAwwLd2VjaGF0QHJvb3QxETAPBgkqhkiG9w0BCQEW
Anh4MB4XDTE4MDUwODAwMzE1MFoXDTI4MDUwNTAwMzE1MFowajELMAkGA1UEBhMC
Q04xCzAJBgNVBAgMAnh4MQswCQYDVQQHDAJ4eDELMAkGA1UECgwCeHgxCzAJBgNV
BAsMAnh4MRQwEgYDVQQDDAt3ZWNoYXRAcm9vdDERMA8GCSqGSIb3DQEJARYCeHgw
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDNcWVlhOpXjqivCpJpVHCG
eM+Q/e83MDpqXmlT33hAFlCfUzWYYhFAHwc+xShkvuoaBwKpAr0fcHT7Kj4TLSJB
F9FwBccPp4Tv0YH0h4pzZ9nWZCMGXB5TlXAreXrcu8Qab0MPAxMKtMO0FZzQawRD
mO/S43u6tDBrhW2zgFzCUo+cRiCYRXoewuFVT4RT7eYvEu579oy4mto7YpR8GfPY
ywiTh7D/4To7nNoWbly84WB9ZSQ8HZoC1KTykqaeSGw7xUhKEoPoeKPke1ECfXB2
lSZMwKar2IU0BVlvxL55EMF8oXJozLVEZeVIZ4Wp9EFY2KaZJowgoOa1I+vcwBbd
AgMBAAGjUzBRMB0GA1UdDgQWBBR2SmMGUpTN/Dl0wcSYcPav6o01iDAfBgNVHSME
GDAWgBR2SmMGUpTN/Dl0wcSYcPav6o01iDAPBgNVHRMBAf8EBTADAQH/MA0GCSqG
SIb3DQEBCwUAA4IBAQAdCSkska27VdLcGqK9/sraopxaX31Nseci/sJbimIHxr+q
DwAHExU5sJ1qT827n2OpF/lWMRhnJZ3ubeJ8oGA3CAKu4EiKKGA1hGOLCbEvagCc
sdBSegk050qkMssJzNaw7boZB8vek1RDK32Fuhsh4m+MUZBj6bJCdGW9K+ZMmpZl
bMwmsgqV6+EMvr+PhFHy8bOAdIs4/eOTjW7R0JwYgFArVXMrVKgiRknkhM+PBBHG
DPWO0j3855SF2X5r4jQs2PvKGJjOMuQeIgsf2GbwSQhXEhM8lGdjn9up8hm7VSXf
x8wZquXczPSdDez7tP+g9nKbxcJtGnxo8+Jntmvs
-----END CERTIFICATE-----`
)

type customCredential struct{}

func CreateWechatMsg(cmd int32, payloads []byte) *pb.WechatMsg {
	randomEncryKey := []byte{80, 117, 128, 85, 2, 55, 180, 126, 141, 93, 185, 220, 112, 142, 15, 128}
	return &pb.WechatMsg{
		Token:     machineCode,
		Version:   version,
		TimeStamp: int32(time.Now().Unix()),
		Ip:        clientIp,
		BaseMsg: &pb.BaseMsg{
			//Ret: "",
			Cmd: cmd,
			//CmdUrl:         "",
			//ShortHost:      "",
			//LongHost:       "",
			//LongHead:       byte(),
			Payloads: payloads,
			User: &pb.User{
				//Uin:            1,
				//Cookies:        byte(),
				SessionKey: randomEncryKey,
				DeviceId:   sixData,
				//DeviceId:   "5326451F200E0D130CE4AE27262B6169",
				//DeviceType:     "",
				//DeviceName:     "",
				//CurrentsyncKey: byte(),
				//MaxSyncKey:     byte(),
				//AutoAuthKey:    byte(),
				//Userame:        "",
				//Nickname:       byte(),
				//UserExt:        byte(),
			},
		},
	}
}

func HelloWechat(wechatMsg *pb.WechatMsg) (*pb.WechatMsg, error) {
	var conn *grpc.ClientConn
	var err error
	if connTLS {
		conn, err = createTLSClient()
	} else {
		conn, err = createClient()
	}
	//conn, err := poolGrpc.Get(context.Background())
	//
	if err != nil {
		LogWrite(LOG_ERROR, "创建grpc客户端失败，error: "+err.Error())
		return nil, err
	}
	defer conn.Close()
	c := pb.NewWechatClient(conn)
	msg, err := c.HelloWechat(context.Background(), wechatMsg)
	if err != nil {
		if strings.Contains(err.Error(), "connection") {
			//如果包含connection字段，那说明可能grpc连接失败了，就重试一下
			LogWrite(LOG_ERROR, "连接grpc服务器失败，error: "+err.Error())
		} else if strings.Contains(err.Error(), "SessionTimeOut") {

			//如果包含SessionTimeOut字段，那么说明登录信息失效了，需要重新登录一下
			LogWrite(LOG_ERROR, "用户session失效，error: "+err.Error())
		} else if strings.Contains(err.Error(), "transport is closing") {
			//连接正在关闭，需要重新处理一下连接操作
			LogWrite(LOG_ERROR, "grpc连接正在关闭，error: "+err.Error())
		} else {
			LogWrite(LOG_ERROR, "error: "+err.Error())
		}
		return nil, err
	}
	return msg, nil
}

func createTLSClient() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	//roots := x509.NewCertPool()
	//roots.AppendCertsFromPEM([]byte(serverCrt))
	//creds := credentials.NewClientTLSFromCert(roots, serverName)
	//opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithTransportCredentials(
		credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		}),
	))
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	return grpc.Dial(address, opts...)
}

func createClient() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	return grpc.Dial(address, opts...)
}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  appid,
		"appkey": appkey,
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return connTLS
}
