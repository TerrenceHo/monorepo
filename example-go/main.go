package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("hello world")
	u := uuid.New()
	fmt.Println(u)
}
