package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func GenerateHOTP(secret string, counter uint64) (passcode string, err error) {
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	buf := make([]byte, 8)
	mac := hmac.New(sha1.New, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)
	mac.Write(buf)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))
	mod_str := fmt.Sprintf("%d", mod)
	passcode = strings.Repeat("0", 6-len(mod_str)) + mod_str
	return passcode, nil
}

func main() {
	shortout := false
	t := time.Now()
	secret := os.Getenv("TOTP_SECRET")
	if len(os.Args) > 1 {
		if os.Args[1] == "-s" {
			secret = os.Args[2]
			shortout = true
		} else {
			secret = os.Args[1]
		}
	}
	if len(secret) < 1 {
		fmt.Printf("Set TOTP_SECRET env var or\n# %s <mysecret>\n", os.Args[0])
		os.Exit(1)
	}
	counter := uint64(math.Floor(float64(t.Unix()) / float64(30)))
	totp, err := GenerateHOTP(secret, counter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if shortout {
		fmt.Printf("%+v", totp)
	} else {
		fmt.Printf("totp: %+v (%ds left)\n", totp, 30-(t.Unix()%30))
	}
	os.Exit(0)

}
