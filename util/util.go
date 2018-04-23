package util

import (
	"math/rand"

	"github.com/golang/glog"
)

func ErrorHandler(err error) {
	if err != nil {
		glog.Error(err)
	}
}

func RandIntOverRange(max, min int) int {
	if min > max {
		return 0
	}
	return rand.Intn(max-min) + min
}
