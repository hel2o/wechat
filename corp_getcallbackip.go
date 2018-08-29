package wechat

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

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
	s.serverIPListLocker.Lock()
	defer s.serverIPListLocker.Unlock()
	url := fmt.Sprintf(WXAPI_Getcallbackip, s.GetAccessToken())
	if err = util.GetJson(url, &iPlist); err != nil {
		return
	}
	return
}

// FetchWxIPList 定期获取服务器IP列表
func (s *Server) FetchServerIPList() {
	i := 0
	go func() {
		for {
			if s.SyncServerIPList() != nil && i < 2 {
				i++
				Println("尝试再次获取服务器IP列表(", i, ")")
				continue
			}
			i = 0
			time.Sleep(FetchDelay)
		}
	}()
}

// SyncServerIPList 同步服务器IP列表
func (s *Server) SyncServerIPList() (err error) {
	var iplist IPlist
	iplist, err = s.GetCallbackIP()
	if err != nil {
		Println("获取微信服务器的ip段失败", err, "Errcode", iplist.Errcode, "Errmsg", iplist.Errmsg)
		return
	}
	s.serverIPList = &iplist
	return
}

func (s *Server) IsTrustIP(r *http.Request) (isTrust bool) {
	ip := r.RemoteAddr
	idx := strings.LastIndex(ip, ":")
	if idx > 0 {
		ip = ip[0:idx]
	}
	for _, ipRang := range s.serverIPList.IPList {
		if isTrust, _ = regexp.MatchString(ipRang, ip); isTrust {
			break
		}
	}
	return
}

// GetCallbackIP 通过Token获取服务器IP列表
func GetCallbackIP() (iPlist IPlist, err error) {
	return std.GetCallbackIP()
}

func IsTrustIP(r *http.Request) (isTrust bool) {
	return std.IsTrustIP(r)
}
