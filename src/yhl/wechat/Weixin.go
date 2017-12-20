package wechat

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io"
	"sort"
	//	"strconv"
	"strings"
	"time"
	"yhl/help"
)

const (
	ApiUrl = "https://api.weixin.qq.com"
)

var (
	Token  string
	Appid  string
	Secret string
)

func init() {
	Token = beego.AppConfig.String("wechat.token")
	Appid = beego.AppConfig.String("wechat.appid")
	Secret = beego.AppConfig.String("wechat.secret")
}

type MsgBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

func Check(timestamp, nonce, signatureIn string) bool {
	sl := []string{Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))

	signatureGen := fmt.Sprintf("%x", s.Sum(nil))

	return signatureGen == signatureIn
}

func GetAccessToken() (token string) {
	cache := help.Cache
	t := cache.Get("access_token_" + Appid)
	if t != nil {
		token = string(t.([]uint8))
		return
	}

	url := ApiUrl + "/cgi-bin/token?grant_type=client_credential&appid=" + Appid + "&secret=" + Secret

	b := httplib.Get(url)
	data := make(map[string]interface{})
	b.ToJSON(&data)
	//fmt.Println(data)

	if v, ok := data["access_token"]; ok {
		token = v.(string)
		ttl := time.Duration(data["expires_in"].(float64))
		cache.Put("access_token_"+Appid, token, ttl*time.Second)
	}

	return
}

func SendMsg(m map[string]interface{}) {
	url := ApiUrl + "/cgi-bin/message/custom/send?access_token=" + GetAccessToken()
	req := httplib.Post(url)
	//fmt.Println(m)
	req.JSONBody(m)
	req.String()
	//fmt.Println(req.String())
}

func SendTextMsg(touser, content string) {
	m := map[string]interface{}{}
	m["touser"] = touser
	m["msgtype"] = "text"
	m["text"] = map[string]string{"content": content}

	SendMsg(m)
}

func GetWxUserinfo(openid, lang string) (m map[string]interface{}) {
	if lang == "" {
		lang = "zh_CN"
	}
	url := ApiUrl + "/cgi-bin/user/info?access_token=" + GetAccessToken() + "&openid=" + openid + "&lang=" + lang
	req := httplib.Get(url)
	m = make(map[string]interface{})
	req.ToJSON(&m)

	return
}