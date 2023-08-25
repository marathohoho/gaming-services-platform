package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Headers map[string][]string

type HttpResponse struct {
	StatusCode int
	Headers    Headers
	Body       []byte
}

var ErrUnsuccessfullHttpStatusCode = errors.New("unsuccessful http status code")
var ErrResponseBodyUnmarshalling = errors.New("failed to unmarshal response body")
var ErrResponseBodyRead = errors.New("failed to read response body")

func SendJsonRequest(method, url string, payload []byte, headers map[string]string) (response *HttpResponse, err error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	headers[fiber.HeaderContentType] = fiber.MIMEApplicationJSON
	response, err = sendRaw(method, url, payload, headers)
	if err != nil {
		if errors.Is(err, ErrUnsuccessfullHttpStatusCode) {
			err = nil
		}
	}

	return
}

func sendRaw(method, url string, payload []byte, headers map[string]string) (resp *HttpResponse, err error) {
	if method == "" {
		err = errors.New("method argument must not be an empty string")
		return
	}
	method = strings.ToUpper(method)
	if method == http.MethodGet && payload != nil {
		err = errors.New("invalid 'method' argument value")
		return
	}
	if url == "" {
		err = errors.New("url argument must not be an empty string")
		return
	}

	request, err := createRequest(method, url, bytes.NewBuffer(payload), headers)
	if err != nil {
		return
	}

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	resp = &HttpResponse{
		StatusCode: response.StatusCode,
		Headers:    Headers{},
	}

	for headerName, headerValues := range response.Header {
		headerNameInLower := strings.ToLower(headerName)
		resp.Headers[headerNameInLower] = headerValues
	}

	resp.Body, err = io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrResponseBodyUnmarshalling, err.Error())
		return
	}

	if response.StatusCode >= 400 {
		err = ErrUnsuccessfullHttpStatusCode
		return
	}

	return
}

func createRequest(method, url string, body io.Reader, headers map[string]string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	for headerName, headerValue := range headers {
		request.Header.Set(headerName, headerValue)
	}

	return
}
