/*
author:		Henri Tran
date:		9/24/2025
*/

package main

import ("fmt"
	"bufio"
	"os"
	"strconv"
        )

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println("Error scanning file: ", err)
		os.Exit(1)
	}

	fmt.Println(BintoASCII(scanner.Text()))
}

func BintoASCII(input string) string {
	bitLen := 0
	if len(input)%8 == 0 {
		bitLen = 8
	} else if len(input)%7 == 0 {
		bitLen = 7
	} else {
		fmt.Println("Incorrect length for ASCII values")
		os.Exit(1)
	}

	output := ""
	for i := 0; i < len(input); i+=bitLen {
		curr, err := strconv.ParseInt(input[i: i+bitLen], 2, bitLen+1)
		if err != nil {
			fmt.Println("Unable to parse input: ", err)
			os.Exit(1)
		}
		output += string(curr)
	}
	return output
}
