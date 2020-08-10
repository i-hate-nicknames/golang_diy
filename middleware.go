package main

import (
	"fmt"
	"strings"
)

// Implement a simple routing system that matches paths exactly,
// and allows to register a handler function for each path

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

func constHandler(in string) string {
	return "kurwa"
}

// an identity handler that returns input as output
func identityHandler(in string) string {
	return in
}

// a handler that appends some data to input:
func appendBangHandler(in string) string {
	return in + "!"
}

// 1.2 Advanced handlers
// Implement the following handlers using function definitions:
// captHandler -> capitalizes input
func captHandler(in string) string {
	return strings.ToUpper(in)
}

// captBangHandler -> capitalizes input and adds "!" to the end
func captBangHandler(in string) string {
	return strings.ToUpper(in) + "!"
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// reverse order of words in s, assume words are space separated
func reverseWords(s string) string {
	words := strings.Fields(s)
	result := ""
	// inefficient but whatever
	for i := len(words) - 1; i >= 0; i-- {
		result += words[i]
	}
	return result
}

// revHandler -> reverses input
func revHandler(in string) string {
	return reverse(in)
}

// revBangHandler -> reverses input and adds "!" to the end
func revBangHandler(in string) string {
	return reverse(in) + "!"
}

// revCapHandler -> reverses order of letters in every word of the input,
// assume the words are space separated
func revCapHandler(in string) string {
	return reverse(reverseWords(in))
}

func handlersTask() {
	// todo: test your handlers here
	// fmt.Println(constant("test"))
	// fmt.Println(identity("test"))
	// ...
}

// 2. Middleware

// You cannot compose handlers: there should be only one handler for a given
// request. Instead, you create middlewares: a function that takes a handler
// and returns another handler. It can mix in its own functionality, while retaining
// original handler functionality by calling it directly on the input

// Middleware takes a handler h1 and returns another handler h2.
// To extend h1 functionality, call h1 internally and use its result,
// to replace h1 completely just do not use it in h2
type Middleware func(Handler) Handler

// Double middleware returns a handler, that doubles its input and then
// calls given handler on the result
func doubleMiddleware(h Handler) Handler {
	return Handler(func(in string) string {
		return h(in + in)
	})
}

// Question: what is the difference between double handler and double middleware?

// 2.1 Defining middlewares
// Implement the following middlewares:

// const middleware that returns a handler that ignores its input and always
// returns some constant string
func constMw(h Handler) Handler {
	return Handler(func(in string) string {
		return "a"
	})
}

// capitalize middleware that returns a handler that capitalizes input and then calls given handler on the result
func capitalizeMw(h Handler) Handler {
	return Handler(func(in string) string {
		return h(strings.ToUpper(in))
	})
}

// bangify middleware that returns a handler that adds a "!" to the end of its input and then calls given handler on the result
func bangifyMw(h Handler) Handler {
	return Handler(func(in string) string {
		return h(in + "!")
	})
}

// reverse middleware that returns a handler that reverses its input and then calls given handler on the result
func reverseMw(h Handler) Handler {
	return Handler(func(in string) string {
		return h(reverse(in))
	})
}

// reverseWords middleware that returns a handler that reverses each word in the input and then calls given handler on the result
// strings with number of words <= 1 are not modified
func reverseWordsMw(h Handler) Handler {
	return Handler(func(in string) string {
		return h(reverseWords(in))
	})
}

// Middleware factory: implement a function that returns a middleware.
// The function should take a string s and
// return a middleware that will return a handler that will append s
// to every input and pass the result to the original handler
func makeAppender(s string) Middleware {
	return Middleware(func(h Handler) Handler {
		return Handler(func(in string) string {
			return h(in + s)
		})
	})
}

// 2.2 Using middlewares

func middlewareTask() {
	// Use DoubleMiddleware and identity handler to define quad handler: a handler that repeats its input
	// four times

	var quadHandler Handler
	quadHandler = doubleMiddleware(doubleMiddleware(identityHandler))

	// Reimplement handlers from 1.2 using only middlewares from 2.1 and identity handler from 1.1. Do not define
	// any new handlers with func and do not use handlers from 1.2

	var capt, captBang, rev, revBang, revCap Handler
	capt = capitalizeMw(identityHandler)
	captBang = capitalizeMw(bangifyMw(identityHandler))
	rev = reverseMw(identityHandler)
	revBang = reverseMw(bangifyMw(identityHandler))
	revCap = reverseMw(capitalizeMw(identityHandler))

	fmt.Printf("quad in: %s, quad out: %s\n", "test", quadHandler("test"))
	fmt.Printf("capt in: %s, capt out: %s\n", "test", capt("test"))
	fmt.Printf("captBang in: %s, quad out: %s\n", "test", captBang("test"))
	fmt.Printf("rev in: %s, quad out: %s\n", "test", rev("test"))
	fmt.Printf("revBang in: %s, quad out: %s\n", "test", revBang("test"))
	fmt.Printf("revCap in: %s, quad out: %s\n", "test", revCap("test"))

	// Implement questionize middleware using makeAppender. This middleware
	// should append "?" to input before calling passed handler

	var questionizeMw Middleware
	questionizeMw = makeAppender("?")

	var q Handler = questionizeMw(bangifyMw(identityHandler))
	fmt.Printf(q("test string") + "\n")

}

// 2.3 Pre and post middlewares
// All middlewares that we defined before were "pre" middlewares. They first modify the input, and then call
// a handler they were given on the result. This allows building a chain, such that for some middlewares
// p, q, r and some initial handler h, we can apply them in order and produce a final handler h' = p(q(r(h))).
// This handler will perform the following effects on its input:
// (input) -> p -> q -> r -> h -> (output)
// There is another place to add your computation in middleware: after running provided handler. This changes order
// in which the computation will be applied, and it's the only way to run middleware code _after_ provided handler.
// So, for some handler h, if p, q and r are post middlewares, and we apply them again in the same order
// h' = p(q(r(h))), the order of computation will be the following:
// (input) -> h -> r -> q -> p -> (output)

func postMiddlewareTask() {
	// Create a handler that adds "..." to the end of its input

	var ellipsify Handler
	ellipsify = func(in string) string {
		return in + "..."
	}

	// Define orNot middleware that returns a handler that first calls provided handler, and then appends
	// "or not?" string to the end

	var orNotMw Middleware
	orNotMw = Middleware(func(h Handler) Handler {
		return Handler(func(in string) string {
			return h(in) + "or not?"
		})
	})

	// Obtain a handler that adds "...or not?" by using ellipsify with orNotMw.
	// Observe that orNot has to be a "post" middleware in this case

	var doubtfulHandler Handler
	doubtfulHandler = orNotMw(ellipsify)

	fmt.Println(doubtfulHandler("test"))
}

// 3 Router

// Router ties everything together, allowing to register a handler to a concrete path,
// as well as attaching zero or more middlewares to a path that will modify the handler registered to that path
type Router interface {
	// Add a handler to the given path, that should be matched exactly
	// when the path is matched, a registered handler will be used to handle data
	RegisterHandler(path string, h Handler)

	// Add middleware to the given path. A handler for this path will be modified by all
	// the middlewares registered for this path, each in turn.
	// Note: the order of registering middlewares affects the order of their application
	UseMiddleware(path string, mw Middleware)

	// Match given path, and run a handler that is registered under that path
	// return error when there is no handler registered for the givne path
	Match(path string, data string) error
}

// 3.1 Implementing router
// Implement router

type MyRouter struct {
	// todo
}

// todo: implement router interface

// 3.2 Using router
// Use router together with middlewares to check how it all works together
func routerTask() {
	// Initialize router as your concrete implementation

	// todo: uncomment and implement
	// var router Router
	// router = ...

	// Define rootHandler as a function that does some processing

	// todo: uncomment and implement
	// var rootHandler Handler
	// rootHandler = ...

	// Register root handler under "/" path

	// todo: uncomment when rootHandler is defined
	// router.RegisterHandler("/", rootHandler)

	// Test router root path with different data

	// todo: uncomment
	// router.Match("/", "some text")
	// router.Match("/", "some other text")

	// Register a path with middleware: add /revCapBangify path that reverses,
	// capitalizes and adds "!" to the end of input strings

	// todo: uncomment and implement
	// router.UseMiddleware("/revCapBangify", ...)
	// router.UseMiddleware("/revCapBangify", ...)
	// router.RegisterHandler("/revCapBangify", identityHandler)

	// Test /revCapBangify path

	// todo: uncomment
	// router.Match("/revCapBangify", "some text")
	// router.Match("/revCapBangify", "some other text")
}

func main() {
	middlewareTask()
	postMiddlewareTask()
	// routerTask()
}
