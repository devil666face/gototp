package main

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"fmt"
	"time"
)

func GeneratePassCode(utf8string string) string {
	// secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	secret := utf8string
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

func main() {
	fmt.Println(GeneratePassCode("5DAJLNUGVMEJWAPCXHP4UCOTXULH33KD"))
}
