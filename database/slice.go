package database

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

//SliceInfo struct use in GetSlice func
type SliceInfo struct {
	UserName   string
	Gallery    []string
	Content    string
	Type       string
	Visibility string
	Location   string
	Time       time.Time
}

//GetSlices Get Slice List
func GetSlices(LineName string, UserID string, PageNum int) (Res []SliceInfo, DBerr error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return nil, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()
	var inGroup bool = false
	var sqlCmd string
	//If User Already Sign In
	if UserID != "" {
		sqlCmd = `SELECT EXISTS(SELECT * FROM 'Group' 'G' INNER JOIN 'Line' 'L' ON 'G'.'LineID'='L'.'ID' 
					WHERE 'L'.'Name'=? AND 'G'.'UserID'=?);`
		row := course.QueryRow(strings.Replace(sqlCmd, "'", "`", -1), LineName, UserID)
		if err := row.Scan(&inGroup); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	var rows *sql.Rows
	if inGroup == false {
		//User not in Group
		//Can Only See Public
		sqlCmd = `SELECT 'U'.'NickName', 'S'.'Gallery', 'S'.'Content', 'S'.'Type', 'S'.'Visibility', 'S'.'Location', 'S'.'Time'
					FROM 'User' 'U', 'Slice' 'S', 'Line' 'L'
					WHERE 'S'.'Visibility'="Public" AND 'L'.'Name'=? AND 'S'.'LineID'='L'.'ID' AND 'U'.'ID'='S'.'UserID'
					ORDER BY 'S'.'Time' DESC LIMIT 0, 20`
		sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
		sqlCmd = strings.Replace(sqlCmd, `"`, `'`, -1)
		rows, DBErr = course.Query(sqlCmd, LineName)
	} else {
		//User in Group
		//Can See Public Protect and self Private
		sqlCmd = `SELECT 'U'.'NickName', 'S'.'Gallery', 'S'.'Content', 'S'.'Type', 'S'.'Visibility', 'S'.'Location', 'S'.'Time'
					FROM ('User' 'U' INNER JOIN 'Slice' 'S' ON 'U'.'ID'='S'.'UserID') INNER JOIN 'Line' 'L' ON 'S'.'LineID'='L'.'ID'
					WHERE 'L'.'Name'=? AND ('S'.'Visibility' IN ("Public", "Protect") OR 'S'.'UserID'=?) 
					ORDER BY 'S'.'Time' DESC LIMIT 0, 20`
		sqlCmd = strings.Replace(sqlCmd, "'", "`", -1)
		sqlCmd = strings.Replace(sqlCmd, `"`, `'`, -1)
		rows, DBErr = course.Query(sqlCmd, LineName, UserID)
	}
	defer rows.Close()
	if DBErr != nil {
		log.Println(DBErr.Error())
		return nil, DBErr
	}
	Res = make([]SliceInfo, 0, 20)
	for rows.Next() {
		//var s SliceInfo
		s := new(SliceInfo)
		var gallery string
		if err := rows.Scan(&(s.UserName), &gallery, &(s.Content), &(s.Type), &(s.Visibility), &(s.Location), &(s.Time)); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		if gallery != "" {
			s.Gallery = strings.Split(gallery, ",")
		}
		Res = append(Res, *s)
	}
	return Res, nil
}

//CreateSlice Create Slice
func CreateSlice(LineName, UserID string, Gallery []string, Content, Type, Visibility, Location, Time string) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	//Gallery into hash1,hahs2,...hashn, format
	var galleryString string
	imgNum := len(Gallery)
	if imgNum > 0 {
		//64*imgNum+imgNum
		buff := bytes.NewBuffer(make([]byte, 65*imgNum))
		buff.Reset()
		for _, Hash := range Gallery {
			buff.WriteString(Hash)
			buff.Write([]byte(","))
		}
		galleryString = string(buff.Bytes()[0 : buff.Len()-2])
	} else {
		galleryString = ""
	}

	_, DBErr = course.Exec("INSERT INTO `Slice` (`UserID`,`LineID`,`Content`,`Gallery`,`Location`,`Type`,`Visibility`,`Time`) SELECT  ?, `Line`.`ID`, ?, ?, ?, ?, ?, ? FROM `Line` WHERE `Line`.`Name`=?",
		UserID, Content, galleryString, Location, Type, Visibility, Time, LineName)
	if DBErr != nil {
		switch {
		//if create success, this will not Duplicate
		//case strings.HasPrefix(DBErr.Error(), "Error 1062: Duplicate entry"):
		//	return errors.New("Line Already Exist")
		default:
			log.Println("CreateSlice:", DBErr.Error())
			return errors.New("CreateSlice Failde")
		}
	}

	return nil
}
