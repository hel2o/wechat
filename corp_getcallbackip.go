package wechat

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号获取企业微信服务器的ip段接口
const (
	WXAPI_Getcallbackip = WXAPI_ENT + "getcallbackip?access_token=%s"
)

// IPlist 企业微信服务器的ip段
type IPlist struct {
	Errcode int      `json:"errcode"`
	Errmsg  string   `json:"errmsg"`
	IPList  []string `json:"ip_list"`
}

// GetCallbackIP 通过Token获取服务器IP列表
func (s *Server) GetCallbackIP() (iPlist IPlist, err error) {
	url := fmt.Sprintf(WXAPI_Getcallbackip, s.GetAccessToken())
	if err = util.GetJson(url, &iPlist); err != nil {
		return
	}
	return
}

// GetCallbackIP 通过Token获取服务器IP列表
func GetCallbackIP() (iPlist IPlist, err error) {
	return std.GetCallbackIP()
}

func (s *Server) IsTrustIP(r *http.Request) (isTrust bool) {
	iplist, err := GetCallbackIP()
	if err != nil {
		log.Println("获取微信服务器的ip段失败", err, "Errcode", iplist.Errcode, "Errmsg", iplist.Errmsg)
		return
	}
	ip := r.RemoteAddr
	idx := strings.LastIndex(ip, ":")
	if idx > 0 {
		ip = ip[0:idx]
	}
	for _, ipRang := range iplist.IPList {
		if isTrust, _ = regexp.MatchString(ipRang, ip); isTrust {
			break
		}
	}
	return
}

func IsTrustIP(r *http.Request) (isTrust bool) {
	return std.IsTrustIP(r)

}
