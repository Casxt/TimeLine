package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"
)

//CreateUser create a unverify user
func CreateUser(Phone, Mail, HashPass string) (NickName, Pass string, err error) {
	course, selfCourse, err := Begin(nil)
	if err != nil {
		return "", "", err
	}
	defer GraceCommit(course, selfCourse, err)
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
	_, err = course.Exec(sqlCmd, Phone, Mail, "Unverify User", Salt, HashSaltPass)
	if err != nil {
		log.Println("CreateUser:", err)
		return "", "", err
	}

	return "Unverify User", Pass, err
}
