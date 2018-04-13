package proto

type TokenVerifyData struct {
	Token string `string:"token"`
	Uid   string `json:"uid"`
	Name  string `json:"name"`
}

func (t *TokenVerifyData) Valid() bool {
	return len(t.Uid) != 0
}

type Session struct {
	SessionID string `json:"id"`
	UID       string `json:"uid"`
	Name      string `json:"name`
}
