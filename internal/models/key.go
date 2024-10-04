package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var (
	algorithms = map[string]otp.Algorithm{
		"SHA1":   otp.AlgorithmSHA1,
		"SHA256": otp.AlgorithmSHA256,
		"SHA512": otp.AlgorithmSHA512,
		"MD5":    otp.AlgorithmMD5,
	}
	digits = map[string]otp.Digits{
		"6": otp.DigitsSix,
		"8": otp.DigitsEight,
	}
)

type Keystore struct {
	Keys []Key
}

type Input struct {
	Name      string
	Period    string
	Digit     string
	Algorithm string
	Secret    string
}

func (k *Keystore) Add(input Input) {
	k.Keys = append(k.Keys, NewKey(
		input.Name,
		input.Period,
		input.Digit,
		input.Algorithm,
		input.Secret,
	))
}

func (k *Keystore) Delete(id int) error {
	if id < 0 || id >= len(k.Keys) {
		return fmt.Errorf("index out of range %d", len(k.Keys))
	}
	k.Keys = append(k.Keys[:id], k.Keys[id+1:]...)
	return nil
}

func (k *Keystore) Suggestions() []string {
	var (
		suggestions []string
		separators  = []string{
			"@",
			".",
			",",
			"/",
			`\`,
			"(",
			")",
			"_",
		}
	)
	for _, k := range k.Keys {
		for _, s := range separators {
			for i := 0; i <= len([]rune(k.Name)); i++ {
				suggestions = append(suggestions, strings.SplitAfterN(k.Name, s, i)...)
			}
		}
	}
	return suggestions
}

type Key struct {
	Name      string
	Period    int
	Digits    otp.Digits
	Algorithm otp.Algorithm
	Secret    string
}

func NewKey(
	name,
	period,
	digit,
	algorithm,
	secret string,
) Key {
	key := Key{
		Name:      name,
		Digits:    digits[digit],
		Algorithm: algorithms[algorithm],
		Secret:    secret,
	}
	key.Period, _ = strconv.Atoi(period)
	return key
}

func (k Key) String() string {
	return fmt.Sprintf("%s", k.Name)
}

func (k Key) GenCode() (string, error) {
	passcode, err := totp.GenerateCodeCustom(k.Secret, time.Now(), totp.ValidateOpts{
		Period:    uint(k.Period),
		Skew:      1,
		Digits:    k.Digits,
		Algorithm: k.Algorithm,
	})
	if err != nil {
		return "", err
	}
	return passcode, nil
}
