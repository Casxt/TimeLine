package database

import (
	"errors"
	"log"
)

//CreateSlice Create Slice
func CreateSlice(LineName, UserID, Content, Gallery, Type, Visibility, Location, Time string) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	_, DBErr = course.Exec("INSERT INTO `Slice` (`UserID`,`LineID`,`Content`,`Gallery`,`Location`,`Type`,`Visibility`,`Time`) SELECT  ?, `Line`.`ID`, ?, ?, ?, ?, ?, ? FROM `Line` WHERE `Line`.`Name`=?",
		UserID, Content, Gallery, Location, Type, Visibility, Time, LineName)
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
