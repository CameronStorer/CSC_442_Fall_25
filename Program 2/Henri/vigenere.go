/*
author:		Henri Tran
date:		9/24/2025
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"unicode"
)

func main() {
	// Parse command-line flags
	encryptFlag := flag.Bool("e", false, "Encrypt input")
	decryptFlag := flag.Bool("d", false, "Decrypt input")
	flag.Parse()

	if (!*encryptFlag && !*decryptFlag) || (*encryptFlag && *decryptFlag) {
		fmt.Fprintln(os.Stderr, "Usage: vigenere -e key OR vigenere -d key")
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Missing key")
		os.Exit(1)
	}

	key := normalizeKey(flag.Arg(0))
	if len(key) == 0 {
		fmt.Fprintln(os.Stderr, "Key must contain at least one alphabetic character")
		os.Exit(1)
	}

	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if *encryptFlag {
			fmt.Println(vigenere(line, key, true))
		} else {
			fmt.Println(vigenere(line, key, false))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
	}
}

// normalizeKey removes spaces and keeps only alphabetic characters
func normalizeKey(key string) []rune {
	var result []rune
	for _, r := range key {
		if unicode.IsLetter(r) {
			result = append(result, unicode.ToUpper(r))
		}
	}
	return result
}

// vigenere applies the Vigenère cipher (encrypt or decrypt)
func vigenere(text string, key []rune, encrypt bool) string {
	var result []rune
	keyLen := len(key)
	j := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			shift := int(key[j%keyLen] - 'A')
			if !encrypt {
				shift = -shift
			}
			result = append(result, shiftRune(r, shift))
			j++
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// shiftRune shifts a letter while preserving case
func shiftRune(r rune, shift int) rune {
	var base rune
	if unicode.IsUpper(r) {
		base = 'A'
	} else {
		base = 'a'
	}
	return (r-base+rune(shift)+26)%26 + base
}
