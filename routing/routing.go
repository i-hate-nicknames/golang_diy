package routing

import (
	"strings"
	"unicode"
)

// TODO: move most of the prose to markdown, leave only implementation-related text

// This guide walks through the implementation of a simple routing system, that is often used in
// web development, but is not limited to it. It includes exercises that are automatically tested.
// To check the correctness of your implementation run go test ./... in this directory

// The goal of a routing system is to process incoming requests, and based on some set of rules
// process these requests. Routing system helps decoupling the processing part from the rule part.
// Additionally, it helps with processing organization by adding a powerful concept of middlewares.

// Imagine a simple routing system that sorts balls by color.
// It should put green balls in box G, and red balls in box R.
// Assume no other balls can come in

// Main components of routing system are requests, route matchers and handlers.

// Request is something that comes in our system from the outside. In web development it's
// some representation of an http request, in our ball machine example requests are balls.
// In this guide a request will be a pair of strings: a path and data, which represents some
// client sending given data to a given path.
type Request struct {
	Path, Data string
}

// Route matcher gets a request and checks if this request is matched. I.e. it answers a yes/no
// question, and can be summarized as a function that takes a request and returns a boolean

// In the balls example, it is a part of the system that checks balls for color.
// In this system we have two matchers: one for green balls and one for red balls.

// In this guide, routing rules will be simple strings, that match paths exactly. I.e. a routing
// rule "test" will match all incoming requests that have Path set to "test".
// A more flexible system might want to use a matcher type that will allow adding
// custom logic to the match system

// Matcher type is not used, and is here for illustration purposes. We use strings and match path exactly.
type Matcher func(Request) bool

// A handler is a function that takes a request and then does some work on it. In web development
// it usually operates on http response. In our example we have two handlers: green handler and red hander.
// Green handler takes a ball, and puts it in G box, red handler puts the incoming ball to the R box.
// The important part: handlers are dumb, they do not know the color of incoming ball, they just
// put it where they are programmed to.

// Handler in this guide handlers will operate on the underlying string data, instead on the requests
// themselves, and will return their results as strings. This is done for brevity and is not crucial
// for conceptual understanding of a routing system.
type Handler func(string) string

// Important note: for every rule in the routing system there is exactly one handler. Imagine if you want
// to register multiple handlers. There can be two possible reasons for this:
// 1. You want all those handlers to run. In this case, you will have to make a new handler, that will
// encapsulate all the logic you want. You can make such handler yourself, or abstract the logic to
// middleware (more on those in Middleware section of this guide)
// 2. You want only some of the handlers to run. In this case, the logic that decides if a handler
// is executed should be moved to the rule part, and the handlers themselves grouped differently under
// new rules.

// Lifetime of a routing system can be divided in two phases: configuration time and run time.
// At configuration time, a programmer specifies routing rules, together with processing "attached"
// to each rule.
// At run time, requests come into the routing system, and routing system forwards these
// requests according to the rules it was given at configuration time.

// Note: nothing prevents adding/removing rules at runtime too, but for the vast majority of cases
// it happens at configuration time

// To follow this guide, do not change types of anything. Remove panics with your implementation
// and run tests

// 1. Handlers

// douleHandler doubles its input
// this is an example
func doubleHandler(in string) string {
	return in + in
}

// Simple handlers
// Implement the following handlers:

// a constant handler that ignores input and always returns a constant
// constantHandler("a") -> "kurwa"
// constantHandler("b") -> "kurwa"
func constantHandler(in string) string {
	const PERMANENT = "kurwa"

	return PERMANENT
}

// an identity handler that returns input as output
// identityHandler("a") -> "a"
// identityHandler("b") -> "b"
func identityHandler(in string) string {
	return in
}

// a handler that appends an exclamation mark to input:
// h("a") -> "a!"
// h("b") -> "b!"
func appendBangHandler(in string) string {
	return in + "!"
}

// Advanced handlers
// Implement the following handlers using function definitions:
// Here and further, capitalization is making first letter of the text to be uppercase
// "test" -> "Test"

// Similarity between handlers is intended. To avoid code repetition,
// you may want to call previously defined handlers in your other handlers. There is a better
// way to do this, using middlewares, that is explained in middleware section

