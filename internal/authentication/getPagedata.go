package authentication

import (
	"io"
	"net/http"
	"strconv"

	"github.com/antchfx/htmlquery"
	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
)

type UserData struct {
	UserId    int
	Xsrftoken string
}

func GetUserData(cookies []*kooky.Cookie) (UserData, error) {
	reqUrl := "https://app.timetastic.co.uk/wallchart"

	req, _ := http.NewRequest("GET", reqUrl, nil)
	for _, cookie := range cookies {
		req.AddCookie(&cookie.Cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserData{}, err
	}
	defer resp.Body.Close()
	userdata, err := extractUserData(resp.Body)

	return userdata, err
}

func extractUserData(body io.Reader) (UserData, error) {
	data := UserData{}
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return data, err
	}

	xsrfnodes := htmlquery.Find(doc, "//*[@id=\"_AjaxCsrfToken\"]")
	for _, attr := range xsrfnodes[0].Attr {
		if attr.Key == "value" {
			data.Xsrftoken = attr.Val
		}
	}

	usernodes := htmlquery.Find(doc, "//*[@id=\"logoutForm\"]")
	for _, attr := range usernodes[0].Attr {
		if attr.Key == "data-id" {
			i, err := strconv.Atoi(attr.Val)
			if err != nil {
				return data, err
			}
			data.UserId = i
		}
	}

	return data, nil
}
