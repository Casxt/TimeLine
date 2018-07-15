package database

import (
	"errors"
	"log"
	"strings"
	"time"
)

func CheckImgInfo(Hash, UserID string) (Time time.Time, err error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return time.Time{}, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	sqlCmd := "SELECT `Time` FROM `Image` WHERE `Hash`=? AND `UserID`=?"
	Row := course.QueryRow(sqlCmd, Hash, UserID)
	if DBErr = Row.Scan(&Time); DBErr != nil {
		switch DBErr.Error() {
		case "sql: no rows in result set":
			return time.Time{}, errors.New("Img-User Not Exist")
		default:
			log.Println(DBErr.Error())
			return time.Time{}, errors.New("Unknow Error")
		}
	}
	return Time, nil
}

//CreateImage will mark a group of img one by one belong to a user
func CreateImage(UserID string, ImageHash []string, ImageSize []int) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	stmt, err := course.Prepare("INSERT INTO Image (`Hash`,`UserID`,`Visibility`,`Size`) VALUES (?,?,?,?)")
	if err != nil {
		log.Println("CreateImage:", DBErr.Error())
		return err
	}

	for index, hash := range ImageHash {
		_, DBErr = stmt.Exec(hash, UserID, "Protect", ImageSize[index])
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
