package config

//SqlConfig contain sql config
type SqlConfig struct {
	Host string
	User string
	Pass string
	Port string
	Base string
}

//GenerateDSN GenerateDSN
func (c SqlConfig) GenerateDSN() string {
	//"TimeLine:TimeLineProject2018@tcp(sh2.casxt.com)/TimeLine?parseTime=true"
	return (c.User + ":" + c.Pass + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Base + "?parseTime=true")
}
