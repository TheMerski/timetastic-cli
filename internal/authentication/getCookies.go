package authentication

import (
	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
)

func GetTimetasticCookiesFromChrome() []*kooky.Cookie {
	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("timetastic.co.uk"))
	return cookies
}
