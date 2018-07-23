package tools

import "regexp"

//CheckImgHash Check ImgHash Format
func CheckImgHash(Hash string) bool {
	return regexp.MustCompile("^[a-z0-9]{64}$").MatchString(Hash)
}

//ChecNickName Check NickName Format
func ChecNickName(Hash string) bool {
	return regexp.MustCompile("^[\\S]{2,32}$").MatchString(Hash)
}
