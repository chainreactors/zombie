package pkg

import (
	"strings"
)

func UseDefaultPassword(service string, top int) []string {
	if pwds, ok := Keywords[service+"_pwd"]; ok {
		if top == 0 || top > len(pwds) {
			return pwds
		} else {
			return pwds[:top]
		}
	} else {
		if top == 0 || top >= 10 {
			return Keywords["top10_pwd"]
		} else {
			return Keywords["top10_pwd"][:top]
		}
	}
}

func UseDefaultUser(service string, top int) []string {
	if users, ok := Keywords[service+"_user"]; ok {
		if top == 0 || top > len(users) {
			return users
		} else {
			return users[:top]
		}
	} else {
		if top == 0 || top >= 10 {
			return Keywords["top10_user"]
		} else {
			return Keywords["top10_user"][:top]
		}
	}
}

func SplitUserDomain(user string) (string, string) {
	var domain string
	if strings.Contains(user, "/") {
		parts := strings.SplitN(user, "/", 2)
		domain = parts[0]
		user = parts[1]
	}
	return user, domain
}
