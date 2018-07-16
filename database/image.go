package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

type ImgInfo struct {
	Hash   string
	Size   int
	Height int
	Width  int
}

func GetImgInfo(UserID, Hash string) (Size, Height, Width int, Time time.Time, err error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return 0, 0, 0, time.Time{}, errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	sqlCmd := "SELECT `Size`,`Height`,`Width`,`Time` FROM `Image` WHERE `Hash`=? AND `UserID`=?"
	Row := course.QueryRow(sqlCmd, Hash, UserID)
	if DBErr = Row.Scan(&Size, &Height, &Width, &Time); DBErr != nil {
		switch DBErr.Error() {
		case "sql: no rows in result set":
			return 0, 0, 0, time.Time{}, errors.New("Img-User Not Exist")
		default:
			log.Println(DBErr.Error())
			return 0, 0, 0, time.Time{}, errors.New("Unknow Error")
		}
	}
	return Size, Height, Width, Time, nil
}

//CreateImage will mark a group of img one by one belong to a user
func CreateImage(UserID string, Images []ImgInfo) error {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	stmt, DBErr := course.Prepare("INSERT INTO Image (`Hash`,`UserID`,`Visibility`,`Size`,`Height`,`Width`) VALUES (?,?,?,?,?,?)")
	if DBErr != nil {
		log.Println("CreateImage:", DBErr.Error())
		return DBErr
	}

	for _, imgInfo := range Images {
		_, DBErr = stmt.Exec(imgInfo.Hash, UserID, "Protect", imgInfo.Size, imgInfo.Height, imgInfo.Width)
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

func UpdateImgVisibility(UserID, ImgHash, Visibility string, course *sql.Tx) (ErrorMsg error) {
	course, selfCourse, DBErr := Begin(course)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "Update `Image` SET `Visibility`=? WHERE `Hash`=? AND `UserID`=?"
	_, DBErr = course.Exec(sqlCmd, Visibility, ImgHash, UserID)
	if DBErr != nil {
		log.Println("UpdateProfilePic:", DBErr.Error())
		return DBErr
	}
	return nil
}
