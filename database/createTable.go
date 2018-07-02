package database

import (
	"database/sql"
	"log"
	"strings"
)

//CreateUserTable Create User Table
func CreateUserTable(course *sql.Tx) (err error) {
	course, selfCourse, err := Begin(course)
	if err != nil {
		return err
	}
	defer GraceCommit(course, selfCourse, err)
	sqlCmd := `CREATE TABLE 'User' (
		'ID'  int NOT NULL AUTO_INCREMENT ,
		'Phone'  varchar(15) CHARACTER SET utf8 NOT NULL,
		'Mail'  varchar(128) CHARACTER SET utf8 NOT NULL,
		'NickName'  varchar(128) CHARACTER SET utf8 NOT NULL,
		'Gender'  enum("Male","Female","Secret") DEFAULT "Secret" NOT NULL,
		'Salt'  char(64) CHARACTER SET utf8 NOT NULL,
		'SaltPass'  char(64) CHARACTER SET utf8 NOT NULL,
		'ProfilePic'  char(64) CHARACTER SET utf8 DEFAULT "" NOT NULL,
		'Time'  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE INDEX ('Phone'),
		UNIQUE INDEX ('Mail'),
		PRIMARY KEY ('ID')
		)`
	//'Salt'  char(64) CHARACTER SET utf8 NOT NULL,
	sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
	result, err := course.Exec(sqlCmd)
	log.Println("Create Table User, Result:", result, "Error:", err)
	return err
}

//CreateLineTable Create Line Table
func CreateLineTable(course *sql.Tx) (err error) {
	course, selfCourse, err := Begin(course)
	if err != nil {
		return err
	}
	defer GraceCommit(course, selfCourse, err)
	sqlCmd := `CREATE TABLE 'Line' (
		'ID'  int NOT NULL AUTO_INCREMENT ,
		'Name'  varchar(128) CHARACTER SET utf8 NOT NULL,
		'Time'  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE INDEX ('Name'),
		PRIMARY KEY ('ID')
		)`
	sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
	result, err := course.Exec(sqlCmd)
	log.Println("Create Table Line, Result:", result, "Error:", err)

	return err
}

//CreateGroupTable Create Group Table, Group recoder the relation between user and line
func CreateGroupTable(course *sql.Tx) (err error) {
	course, selfCourse, err := Begin(course)
	if err != nil {
		return err
	}
	defer GraceCommit(course, selfCourse, err)
	sqlCmd := `CREATE TABLE 'Group' (
		'ID'  int NOT NULL AUTO_INCREMENT ,
		'LineID' int NOT NULL,
		'UserID'  int NOT NULL,
		'Time'  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY ('LineID') REFERENCES 'Line' ('ID'),
		FOREIGN KEY ('UserID') REFERENCES 'User' ('ID'),
		INDEX ('LineID'),
		INDEX ('UserID'),
		PRIMARY KEY ('ID')
		)`
	sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
	result, err := course.Exec(sqlCmd)
	log.Println("Create Table Group, Result:", result, "Error:", err)

	return err
}

//CreateImageTable Create image Table, image recoder the image
func CreateImageTable(course *sql.Tx) (err error) {
	course, selfCourse, err := Begin(course)
	if err != nil {
		return err
	}
	defer GraceCommit(course, selfCourse, err)
	sqlCmd := `CREATE TABLE 'Image' (
		'ID'  int NOT NULL AUTO_INCREMENT ,
		'Hash' char(64) CHARACTER SET utf8 NOT NULL,
		'UserID'  int NOT NULL,
		'Time'  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY ('UserID') REFERENCES 'User' ('ID'),
		INDEX ('UserID'),
		UNIQUE INDEX ('Hash','UserID'),
		PRIMARY KEY ('ID')
		)`
	sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
	result, err := course.Exec(sqlCmd)
	log.Println("Create Table Group, Result:", result, "Error:", err)

	return err
}

//CreateSliceTable Create Slice Table, Slice is the unite in the line, it has two basic type, memory and Anniversary
//the visibillty Protect means only people join this group can view, Private means only people create this slice can view
func CreateSliceTable(course *sql.Tx) (err error) {
	course, selfCourse, err := Begin(course)
	if err != nil {
		return err
	}
	//'Title'  varchar(128) CHARACTER SET utf8 NOT NULL,
	//		INDEX ('Title'),
	defer GraceCommit(course, selfCourse, err)
	sqlCmd := `CREATE TABLE 'Slice' (
		'ID'  int NOT NULL AUTO_INCREMENT ,
		'UserID'  int NOT NULL,
		'LineID'  int NOT NULL,
		'Content'  varchar(2048) CHARACTER SET utf8 NOT NULL,
		'Gallery'  varchar(1024) CHARACTER SET utf8 NOT NULL,
		'Location'  varchar(128) CHARACTER SET utf8 ,
		'Type'  enum("Memory","Anniversary") DEFAULT "Memory",
		'Visibility'  enum("Private","Public","Protect") DEFAULT "Protect",
		'Time'  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY ('LineID') REFERENCES 'Line' ('ID'),
		FOREIGN KEY ('UserID') REFERENCES 'User' ('ID'),
		INDEX ('LineID'),
		INDEX ('UserID'),
		INDEX ('Content'),
		PRIMARY KEY ('ID')
		)`
	sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
	result, err := course.Exec(sqlCmd)
	log.Println("Create Table Slice, Result:", result, "Error:", err)

	return err
}
