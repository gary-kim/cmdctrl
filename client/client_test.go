package client

import (
	"math/rand"
	"testing"
	"time"
)

func TestVerifySharedPass(t *testing.T) {

	randomSharedPass := ""
	rand.Seed(time.Now().UnixNano())
	possibleLetters := []byte("123567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 20; i++ {
		randomSharedPass += string(possibleLetters[rand.Intn(len(possibleLetters))])
	}
	got := RemoteRESTServer{
		opt: Options{
			SharedPass: randomSharedPass,
		},
	}.verifySharedPass(randomSharedPass)
	if !got {
		t.Errorf(`verifySharedPass("%s") returned %v, was expecting true`, randomSharedPass, got)
	}
}
