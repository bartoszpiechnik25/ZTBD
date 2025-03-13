package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("BEGIN_YEAR"))
	fmt.Println(os.Getenv("END_YEAR"))
	fmt.Println("Hi mom!")
}
