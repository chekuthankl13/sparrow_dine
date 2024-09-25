package middlewares

import "golang.org/x/crypto/bcrypt"

func HashPassword(psw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(psw), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckHashPsw(psw string, hashPSw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPSw), []byte(psw))
	return err == nil
}
