package tools

import "regexp"

func CheckImgHash(Hash string) bool {
	return regexp.MustCompile("^[a-z0-9]{64}$").MatchString(Hash)
}
