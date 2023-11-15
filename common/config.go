package common

import (
	"os"
	"strconv"
	"time"
)

func Timeout() time.Duration {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		panic(err)
	}

	return time.Duration(timeout) * time.Second
}
