package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"
)

// EncodePassword .
// psw: 原密码
// 返回加密之后的密码
func EncodePassword(psw string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(psw), 5)
	if err != nil {
		log.Println("密码加密失败")
	}
	return string(hashedPassword)
}

// VerifyPsw 验证密码是否是正确
// rawPsw: 原密码
// enCodePassword: 加密过后的代码
func VerifyPsw(rawPsw, enCodePassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(enCodePassword), []byte(rawPsw)) == nil
}

func GetFileMd5(file *multipart.FileHeader) string {
	m := md5.New()
	f, err := file.Open()
	if err != nil {
		return ""
	}
	io.Copy(m, f)
	return hex.EncodeToString(m.Sum(nil))
}
