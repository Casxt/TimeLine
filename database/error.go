package database

//DBerror is Database error
type DBerror struct {
	ErrorMsg string
}

//Error 实现 `error` 接口
func (e DBerror) Error() string {
	return e.ErrorMsg
}
