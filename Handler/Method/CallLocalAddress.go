package Method

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gport/Common"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func compressData(data string) (string, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write([]byte(data))
	if err != nil {
		return "", err
	}

	_ = writer.Close()

	compressedData := buffer.Bytes()
	return base64.StdEncoding.EncodeToString(compressedData), nil
}

func statusColor(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return "\x1b[32m"
	} else if statusCode >= 300 && statusCode < 400 {
		return "\x1b[38;5;247m"
	} else if statusCode >= 400 && statusCode < 500 {
		return "\x1b[33m"
	} else {
		return "\x1b[31m"
	}
}

func CallLocalAddress(rawMessage []byte, writeChan chan []byte) {
	var message Common.CallLocalAddress

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		fmt.Println("Cannot decode message:", err)
		return
	}

	requestBody, reader := message.Payload.Request.Body, new(bytes.Buffer)

	err := json.NewEncoder(reader).Encode(requestBody)

	if err != nil {
		fmt.Println("Error encoding requestBody", err)
		return
	}

	request, err := http.NewRequest(message.Payload.Request.Method, message.Payload.Address, reader)

	client := http.Client{}

	response, err := client.Do(request)

	u, _ := url.Parse(message.Payload.Address)

	if err != nil {
		fmt.Println("\x1b[31m", "Error", u.Path, "-", err, "\x1b[0m")
		return
	}

	fmt.Println(strings.ToUpper(message.Payload.Request.Method), statusColor(response.StatusCode), response.Status, u.Path, "\x1b[0m")

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	compressedContent, err := compressData(string(body))

	if err != nil {
		fmt.Println("Error compressing data:", err)
		return
	}

	callback := Common.ResponseContent{}

	callback.Method = "remote-client:response"
	callback.Payload.RequestId = message.Payload.RequestId
	callback.Payload.ConnectionId = message.Payload.ConnectionId
	callback.Payload.Content.Type = response.Header.Get("Content-Type")
	callback.Payload.Content.Base64 = compressedContent

	callbackMessage, err := json.Marshal(callback)

	if err != nil {
		fmt.Println("Error marshalling callback message:", err)
		return
	}

	writeChan <- callbackMessage
}
