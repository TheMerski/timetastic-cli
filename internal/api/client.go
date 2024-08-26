package api

import (
	"log/slog"
	"net/http"

	"github.com/browserutils/kooky"
	"github.com/themerski/timetastic-cli/internal/authentication"
)

type TimetasticClient struct {
	cookies  []*kooky.Cookie
	userData authentication.UserData
}

func NewTimetasticClient() *TimetasticClient {
	cookies, err := authentication.GetTimetasticCookiesFromChrome()
	if err != nil {
		slog.Error("Failed to get cookies", "error", err)
		return nil
	}

	userData, err := authentication.GetUserData(cookies)
	if err != nil {
		slog.Error("Failed to get user data", "error", err)
		return nil
	}

	return &TimetasticClient{
		cookies:  cookies,
		userData: userData,
	}
}

func (c *TimetasticClient) makeRequest(req *http.Request) (*http.Response, error) {
	for _, cookie := range c.cookies {
		req.AddCookie(&cookie.Cookie)
	}

	req.Header.Set("x-xsrf-token", c.userData.Xsrftoken)
	client := &http.Client{}
	return client.Do(req)
}
