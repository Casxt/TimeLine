package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"time"
)

//GetUserByMail Found User By Mail
func GetUserByMail(Mail string) (Phone, NickName, Gender, Salt, SaltPass, ProfilePic string, SignInTime time.Time, ErrorMsg error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "SELECT `Phone`,`NickName`,`Gender`,`Salt`,`SaltPass`,`ProfilePic`,`Time` FROM `User` WHERE `Mail`=?"
	Row := course.QueryRow(sqlCmd, Mail)

	if DBErr = Row.Scan(&Phone, &NickName, &Gender, &Salt, &SaltPass, &ProfilePic, &SignInTime); DBErr != nil {
		switch DBErr.Error() {
		case "no rows in result set":
			return "", "", "", "", "", "", time.Time{}, errors.New("User Not Exist")
		default:
			log.Println(DBErr.Error())
			return "", "", "", "", "", "", time.Time{}, errors.New("Unknow Error")
		}
	}

	return Mail, NickName, Gender, Salt, SaltPass, ProfilePic, SignInTime, nil
}

//GetUserByPhone Found User By Phone
func GetUserByPhone(Phone string) (Mail, NickName, Gender, Salt, SaltPass, ProfilePic string, SignInTime time.Time, ErrorMsg error) {
	course, selfCourse, DBErr := Begin(nil)
	if DBErr != nil {
		log.Println(DBErr.Error())
		return "", "", "", "", "", "", time.Time{}, errors.New("DataBase Connection Error")
	}
	defer GraceCommit(course, selfCourse, DBErr)

	sqlCmd := "SELECT `Mail`,`NickName`,`Gender`,`Salt`,`SaltPass`,`ProfilePic`,`Time` FROM `User` WHERE `Phone`=?"
	Row := course.QueryRow(sqlCmd, Phone)

	if DBErr = Row.Scan(&Mail, &NickName, &Gender, &Salt, &SaltPass, &ProfilePic, &SignInTime); DBErr != nil {
		switch DBErr.Error() {
		case "no rows in result set":
			return "", "", "", "", "", "", time.Time{}, errors.New("User Not Exist")
		default:
			log.Println(DBErr.Error())
			return "", "", "", "", "", "", time.Time{}, errors.New("Unknow Error")
		}
	}

	return Mail, NickName, Gender, Salt, SaltPass, ProfilePic, SignInTime, nil
}

//CreateUser create a unverify user
func CreateUser(Phone, Mail, HashPass string) (NickName, Pass string, ErrorMsg error) {
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
	sqlCmd := "INSERT INTO User (`Phone`,`Mail`,`NickName`,`Salt`,`SaltPass`) VALUES (?,?,?,?,?)"
	_, DBErr = course.Exec(sqlCmd, Phone, Mail, "Unverify User", Salt, HashSaltPass)
	if DBErr != nil {
		log.Println("CreateUser:", DBErr.Error())
		return "", "", errors.New("User SignIn Failde")
	}

	return "Unverify User", Pass, nil
}
