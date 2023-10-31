package env

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	ErrInvalidInput        = errors.New("Invalid input")
	ErrBadRequest          = errors.New("Bad request")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrTooManyRequests     = errors.New("Too many requests")
	ErrServiceUnavailable  = errors.New("Service unavailable")
	ErrInternalServerError = errors.New("Internal server error")
)

// ================================================================
//
// ================================================================
type TwSms struct {
	Username string
	Password string
	URL      *url.URL
}

func NewTwSms() (*TwSms, error) {
	u, err := url.Parse("https://api.twsms.com/json/sms_send.php")
	if err != nil {
		return nil, err
	}

	e := &TwSms{
		Username: os.Getenv("TWSMS_USERNAME"),
		Password: os.Getenv("TWSMS_PASSWORD"),
	}

	q := u.Query()
	q.Set("username", e.Username)
	q.Set("password", e.Password)
	u.RawQuery = q.Encode()
	e.URL = u

	return e, nil
}

type TwSmsSendApiResp struct {
	Code  string `json:"code"`
	Text  string `json:"text"`
	Msgid int64  `json:"msgid"`
}

func (r TwSmsSendApiResp) Error() error {
	if code, err := strconv.Atoi(r.Code); err != nil {
		return ErrServiceUnavailable
	} else {
		switch {
		case code <= 1:
			return nil
		case code >= 10 && code <= 12:
			return ErrBadRequest
		case code == 20:
			return ErrTooManyRequests
		case code >= 30 && code <= 41:
			return ErrUnauthorized
		case code >= 50 && code <= 140:
			return ErrBadRequest
		case code >= 99998:
			return ErrServiceUnavailable
		default:
			return ErrInternalServerError
		}
	}
}

func (e TwSms) SendSms(to []string, subject, body string) error {
	if len(to) != 1 {
		return ErrInvalidInput
	}
	if subject != "" {
		body = subject + body
	}

	q := e.URL.Query()
	q.Set("mobile", to[0])
	q.Set("message", body)
	e.URL.RawQuery = q.Encode()

	if req, err := http.NewRequest("POST", e.URL.String(), nil); err != nil {
		return ErrInternalServerError
	} else {
		client := &http.Client{}
		if resp, err := client.Do(req); err != nil {
			return ErrInternalServerError
		} else {
			defer resp.Body.Close()
			switch {
			case resp.StatusCode >= 500:
				return ErrServiceUnavailable
			case resp.StatusCode >= 400:
				return ErrInternalServerError
			}

			if body, err := io.ReadAll(resp.Body); err != nil {
				return ErrInternalServerError
			} else {
				apiresp := new(TwSmsSendApiResp)
				if err := json.Unmarshal(body, apiresp); err != nil {
					return ErrServiceUnavailable
				} else {
					return apiresp.Error()
				}
			}
		}
	}
}
