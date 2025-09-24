package main

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	table := makeTable()
	tablelower := makeTableLower()
	if os.Args[1] == "-e" || os.Args[1] == "-E" {
		key := os.Args[2]
		if key == "" {
			println("Key is required add a key after entering -e")
		}
		for i := 0; i < len(key); i++ {
			if !unicode.IsLetter(rune([]byte(key[i : i+1])[0])) {
				key = strings.ReplaceAll(key, (key[i : i+1]), "")
				i = 0

			}
		}

		for true {
			input := bufio.NewScanner(os.Stdin)
			output := ""
			for input.Scan() {

				str := input.Text()
				if len(str) >= len(key) {
					key = strings.Repeat(key, (len(str)/len(key))+1)
					println(key)
				}

				for i := 0; i < len(str); i++ {
					by := []byte(str[i : i+1])
					a := by[0]
					if (a < 65 || a > 122) || (90 < a && a < 97) {
						output += str[i : i+1]
					} else if unicode.IsLower(rune(a)) {
						t := strings.ToLower(key[i : i+1])
						b := []byte(t)[0]
						output += tablelower[a-97][b-97]
					} else {
						t := strings.ToUpper(key[i : i+1])

						b := []byte(t)[0]
						output += table[a-65][b-65]

					}
				}
				println(output)
				break
			}

		}
		os.Exit(0)

	} else if os.Args[1] == "-d" || os.Args[1] == "-D" {
		key := os.Args[2]
		if key == "" {
			println("Key is required add a key after entering -e")
		}
		for i := 0; i < len(key); i++ {
			if !unicode.IsLetter(rune([]byte(key[i : i+1])[0])) {
				key = strings.ReplaceAll(key, (key[i : i+1]), "")
				i = 0

			}
		}
		for true {
			input := bufio.NewScanner(os.Stdin)
			output := ""
			for input.Scan() {

				str := input.Text()
				if len(str) >= len(key) {
					key = strings.Repeat(key, (len(str)/len(key))+1)
					println(key)
				}

				for i := 0; i < len(str); i++ {
					by := []byte(str[i : i+1])
					a := by[0]
					if (a < 65 || a > 122) || (90 < a && a < 97) {
						output += str[i : i+1]
					} else if unicode.IsUpper(rune([]byte(str[i : i+1])[0])) {
						target := str[i : i+1]
						b := []byte(strings.ToUpper(key[i : i+1]))[0]
						for x := 0; x < 26; x++ {
							if table[b-65][x] == target {
								output += string(rune(x + 65))
								break

							}
						}
					} else {
						target := str[i : i+1]
						b := []byte(strings.ToLower(key[i : i+1]))[0]
						for x := 0; x < 26; x++ {
							if tablelower[b-97][x] == target {
								output += string(rune(x + 97))
								break
							}
						}
					}
				}
				println(output)
				break
			}
		}

	} else {
		println("Please specify encryption or decryption by using -e or -d")
		os.Exit(1)
	}
}

func makeTable() [26][26]string {
	table := [26][26]string{}
	index := 0
	buffer := 0
	for i := 'A'; i <= 'Z'; i++ {
		for j := 'A'; j <= 'Z'; j++ {
			if j+int32(index) <= 'Z' {
				table[i-65][j-65] = string(j + int32(index))

			} else {
				table[i-65][j-65] = string('A' + int32(buffer))
				buffer++
			}
		}
		buffer = 0
		index++
	}
	return table
}
func makeTableLower() [26][26]string {
	table := [26][26]string{}
	index := 0
	buffer := 0
	for i := 'a'; i <= 'z'; i++ {
		for j := 'a'; j <= 'z'; j++ {
			if j+int32(index) <= 'z' {
				table[i-97][j-97] = string(j + int32(index))

			} else {
				table[i-97][j-97] = string('a' + int32(buffer))
				buffer++
			}
		}
		buffer = 0
		index++
	}
	return table
}
