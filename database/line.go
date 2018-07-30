package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

//GetLines get Lines of user
func GetLines(UserID string) (Lines []string, DBerr error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return nil, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	Lines = make([]string, 0, 1)
	rows, DBErr := course.Query("SELECT `L`.`Name` FROM `Line` `L` INNER JOIN `Group` `G` ON `G`.`LineID`=`L`.`ID` WHERE `G`.`UserID`=?", UserID)
	defer rows.Close()
	if DBErr != nil {
		log.Println(DBErr.Error())
		return nil, DBErr
	}
	for rows.Next() {
		var LineName string
		if err := rows.Scan(&LineName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		Lines = append(Lines, LineName)
	}
	return Lines, nil
}

//CreateLine Create Line and add a user to it's group
func CreateLine(LineName, UserID string) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	_, DBErr = course.Exec("INSERT INTO Line (`Name`) VALUES (?)",
		LineName)
	if DBErr != nil {
		switch {
		case strings.HasPrefix(DBErr.Error(), "Error 1062: Duplicate entry"):
			return errors.New("Line Already Exist")
		default:
			log.Println("CreateLine:", DBErr.Error())
			return errors.New("CreateLine Failde")
		}
	}

	_, DBErr = course.Exec("INSERT INTO `Group` (`Group`.`LineID`,`Group`.`UserID`) SELECT `Line`.`ID`, ? FROM `Line` WHERE `Line`.`Name`=?",
		UserID, LineName)
	if DBErr != nil {
		switch {
		//if create success, this will not Duplicate
		//case strings.HasPrefix(DBErr.Error(), "Error 1062: Duplicate entry"):
		//	return errors.New("Line Already Exist")
		default:
			log.Println("CreateLine:", DBErr.Error())
			return errors.New("CreateLine Failde")
		}
	}
	return nil
}

//GetLineInfo Get Line Info
func GetLineInfo(LineName string, course *sql.Tx) (LineID, Name string, CreateTime time.Time, DBErr error) {
	course, selfCourse, DBErr := Begin(course)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	row := course.QueryRow("SELECT `ID`, `Time` FROM `Line` WHERE `Name`=?", LineName)
	if DBErr = row.Scan(&LineID, &CreateTime); DBErr != nil {
		return "", "", time.Time{}, DBErr
	}
	return LineID, Name, CreateTime, nil
}

//GetLineDetail Get some statics info of line
//contain info of GetLineInfo
func GetLineDetail(LineName string, course *sql.Tx) (LineID, Name, LatestImg string, Users []string, SliceNum, ImgNum int, CreateTime, LatestTime time.Time, DBErr error) {
	course, selfCourse, DBErr := Begin(course)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	if LineID, LineName, CreateTime, DBErr = GetLineInfo(LineName, course); DBErr != nil {
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
	}
	var (
		imgNum     *int
		latestTime *time.Time
	)
	const sqlCmd string = `
		SELECT
			SUM( LENGTH( "Gallery" ) - LENGTH( REPLACE ( "Gallery", ',', '' ) ) + 1 ),
			COUNT( * ),
			MAX( "Time" ) 
		FROM
			"Slice" 
		WHERE
			"LineID" = ?`
	row := course.QueryRow(strings.Replace(sqlCmd, `"`, "`", -1), LineID)
	if DBErr = row.Scan(&imgNum, &SliceNum, &latestTime); DBErr != nil {
		log.Println("GetLineDetail-sql1", DBErr.Error())
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
	}
	if latestTime == nil {
		latestTime = &CreateTime
	}
	LatestTime = *latestTime
	if imgNum == nil {
		imgNum = new(int)
	}
	ImgNum = *imgNum

	const sqlCmd2 string = `
		SELECT 
			"User"."NickName" 
		FROM
			"User"
			INNER JOIN "Group" ON "User"."ID" = "Group"."UserID" 
		WHERE
			"Group"."LineID" = ?`
	var rows *sql.Rows

	if rows, DBErr = course.Query(strings.Replace(sqlCmd2, `"`, "`", -1), LineID); DBErr != nil {
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
	}

	for rows.Next() {
		var user string
		if DBErr = rows.Scan(&user); DBErr != nil {
			log.Println("GetLineDetail-sql2", DBErr.Error())
			return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
		}
		Users = append(Users, user)
	}

	const sqlCmd3 string = `
		SELECT
			"Gallery" 
		FROM
			"Slice" 
		WHERE
			"LineID" = ? AND "Gallery" != '' 
			ORDER BY "Time" DESC
			LIMIT 1`

	if row = course.QueryRow(strings.Replace(sqlCmd3, `"`, "`", -1), LineID); DBErr != nil {
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
	}

	switch DBErr = row.Scan(&LatestImg); {
	case DBErr == nil:
	case DBErr.Error() == "sql: no rows in result set":
		LatestImg = ""
		DBErr = nil
	default:
		log.Println("GetLineDetail-sql3", DBErr.Error())
		return "", "", "", nil, 0, 0, time.Time{}, time.Time{}, DBErr
	}
	LatestImg = strings.Split(LatestImg, ",")[0]

	return LineID, Name, LatestImg, Users, SliceNum, ImgNum, CreateTime, LatestTime, DBErr
}
