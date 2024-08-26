package authentication

import (
	"errors"
	"strings"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
)

func GetTimetasticCookiesFromChrome() ([]*kooky.Cookie, error) {
	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("timetastic.co.uk"))
	var containsSecureCookie = false
	for _, cookie := range cookies {
		if strings.Contains(cookie.Name, "AspNetCore.Identity.Application") {
			containsSecureCookie = true
		}
	}
	if !containsSecureCookie {
		return nil, errors.New("No secure cookie found, please login to Timetastic in Chrome first")
	}

	return cookies, nil
}