// captHandler capitalizes input
func captHandler(in string) string {
	asRunes := []rune(in)
	asRunes[0] = unicode.ToUpper(asRunes[0])

	return string(asRunes)
}

// captBangHandler capitalizes input and adds "!" to the end
func captBangHandler(in string) string {
	return captHandler(in) + "!"
}

// revHandler reverses input
func revHandler(in string) string {
	return reverseString(in)
}

// revBangHandler reverses input and adds "!" to the end
func revBangHandler(in string) string {
	return reverseString(in) + "!"
}

// revCaptHandler reverses the input and then capitalizes it
func revCaptHandler(in string) string {
	panic("not implemented")
}

// revCapBangHandler: capitalizes input, then reverses it and then adds "!" to the end
func captRevBangHandler(in string) string {
	panic("not implemented")
}

// 2. Middlewares

// Handlers are functions that take strings and return strings. Naturally, they can be composed,
// to produce new handlers. You probably notices this when implementing handlers in previous section.
// If you implemented all your handlers from scratch, I advise you to go back and try to rewrite them
// using existing handlers.
// In general, if you have two handlers f and g, and want to create a new handler h that uses first f,
// and then g on the input, you can do it as follows
// h := func (in string) string { return g(f(in)) }

// However, this is cumbersome. Each time you want to add functionality to an existing handler, you
// need to create a new handler.
// The solution is middlewares. Basically, middlewares facilitate creation of handlers, and allow mixing
// their own logic with the logic of an existing handler.
// In routing systems, middlewares act as plugins that you can install for a particular route, but the route
// still needs a handler.
// Middlewares are one level above handlers: while handlers operate on requests, middlewares operate on handlers.
// Middleware takes a handler, and modifies it to add it's own logic. Since it cannot inject its own code
// inside a handler, the only solution is to create a  _new_ handler, that will call the original handler,
// plus do some work.
// This way, the resulting handler can retain the original functionality (from the passed handler),
// while also add its own functionality.

// Middleware takes a handler h and returns another handler.
// To extend h1 functionality, call h internally and use its result,
// to replace h1 completely just do not use h at all
type Middleware func(Handler) Handler

// Double middleware returns a handler, that doubles its input and then
// calls given handler on the result
func doubleMiddleware(h Handler) Handler {
	return Handler(func(in string) string {
		return h(in + in)
	})
}

// Question: what is the difference between double handler and double middleware?
// Look at their type signatures, and consider how each can be used.

// Defining middlewares
// Implement the following middlewares:

// const middleware that returns a handler that ignores its input and always
// returns string "kurwa"
func constMiddleware(h Handler) Handler {
	panic("not implemented")
}

// capitalizeMiddleware returns a handler that capitalizes input and then calls given handler on the result
func capitalizeMiddleware(h Handler) Handler {
	panic("not implemented")
}

// bangifyMiddleware returns a handler that adds a "!" to the end of its input and then calls given handler on the result
func bangifyMiddleware(h Handler) Handler {
	panic("not implemented")
}

// reverseMiddleware returns a handler that reverses its input and then calls given handler on the result
func reverseMiddleware(h Handler) Handler {
	panic("not implemented")
}

// composeMiddleware takes many handlers, returns a handler that takes a string,
// and then applies first handler to this string, then applies second handler to the result, and so on...
// Remember when someone said you can only use one handler per route? Pfff.
func composeMiddleware(hs ...Handler) Handler {
	panic("not implemented")
}

// Middleware factory: implement a function that returns a middleware.
// The function should take a string s and
// return a middleware that will return a handler that will append s
// to every input and pass the result to the original handler
func makeAppender(s string) Middleware {
	// todo: implement
	panic("not implemented")
}

// Middlewares are used by routing system, to mix in functionality. However, to understand them
// better it might be helpful to run them directly.
// Let's make some handlers by manually calling middlewares. Functionality of these handlers
// resembles what we already implemented in the handlers section. This is intentional, to show
// a different and more compact way handlers can be defined, when common functionality is abscracted
// into middlewares.
// Consequently, when defining handlers in this section you SHOULD NOT use any other existing handlers or
// helper functions, or define any new handlers the regular way.
// The only handler you need is identityHandler, everything else should be implemented by existing middlewares

