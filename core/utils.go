package core

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

var nameStripRE = regexp.MustCompile("^((20)|(25)|(2b)|(2f)|(3d)|(3a)|(40))+")

func GetRawCookie(cookies []*http.Cookie) string {
	var rawCookies []string
	for _, c := range cookies {
		e := fmt.Sprintf("%s=%s", c.Name, c.Value)
		rawCookies = append(rawCookies, e)
	}
	return strings.Join(rawCookies, "; ")
}

func GetDomain(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	domain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return ""
	}
	return domain
}

func GetHostname(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func FixUrl(url, site string) string {
	var newUrl string

	if strings.HasPrefix(url, "http") {
		// http://google.com || https://google.com
		newUrl = url
	} else if strings.HasPrefix(url, "//") {
		// //google.com/example.php
		newUrl = "https:" + url
	} else if !strings.HasPrefix(url, "http") && len(url) > 0 {
		if url[:1] == "/" {
			// Ex: /?thread=10
			newUrl = site + url
		} else {
			// Ex: ?thread=10
			newUrl = site + "/" + url
		}
	}
	return newUrl
}

func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func LoadCookies(rawCookie string) []*http.Cookie {
	httpCookies := []*http.Cookie{}
	cookies := strings.Split(rawCookie, ";")
	for _, cookie := range cookies {
		cookieArgs := strings.SplitN(cookie, "=", 2)
		if len(cookieArgs) > 2 {
			continue
		}

		ck := &http.Cookie{Name: strings.TrimSpace(cookieArgs[0]), Value: strings.TrimSpace(cookieArgs[1])}
		httpCookies = append(httpCookies, ck)
	}
	return httpCookies
}

func GetExtType(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	return path.Ext(u.Path)
}

func CleanSubdomain(s string) string {
	s = cleanName(s)
	s = strings.TrimPrefix(s, "*.")
	s = strings.TrimPrefix(s,"u002f")
	return s
}

// Clean up the names scraped from the web.
// Get from Amass
func cleanName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))

	for {
		if i := nameStripRE.FindStringIndex(name); i != nil {
			name = name[i[1]:]
		} else {
			break
		}
	}

	name = strings.Trim(name, "-")
	// Remove dots at the beginning of names
	if len(name) > 1 && name[0] == '.' {
		name = name[1:]
	}
	return name
}
