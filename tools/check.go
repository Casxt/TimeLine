package tools

import "regexp"

func CheckImgHash(Hash string) bool {
	return regexp.MustCompile("^[a-z0-9]{64}$").MatchString(Hash)
}

func ChecNickName(Hash string) bool {
	return regexp.MustCompile("^[\\S]{2,32}$").MatchString(Hash)
}
