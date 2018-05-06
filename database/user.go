package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/MapleFadeAway/timeline/mail"
)

//{'asdd','asdsadas'}
//CreateUser create a unverify user
func CreateUser(Phone string, Mail string, HashPass string) (err error) {
	course, selfCourse, err := Begin(nil)
	if err != nil {
		return err
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
		Pass := Phone + rnd.String()
		Hash256.Reset()
		Hash256.Write([]byte(Pass))
		HashPass = hex.EncodeToString(Hash256.Sum(nil))
		mail.SendMail(Mail, "TimeLine 注册验证", "<h1>TimeLine密码:"+Pass+"</h1>", nil)
	} else {
		//计算一个初始的随机pass
		mail.SendMail(Mail, "TimeLine 注册验证", "<h1>您已注册TimeLine</h1>", nil)
	}

	//HashSaltPass
	Hash256.Reset()
	Hash256.Write([]byte(Salt + HashPass))
	HashSaltPass := hex.EncodeToString(Hash256.Sum(nil))
	sqlCmd := "INSERT INTO User (`Phone`,`Mail`,`NickName`,`Salt`,`SaltPass`) VALUES (?,?,?,?,?)"
	_, err = course.Exec(sqlCmd, Phone, Mail, "Unverify User", Salt, HashSaltPass)
	if err != nil {
		log.Fatalln("CreateUser error:", err)
		return err
	}

	return err
}