// top-level definitions so that tests can see them
var quadHandler, captH, captBangH, revH, revBangH, revCaptH, captRevBangH Handler
var questionizeMiddleware Middleware

// uncomment and implement all assignments in this function
func usingMWTask() {
	// Use DoubleMiddleware and identity handler to define quad handler:
	// a handler that repeats its input four times

	// todo: uncomment and implement
	// quadHandler = ...

	// Capitalize handler capitalizes its input
	// todo: uncomment and implement
	// captH = ...

	// Capitalize Bang handler capitalizes its input  and adds "!" to the end
	// todo: uncomment and implement
	// captBangH = ...

	// Reverse handler reverses its input
	// todo: uncomment and implement
	// revH = ...

	// Reverse Bang handler reverses its input and then adds "!" to the end
	// todo: uncomment and implement
	// revBangH = ...

	// Reverse Capitalize handler reverses its input and then capitalizes it
	// todo: uncomment and implement
	// revCaptH = ...

	// Capitalize Reverse Bang handler capitalizes its input, then reverses it
	// and then adds "!" to the end
	// captRevBangH = ...

	// Implement questionize middleware using makeAppender. This middleware
	// should append "?" to input before calling passed handler
	// todo: uncomment and implement
	// questionizeMiddleware = ...
}

// Pre and post middlewares
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
// The third option might be an "around" middleware, that runs something both before and after everything
// it was passed. It's somewhat more rare than the other two. One example might be logging the time it
// takes to run passed handler

// So here comes another application of middlewares: you can import a middleware from a library, and you
// don't have to think about order of computation, it's already encapsulated in the middleware.
// If you import raw functions from a library and make your own handlers with them, you would have to
// decide where the logic should go.

var ellipsifyHandler, doubtfulHandler Handler
var orNotMiddleware Middleware

func postMiddlewareTask() {
	// Create a handler that adds "..." to the end of its input
	// todo: uncomment and implement
	// ellipsifyHandler = ...

	// orNot middleware returns a handler that first calls provided handler, and then appends
	// "or not?" string to the end
	// todo: uncomment and implement
	// orNotMiddleware = ...

	// Obtain a handler that adds "...or not?" by using ellipsify with orNotMw.
	// Observe that orNot has to be a "post" middleware in this case
	// todo: uncomment and implement
	// doubtfulHandler = ...
}

// To sum up
// A middleware is a function that takes a handler and constructs another handler
func exampleMiddleware(h Handler) Handler {
	// the logic here will be executed at configuration time, when this middleware is registered
	// Typically, not much happens here, because we only have access to the original handler,
	// and to global scope
	return Handler(func(request string) string {
		// The code here will be executed at run time, when a request is matched to the route this
		// middleware is registered on
		// the request may have been modified by previous middlewares, but as a middleware builder
		// this should not concern us.
		// The typical logic here would be to put some information inside request or some request context,
		// instead of modifying the request
		// Here, we are in charge of what goes down the chain: we can ignore the request, pass it as it is,
		// modify it, write to log. We can prevent the original handler from running, if we want to.

		// In this case we do nothing, effectively making our middleware redundant.
		// Make sure you understand why "do nothing" means "return h(request)" and not "return request"
		return h(request)
	})
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
	// Middleware that is installed first should run first
	UseMiddleware(path string, mw Middleware)

	// Match given path, and run a handler that is registered under that path
	// return error when there is no handler registered for the given path
	// return result of running registered handler together with all registered
	// middlewares, for the given path
	Match(request Request) (string, error)
}

// 3.1 Implementing router
// Implement router

var router Router

// todo: define your own type that implements Router interface

func routerTask() {

	// todo: uncomment and assign router to an instance of your type
	// that implements Router interface
	// router = ...
}

func reverseString(in string) string {
	var sb strings.Builder
	runeSlice := []rune(in)
	len := len(runeSlice)

	for i := len - 1; i >= 0; i-- {
		sb.WriteRune(runeSlice[i])
	}

	return sb.String()
}
