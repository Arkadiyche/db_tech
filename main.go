package main

import (
	"fmt"
	"githun.com/Arkadiyche/bd_techpark/check"
)

func main() {
	if check.Check() {
		fmt.Println("a")
	}
}
