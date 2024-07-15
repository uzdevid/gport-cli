package Common

type Message struct {
	Method string `json:"method"`
}

type CallLocalAddress struct {
	Message
	Payload struct {
		TrafficId int64  `json:"trafficId"`
		ClientId  string `json:"clientId"`
		Request   struct {
			Url     string                  `json:"url"`
			Method  string                  `json:"method"`
			Headers map[string][]string     `json:"headers"`
			Body    *map[string]interface{} `json:"body"`
		} `json:"request"`
	} `json:"payload"`
}

type ResponseContent struct {
	Method  string `json:"method"`
	Payload struct {
		TrafficId int64  `json:"trafficId"`
		ClientId  string `json:"clientId"`
		Response  struct {
			Headers map[string][]string `json:"headers"`
			Body    string              `json:"body"`
		} `json:"response"`
	} `json:"payload"`
}

type Sharing struct {
	Method  string `json:"method"`
	Payload struct {
		Domain       string `json:"domain"`
		LocalAddress string `json:"localAddress"`
	} `json:"payload"`
}

type SharingResponse struct {
	Method  string `json:"method"`
	Payload struct {
		Id      int64             `json:"id"`
		Key     string            `json:"key"`
		Proxies map[string]string `json:"proxies"`
	} `json:"payload"`
}

type PrintMessage struct {
	Message
	Payload struct {
		Message string `json:"message"`
		Exit    bool   `json:"exit"`
	}
}
