package authentication

import (
	"context"
	"errors"
	"strings"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

func GetTimetasticCookiesFromChrome() ([]*kooky.Cookie, error) {
	cookies, err := kooky.ReadCookies(context.Background(), kooky.Valid, kooky.DomainHasSuffix("timetastic.co.uk"))
	if err != nil {
		return nil, err
	}
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
