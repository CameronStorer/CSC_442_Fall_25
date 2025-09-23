// Description: Binary Decoder
// Name: Cameron Storer from Cyberstorm team AlaskanBullWorm

package main

// import appropriate libraries
import (
	"fmt"		// library used for printing
	"os"		// libary used for making system calls
	"io"		// libarary used for taking stdin
	"strconv"	// library used for converting strings
)

// helper method to store the standard input
func storeStdin() string {
	// first, ensure that standard input has been redirected
	stat, _ := os.Stdin.Stat()
	// if not ask the user to do so next time
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Please redirect standard input.")
		fmt.Println("EX: './binary < example.txt'")
		// kill program to allow user to retry
		os.Exit(0)
	}
	// attempt to read the stdin by storing it and error checking
	input, err := io.ReadAll(os.Stdin)
	// if there is an error while reading stdin
	if err != nil {
		// print error message
		fmt.Fprintf(os.Stderr, "Error reading the redirected standard input.\n")
		// kill program
		os.Exit(0)
	}
	// DEBUG: print the stdin
	// fmt.Println(string(input))
	// convert the input into string and return
	return string(input)
}

// helper method to decode the binary
func binaryToASCII(stdin string){
	// determine if the number of bits is even or odd
	// declare some later-needed variables
	number, err := 0,0
	// if EVEN
	if len(stdin) % 2 == 0{
		// convert the binary into numbers composed of 8 bits each
		number, err = strconv.ParseInt(stdin, 2, 8)
	} else { // if ODD
		// convert the binary into numbers composed of 7 bits each
		number, err = strconv.ParseInt(stdin, 2, 8)
	}
	// handle error
	if err != nil {
		// print error message
		fmt.Fprintf(os.Stderr, "Error reading the redirected standard input.\n")
		// kill program
		os.Exit(0)
	}

	fmt.Printf("The binary string '%s' is equal to the number %d\n", binaryString, number)


}

// main function
func main(){
	// collect the redirected stdin if there is any...
	stdin := storeStdin()
	// decode the binary
	binaryToASCII(stdin)


	
}