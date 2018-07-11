package database

import (
	"errors"
	"log"
	"strings"
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
