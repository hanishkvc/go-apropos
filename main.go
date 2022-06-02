package main

import (
	"flag"
	"fmt"
)

const PRG_NAME = "GoApropos"
const PRG_VERSION = "v0-20220602IST0954"

var gFind string

func handle_args() {
	flag.StringVar(&gFind, "find", "", "Specify the word to find")
	flag.Parse()
	fmt.Printf("gFind: %v\n", gFind)
}

func main() {
	fmt.Println(PRG_NAME, PRG_VERSION)
	handle_args()
}
