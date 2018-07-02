package database

import (
	"errors"
	"log"
	"strings"
)

//CreateLine Create Line and add a user to it's group
func CreateLine(LineName, UserID string) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

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
