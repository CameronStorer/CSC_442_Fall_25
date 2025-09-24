package main

import (
	"fmt"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	data := ""
	_, err := fmt.Fscanf(os.Stdin, "%s", &data)
	if err != nil {
		return
	}
	//change later on linux

	fmt.Print(decode(data))

}

func decode(s string) string {
	seven := false
	eight := false
	result := ""

	if len(s)%8 == 0 {
		eight = true
	}
	if len(s)%7 == 0 {
		seven = true
	}
	if seven && eight {
		seven = false
		eight = true

	} else if !(seven || eight) {
		seven = true
		eight = false
	}

	if seven {
		for i := 0; i < len(s)-1; i += 7 {
			currbyte := ""
			for j := i; j < i+7; j++ {
				currbyte += string(s[j])
			}
			result += string(rune(parsebinary(currbyte, 7)))

		}
	}
	if eight {

		for i := 0; i < len(s)-1; i += 8 {
			currbyte := ""
			for j := i; j < i+8; j++ {
				currbyte += string(s[j])
			}
			result += string(rune(parsebinary(currbyte, 8)))

		}
	}
	return result

}
func parsebinary(binary string, mode int) int {
	if mode == 7 {
		out := 0
		for i := 0; i < len(binary); i++ {
			if binary[i:i+1] == "1" {
				out += 1 << (7 - (i + 1))
			}
		}

		return out
	}
	if mode == 8 {
		out := 0
		for i := 0; i < len(binary); i++ {
			if binary[i:i+1] == "1" {
				out += 1 << (8 - (i + 1))
			}
		}
		return out
	}
	return 0
}
