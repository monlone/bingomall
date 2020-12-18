package helper

import (
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"math/rand"
	"time"
)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GenerateId() uuid.UUID {
	id := uuid.NewV4()
	return id
}

func GenerateId36() string {
	return uuid.NewV4().String()
}

func GenerateId32() string {
	return ksuid.New().String()
}

type BaseUserInfo struct {
	OpenId    string `json:"openId"`
	Nickname  string `json:"nickName"`
	Gender    string `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"-"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Token     string `json:"token"`
	UserID    uint64 `json:"userId"`
}
