package utils

import (
	"gopkg.in/h2non/filetype.v1"
	"math/rand"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	VideoType   = "VIDEO"
	ImageType   = "IMAGE"
)

func GuessContentType(content []byte) string {
	if filetype.IsVideo(content) {
		return VideoType
	}
	if filetype.IsImage(content) {
		return ImageType
	}
	return ""
}

func MakeUniqueKey(length uint) string {
	startingTime := time.Now()
	rand.Seed(startingTime.Unix())
	buff := make([]rune, length)
	for idx := range buff {
		buff[idx] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(buff)
}
