package DingTalk_go

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var (
	// Marshal defines function used to encode JSON payloads
	Marshal func(v interface{}) ([]byte, error) = json.Marshal
	// Unmarshal defines function used to decode JSON payloads
	Unmarshal func(src []byte, v interface{}) error = json.Unmarshal
)

type RequestConfig struct {
	Request                *http.Request
	ShouldRetryOnRateLimit bool
	MaxRestRetries         int
	Client                 *http.Client
}

func newRequestConfig(s *Session, req *http.Request) *RequestConfig {
	return &RequestConfig{
		ShouldRetryOnRateLimit: s.ShouldRetryOnRateLimit,
		MaxRestRetries:         s.MaxRestRetries,
		Client:                 s.Client,
		Request:                req,
	}
}

// RequestOption is a function which mutates request configuration.
// It can be supplied as an argument to any REST method.
type RequestOption func(cfg *RequestConfig)

// WithClient changes the HTTP client used for the request.
func WithClient(client *http.Client) RequestOption {
	return func(cfg *RequestConfig) {
		if client != nil {
			cfg.Client = client
		}
	}
}

// WithRetryOnRatelimit controls whether session will retry the request on rate limit.
func WithRetryOnRatelimit(retry bool) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.ShouldRetryOnRateLimit = retry
	}
}

// WithHeader sets a header in the request.
func WithHeader(key, value string) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request.Header.Set(key, value)
	}
}

func WithAccessToken(accessToken string) RequestOption {
	return WithHeader("x-acs-dingtalk-access-token", accessToken)
}

type accessTokenRequestBody struct {
	ClientID string `json:"client_id"`
	Token    string `json:"token"`
}

func (s *Session) getAccessToken() (token string, err error) {
	s.Lock()
	defer s.Unlock()
	storedTime := time.Unix(s.AccessTokenTimeStamps, 0)

	// 获取当前时间
	currentTime := time.Now()

	// 计算时间差
	diff := currentTime.Sub(storedTime)

	// 检查时间差是否超过1小时
	if diff.Hours() > 1 {
		// 超过1小时，重新获取token
		resp, errs := s.request(EndPointAccessToken, http.MethodPost, &accessTokenRequestBody{
			ClientID: s.ClientID,
			Token:    s.Token,
		})
		if errs != nil {
			return "", errs
		}
		var accessTokenResponse struct {
			AccessToken string `json:"accessToken"`
			ExpireIn    int64  `json:"expireIn"`
		}
		err = Unmarshal(resp, &accessTokenResponse)
		if err != nil {
			return
		}
		s.AccessToken = accessTokenResponse.AccessToken
		s.AccessTokenTimeStamps = time.Now().Unix()
		token = accessTokenResponse.AccessToken
	} else {
		// 未超过1小时，直接返回
		token = s.AccessToken
	}
	return
}

func (s *Session) MessageGroupSend(conversationID string, msg Message) (processQueryKey string, err error) {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return
	}
	var body struct {
		MsgParam           Message `json:"msgParam"`
		MsgKey             string  `json:"msgKey"`
		OpenConversationID string  `json:"openConversationId"`
		RobotCode          string  `json:"robotCode,omitempty"`
		Token              string  `json:"token,omitempty"`
		CoolAppCode        string  `json:"coolAppCode,omitempty"`
	}
	body.MsgParam = msg
	body.MsgKey = string(msg.Type())
	body.OpenConversationID = conversationID
	resp, err := s.request(EndPointGroupSend, http.MethodPost, body, WithAccessToken(accessToken))
	if err != nil {
		return
	}
	var response struct {
		ProcessQueryKey string `json:"processQueryKey"`
	}
	err = Unmarshal(resp, &response)
	if err != nil {
		return
	}
	processQueryKey = response.ProcessQueryKey
	return
}

func (s *Session) request(url string, method string, data interface{}, options ...RequestOption) (response []byte, err error) {
	var body []byte
	if data != nil {
		body, err = Marshal(data)
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	cfg := newRequestConfig(s, req)
	for _, opt := range options {
		opt(cfg)
	}
	req = cfg.Request
	resp, err := cfg.Client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
