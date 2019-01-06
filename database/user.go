package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/Casxt/TimeLine/tools"
	"log"
	"math/big"
	"strings"
	"time"
)

//GetUserByMail Found User By Mail
func GetUserByMail(Mail string) (ID, Phone, NickName, Gender, Salt, SaltPass, ProfilePic string, SignInTime time.Time, ErrorMsg error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "SELECT `ID`,`Phone`,`NickName`,`Gender`,`Salt`,`SaltPass`,`ProfilePic`,`Time` FROM `User` WHERE `Mail`=?"
	Row := course.QueryRow(sqlCmd, Mail)

	if DBErr = Row.Scan(&ID, &Phone, &NickName, &Gender, &Salt, &SaltPass, &ProfilePic, &SignInTime); DBErr != nil {
		switch DBErr.Error() {
		case "sql: no rows in result set":
			return "", "", "", "", "", "", "", time.Time{}, errors.New("User Not Exist")
		default:
			log.Println(DBErr.Error())
			return "", "", "", "", "", "", "", time.Time{}, errors.New("Unknow Error")
		}
	}
	return ID, Mail, NickName, Gender, Salt, SaltPass, ProfilePic, SignInTime, nil
}

//GetUserByPhone Found User By Phone
//Err :
//DataBase Connection Error
//User Not Exist
//DBErr
func GetUserByPhone(Phone string, course *sql.Tx) (ID, Mail, NickName, Gender, Salt, SaltPass, ProfilePic string, SignInTime time.Time, DBErr error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "SELECT `ID`,\n       `Mail`,\n       `NickName`,\n       `Gender`,\n       `Salt`,\n       `SaltPass`,\n       `ProfilePic`,\n       `Time`\nFROM `User`\nWHERE `Phone` = ?"
	Row := course.QueryRow(sqlCmd, Phone)

	if DBErr = Row.Scan(&ID, &Mail, &NickName, &Gender, &Salt, &SaltPass, &ProfilePic, &SignInTime); DBErr != nil {
		switch DBErr.Error() {
		case "sql: no rows in result set":
			return "", "", "", "", "", "", "", time.Time{}, errors.New("User Not Exist")
		default:
			log.Println("GetUserByPhone", DBErr.Error())
			return "", "", "", "", "", "", "", time.Time{}, DBErr
		}
	}

	return ID, Mail, NickName, Gender, Salt, SaltPass, ProfilePic, SignInTime, nil
}

func GetUserByOpenID(OpenID string) (ID, UnionID, Phone, Mail, NickName, Gender, Salt, SaltPass, ProfilePic string, SignInTime time.Time, DBErr error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "SELECT `ID`,\n       `UnionID`,\n       `Mail`,\n       `NickName`,\n       `Gender`,\n       `Salt`,\n       `SaltPass`,\n       `ProfilePic`,\n       `Time`\nFROM `User`\nWHERE `OpenID` = ?"
	Row := course.QueryRow(sqlCmd, OpenID)

	if DBErr = Row.Scan(&ID, &UnionID, &Mail, &NickName, &Gender, &Salt, &SaltPass, &ProfilePic, &SignInTime); DBErr != nil {
		switch DBErr.Error() {
		case "sql: no rows in result set":
			return "", "", "", "", "", "", "", "", "", time.Time{}, errors.New("User Not Exist")
		default:
			log.Println("GetUserByPhone", DBErr.Error())
			return "", "", "", "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
		}
	}

	return ID, UnionID, Phone, Mail, NickName, Gender, Salt, SaltPass, ProfilePic, SignInTime, nil
}

//CreateUser create a unverify user
func CreateUser(Phone, Mail, HashPass string, OpenID, UnionID *string) (NickName, Pass string, ErrorMsg error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	//随机
	rnd, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	//计算一个初始的随机salt
	Hash256 := sha256.New()
	Hash256.Write([]byte(rnd.String()))
	Salt := hex.EncodeToString(Hash256.Sum(nil))
	if HashPass == "" {
		//计算一个初始的随机pass并Hash
		Pass = rnd.String()
		Hash256.Reset()
		Hash256.Write([]byte(Pass))
		HashPass = hex.EncodeToString(Hash256.Sum(nil))
	}
	//HashSaltPass
	Hash256.Reset()
	Hash256.Write([]byte(Salt + HashPass))
	HashSaltPass := hex.EncodeToString(Hash256.Sum(nil))
	sqlCmd := "INSERT INTO User (`OpenID`, `UnionID`, `Phone`, `Mail`, `NickName`, `Salt`, `SaltPass`)\nVALUES (?, ?, ?, ?, ?, ?, ?)"
	_, DBErr = course.Exec(sqlCmd, tools.AutoNullString(OpenID), tools.AutoNullString(UnionID), Phone, Mail, "UnverifyUser", Salt, HashSaltPass)
	if DBErr != nil {
		switch {
		case strings.HasPrefix(DBErr.Error(), "Error 1062: Duplicate entry"):
			return "", "", errors.New("User Already Exist")
		default:
			log.Println("CreateUser:", DBErr.Error())
			return "", "", errors.New("User Create Failde")
		}
	}

	return "Unverify User", Pass, nil
}

func UpdateWeiXin(OpenID, Phone, Mail string, UnionID *string) (DBErr error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	sqlCmd := "Update `User` SET `OpenID`=?, `UnionID`=? WHERE `Phone`=? AND `Mail`=?"
	res, DBErr := course.Exec(sqlCmd, OpenID, tools.AutoNullString(UnionID), Phone, Mail)
	if DBErr != nil {
		log.Println("UpdateWeiXin:", DBErr.Error())
		return DBErr
	}

	num, DBErr := res.RowsAffected()
	if DBErr != nil {
		log.Println("UpdateWeiXin:", DBErr.Error())
		return DBErr
	}

	if num > 1 {
		DBErr = errors.New("Update too many user ")
		log.Println("UpdateWeiXin:", DBErr.Error())
		return DBErr
	}
	if num < 1 {
		DBErr = errors.New("User Not Exist ")
		log.Println("UpdateWeiXin:", DBErr.Error())
		return DBErr
	}
	return nil
}

//UpdateProfilePic Update ProfilePicture
//Mention this func will set img Visibility to Public
func UpdateProfilePic(UserID, ImgHash string) (ErrorMsg error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	//set img Visibility to Public
	UpdateImgVisibility(UserID, ImgHash, "Public", course)

	sqlCmd := "Update `User` SET `ProfilePic`=? WHERE `ID`=?"
	_, DBErr = course.Exec(sqlCmd, ImgHash, UserID)
	if DBErr != nil {
		log.Println("UpdateProfilePic:", DBErr.Error())
		return DBErr
	}

	return nil
}

//UpdateNickName Update User NickName
func UpdateNickName(UserID, NewName string, course *sql.Tx) (DBErr error) {
	course, selfCourse, DBErr := Begin(course)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return errors.New("DataBase Connection Error")
	}
	defer func() { GraceCommit(course, selfCourse, DBErr) }()

	sqlCmd := "Update `User` SET `NickName`=? WHERE `ID`=?"
	_, DBErr = course.Exec(sqlCmd, NewName, UserID)
	if DBErr != nil {
		log.Println("UpdateNickName:", DBErr.Error())
		return DBErr
	}
	return nil
}
