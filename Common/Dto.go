package Common

type Message struct {
	Method string `json:"method"`
}

type CallLocalAddress struct {
	Message
	Payload struct {
		RequestId    string `json:"requestId"`
		ConnectionId int64  `json:"connectionId"`
		Address      string `json:"address"`
		Request      struct {
			Method  string                  `json:"method"`
			Headers map[string][]string     `json:"headers"`
			Body    *map[string]interface{} `json:"body"`
		} `json:"request"`
	} `json:"payload"`
}

type ResponseContent struct {
	Method  string `json:"method"`
	Payload struct {
		RequestId    string `json:"requestId"`
		ConnectionId int64  `json:"connectionId"`
		Content      struct {
			Type   string `json:"type"`
			Base64 string `json:"base64"`
		} `json:"content"`
	} `json:"payload"`
}

type Sharing struct {
	Method  string `json:"method"`
	Payload struct {
		LocalAddress string `json:"localAddress"`
	} `json:"payload"`
}

type SharingResponse struct {
	Method  string `json:"method"`
	Payload struct {
		Id            int64  `json:"id"`
		Key           string `json:"key"`
		RemoteAddress string `json:"remoteAddress"`
		LocalAddress  string `json:"localAddress"`
	} `json:"payload"`
}
