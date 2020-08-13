package main

import (
	"fmt"
	"testing"
)

func TestOne(t *testing.T) {
	valid := isHashValid("000000xxxx000000", 6)
	fmt.Println(valid)

	i := 15
	fmt.Println(fmt.Sprintf("%x", i))
}
