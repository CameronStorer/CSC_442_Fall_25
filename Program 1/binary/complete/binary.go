package main

import (
	"fmt"
	//"io"
	"os"
	"strconv"
)

//import "bufio"

func main() {

	// Print the entire slice of arguments

	// Check if at least one argument was provided (beyond the program name)
	// if len(os.Args) > 1 {
	// 	fmt.Println("The first argument you provided is:", os.Args[1])
	// }

	test, err := os.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println(err)
		return
	}
	raw := string(test)

	// detect 7-bit or 8-bit
	var kind int = 7
	if (len(raw)-1)%8 == 0 {
		kind = 8
	}

	// split into groups of 7/8
	letters := []int64{}
	for i := 0; i < len(raw)-1; i += kind {
		end := i + kind
		if end > len(raw)-1 {
			end = len(raw) - 1
		}
		tempe, err := strconv.ParseInt(raw[i:end], 2, 0)
		if err != nil {
			fmt.Printf("shat itself during parseint, %v", err)
		}
		letters = append(letters, tempe)

	}
	//fmt.Println("")

	//fmt.Println((len(raw) - 1) % 8)
	//fmt.Println(kind)

	//fmt.Println(letters)
	//out_int := []int{}
	out := []string{}

	for i := range letters {
		out = append(out, string(rune(int(letters[i]))))
	}

	for i := range out {
		fmt.Printf("%v", out[i])
	}
}
