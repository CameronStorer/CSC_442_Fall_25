// Description: Binary Decoder
// Name: Cameron Storer from team AlaskanBullWorm

package main

// import appropriate libraries
import (
	"fmt"
	"os"
)

// helper method to obtain the user input
func storeInput() []string {
	// store the user's input parameters
	args := os.Args[1:]
	// DEBUG: get number of arguments
/**/	numArgs := len(args)
	fmt.Println("Number of arguments provided: ", numArgs)
	fmt.Println("The arguments are: ")
	// DEBUG: print the arguments
	for i := range args{
		// print each argument
		fmt.Print(args[i], " ")
	}
	return args
	//numArgs := len(args) -1
}


func main(){
	fmt.Println("getting arguments ")
	binary := storeInput()
	fmt.Println(binary)
	
}