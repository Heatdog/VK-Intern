package cryptohash

import "golang.org/x/crypto/bcrypt"

func Hash(str string) ([]byte, error) {
	cost := 5
	return bcrypt.GenerateFromPassword([]byte(str), cost)
}

func VerifyHash(hashedStr []byte, str string) bool {
	if err := bcrypt.CompareHashAndPassword(hashedStr, []byte(str)); err != nil {
		return false
	}
	return true
}
