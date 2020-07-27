package ws_server

type Notification struct {
	ToUser   uint   `json:"toUser"`
	FromUser uint   `json:"-"`
	Msg      string `json:"msg"`
}
