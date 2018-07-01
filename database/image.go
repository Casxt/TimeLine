package database

import (
	"errors"
	"log"
	"strings"
)

//CreateImage will mark a group of img one by one belong to a user
func CreateImage(UserID string, ImageHashs []string) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "INSERT INTO Image (`Hash`,`UserID`) VALUES (?,?)"
	stmt, err := course.Prepare(sqlCmd)
	if err != nil {
		log.Println("CreateImage:", DBErr.Error())
		return err
	}
	for hash := range ImageHashs {
		_, DBErr = stmt.Exec(hash, UserID)
		if DBErr != nil {
			switch {
			case strings.HasPrefix(DBErr.Error(), "Error 1062: Duplicate entry"):
				DBErr = nil
				continue
			default:
				log.Println("CreateImage:", DBErr.Error())
				return errors.New("Image CreateImage Failde")
			}
		}
	}

	return nil
}
