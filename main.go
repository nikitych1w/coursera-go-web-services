package main

import (
	"fmt"
	"os"
	"time"
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

func callWeek2() {
	//inputData := []int{0, 1}
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}

	start = time.Now()
	ExecutePipeline(hashSignJobs...)
	fmt.Println("MAIN END", time.Since(start))
}

func main() {
	callWeek1()
}
