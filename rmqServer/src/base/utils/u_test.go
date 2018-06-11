package utils

import (
	"fmt"
	"testing"
)

func TestRandStr(t *testing.T) {
	for index := 0; index < 5; index++ {
		fmt.Println(RandomString(6))
	}

	t.Error()
}
