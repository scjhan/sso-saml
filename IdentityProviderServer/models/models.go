package models

import (
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/mysql"
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
)

func NotifyLogout(uid string) {
	hosts, _ := redis.SMembers(CreateRedisKey(uid, HostSetKey))
	for _, host := range hosts {
		q := url.Values{}
		q.Add("uid", uid)
		q.Add(proto.NotifyLabel, proto.NotifyLogout)
		u := url.URL{
			Scheme:   "http",
			Host:     host,
			Path:     "/idp_notify",
			RawQuery: q.Encode(),
		}

		client := http.Client{}
		client.Get(u.String())
	}
}

// CheckHost check host is legal
// host format example.com:port
func CheckHost(host string) bool {
	if len(host) == 0 {
		return false
	}

	o, err := mysql.NewMySQL()
	if err != nil {
		return false
	}

	lowerHost := strings.ToLower(host)
	sSql := fmt.Sprintf("select 1 from idp_cluster_info where host=%q", lowerHost)
	util.Debug(fmt.Sprintf("sql = %s", sSql))

	_, num := o.Query(sSql)
	if num <= 0 {
		return false
	}

	return true
}

type KeyType string

const (
	UIDSessionKey   = KeyType("US") // uid -> session
	SessionTokenKey = KeyType("ST") // session -> token
	TokenValueKey   = KeyType("TV") // token -> token_value
	HostSetKey      = KeyType("HS") // uid -> host_set
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
