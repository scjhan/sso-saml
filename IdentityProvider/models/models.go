package models

import (
	"bytes"
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/log"
	"chenjunhan/sso-saml/utils/mysql"
	"chenjunhan/sso-saml/utils/redis"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
)

const (
	Page403 = "403.tpl"
	Page500 = "500.tpl"
)

const (
	MagicExpire = 5 * 60
)

func GetHost() string {
	host := beego.AppConfig.String("host")
	if len(host) != 0 {
		return host
	}
	return "idp.com"
}

func GetPort() string {
	port := beego.AppConfig.String("httpport")
	if len(port) != 0 {
		return port
	}
	return "9090"
}

func GetHostPort() string {
	return GetHost() + ":" + GetPort()
}

func GetClusterHostPort() map[string]string {
	hosts := strings.Split(beego.AppConfig.String("cluster::hosts"), beego.AppConfig.String("cluster::commas"))
	ports := strings.Split(beego.AppConfig.String("cluster::ports"), beego.AppConfig.String("cluster::commas"))

	retval := map[string]string{}
	for i, v := range hosts {
		retval[v] = v + ":" + ports[i]
	}

	return retval
}

func NotifyLogout(uid string) {
	hosts, _ := redis.SMembers(CreateRedisKey(uid, HostSetKey))
	hp := GetClusterHostPort()

	fmt.Println(hosts, hp)

	for _, v := range hosts {
		go func(host string) {
			u := url.URL{
				Scheme: "http",
				Host:   hp[host],
				Path:   "/push",
			}
			c := http.Client{}

			msg := proto.PushMsg{
				Type:    proto.IdpLogout,
				Content: uid,
			}

			mb, _ := json.Marshal(msg)

			fmt.Println("=================NotifyLogout url =", u.String())

			c.Post(u.String(), "application/json", bytes.NewReader(mb))
		}(v)
	}
}

// CheckHost check host is legal
// host format example.com:port
func CheckHost(host string) (string, bool) {
	if len(host) == 0 {
		return "", false
	}

	o, err := mysql.NewMySQL()
	if err != nil {
		return "", false
	}

	lowerHost := strings.ToLower(host)
	sSql := fmt.Sprintf("select 1 from idp_cluster_info where host=%q", lowerHost)
	log.Debug(fmt.Sprintf("sql = %s", sSql))

	_, num := o.Query(sSql)
	if num <= 0 {
		return "", false
	}

	return host, true
}

type KeyType string

const (
	UIDSessionKey   = KeyType("US") // uid -> session
	SessionTokenKey = KeyType("ST") // session -> token
	TokenValueKey   = KeyType("TV") // token -> token_value
	HostSetKey      = KeyType("HS") // uid -> host_set
	MagicKey        = KeyType("MK") // magic -> host
)

func CreateRedisKey(key string, kt KeyType) string {
	return beego.AppConfig.String("appname") + "_" + string(kt) + "_" + key
}

func DeleteUIDCache(uid string) {
	k1 := CreateRedisKey(uid, UIDSessionKey)
	sessionid, _ := redis.GetString(k1)
	if len(sessionid) != 0 {
		k2 := CreateRedisKey(sessionid, SessionTokenKey)
		token, _ := redis.GetString(k2)
		if len(token) != 0 {
			redis.Delete(CreateRedisKey(token, TokenValueKey)) // delete TokenValueKey
		}
		redis.Delete(k2) // delete SessionTokenKey
	}

	redis.Delete(k1) // delete UIDSessionKey

	redis.Delete(CreateRedisKey(uid, HostSetKey)) // delete HostSetKey
}
