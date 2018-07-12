package config

//GenerateDSN GenerateDSN
func (c SqlConfig) GenerateDSN() string {
	//"TimeLine:TimeLineProject2018@tcp(sh2.casxt.com)/TimeLine?parseTime=true"
	return (c.User + ":" + c.Pass + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Base + "?parseTime=true")
}
