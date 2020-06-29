package main

import (
	"coursera-web-services/week1"
	"os"
)

func callWeek1() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := week1.DirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	callWeek1()
}
