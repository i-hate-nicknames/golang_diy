package main

import (
	"fmt"
	"strings"
)

// 1. Handlers

// Handler processes request string and returns a (possibly modified) string
type Handler func(string) string

// DouleHandler doubles its input
func doubleHandler(in string) string {
	return in + in
}

// 1.1 Simple handlers
// Implement the following handlers:

// a constant handler that ignores input and always returns a constant
// constant("a") -> "kurwa"
// constant("b") -> "kurwa"

var Constant Handler = func(str string) string {
	const bitch = "kurwa"

	return bitch
}

// an identity handler that returns input as output
// identity("a") -> "a"
// identity("b") -> "b"
// todo: declare and implement

var Identity Handler = func(str string) string {
	return str
}

// a handler that appends some data to input:
// h("a") -> "a!"
// h("b") -> "b!"

var Append Handler = func(s string) string {
	const appendix = "!"
	var str strings.Builder
	in := append([]byte(s), []byte(appendix)...)
	str.Write(in)

	return str.String()
}

// 1.2 Advanced handlers
// Implement the following handlers using function definitions:
// captHandler -> capitalizes input
// captBangHandler -> capitalizes input and adds "!" to the end
// revHandler -> reverses input
// revBangHandler -> reverses input and adds "!" to the end
// revCapHandler -> reverses order of letters in every word of the input

var CapitalizeHandler Handler = func(s string) string {
	return strings.Title(s)
}

func CapitalizeBangHandler(s string, capitalize Handler, append Handler) string {
	capitalized := capitalize(s)
	banged := append(capitalized)

	return banged
}

// revHandler -> reverses input
var ReverseHandler Handler = func(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes) - 1; i < j; i, j = i + 1, j - 1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// revBangHandler -> reverses input and adds "!" to the end
func ReverseBang(s string, reverse, bang Handler) string {
	reversed := reverse(s)

	return bang(reversed)
}

// revCapHandler -> reverses order of letters in every word of the input
func ReverseCapitalize(s string, reverse, capitalize Handler) string {
	reversed := reverse(s)

	return capitalize(reversed)
}


func HandlersTask() {
	const testString = "Kurwa cum back"
	fmt.Println(Constant(testString))
	fmt.Println(Identity(testString))
	fmt.Println(Append(testString))
	fmt.Println(CapitalizeHandler(testString))
	fmt.Println(ReverseHandler(testString))
	fmt.Println(CapitalizeBangHandler(testString, CapitalizeHandler, Append))
	fmt.Println(ReverseBang(testString, ReverseHandler, Append))
	fmt.Println(ReverseBang(testString, ReverseHandler, CapitalizeHandler))
}
