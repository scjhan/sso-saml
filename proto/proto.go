package proto

import (
	"encoding/json"
)

type TokenVerifyData struct {
	Token string `string:"token"`
	Uid   string `json:"uid"`
	Name  string `json:"name"`
}

func (t *TokenVerifyData) Valid() bool {
	return len(t.Uid) != 0
}

// Session Sp session struct
type Session struct {
	SessionID string `json:"id"`
	UID       string `json:"uid"`
	Name      string `json:"name`
}

// Push message define
const (
	Error = 0
	Ok    = 1
	// from idp to cluster, must litter than TypeCommas
	IdpLogout = 1

	TypeCommas = 255

	// frome cluster to idp, must bigger than TypeCommas
	ClusterLogout      = 256
	ClusterVerifyToken = 257
)

// PushMsg push message struct
type PushMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

func (m *PushMsg) String() string {
	s, _ := json.Marshal(*m)

	return string(s)
}

func ToPushMsg(msg []byte) PushMsg {
	pm := PushMsg{
		Type: Error,
	}

	json.Unmarshal(msg, &pm)

	return pm
}

type cryptoMsg struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

// EncryptString use des algorithm to encrypt the pushmsg,
// and return a json string which represent cryptoMsg value,
// cryptoMsg.Key means the des key string which encrypt by rsa_pub
func (m *PushMsg) EncryptString(pub, priv string) (string, error) {

	return "", nil
}
