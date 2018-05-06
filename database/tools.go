package database

import (
	"database/sql"

	//init database
	_ "github.com/go-sql-driver/mysql"
)

//DataBase is global database var
var DataBase *sql.DB

//Open using a sql.Open to open database
func Open() error {
	//此处小心，:= 赋值会建立一个叫DataBase的局部变量,并且在函数执行期间替换掉全局的DataBase
	DataBase, _ = sql.Open("mysql", "TimeLine:TimeLineProject2018@tcp(sh2.casxt.com)/TimeLine")
	DataBase.SetMaxOpenConns(100)
	DataBase.SetMaxIdleConns(50)
	err := DataBase.Ping()
	return err
}

//Close database
func Close() error {
	err := DataBase.Close()
	return err
}

//Begin Automatic call DB.Begin to begin a course
func Begin(course *sql.Tx) (*sql.Tx, bool, error) {
	if course == nil {
		course, err := DataBase.Begin()
		return course, true, err
	}
	return course, false, nil
}

//Commit Automatic call DB.Commit to Commit a course
func Commit(course *sql.Tx, selfCourse bool) error {
	if selfCourse == true {
		err := course.Commit()
		return err
	}
	return nil
}

//Rollback Automatic call DB.Rollback to Commit a Rollback
func Rollback(course *sql.Tx, selfCourse bool) error {
	if selfCourse == true {
		err := course.Rollback()
		return err
	}
	return nil
}

//GraceCommit will automitic commit or rollback
func GraceCommit(course *sql.Tx, selfCourse bool, err error) error {
	if err != nil {
		return Rollback(course, selfCourse)
	}
	return Commit(course, selfCourse)
}
