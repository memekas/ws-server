package ws_server

type Notification struct {
	ToUser   uint   `json:"toUser"`
	FromUser uint   `json:"-"`
	Msg      string `json:"msg"`
}

type UsersOnline struct {
	Count uint32 `json:"count"`
}
