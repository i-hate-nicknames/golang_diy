package main

import "strings"

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

func handlersTask() {
	// todo: test your handlers here
	// fmt.Println(constant("test"))
	// fmt.Println(identity("test"))
	// ...
}
