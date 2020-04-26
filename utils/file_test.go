package utils

import (
	"fmt"
	"testing"
)

func TestFileIsExists(t *testing.T) {
	if FileIsExists("/Users/blacklad/tmp/") {
		fmt.Println("exits")
	}
}
