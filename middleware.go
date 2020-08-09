// Implement a simple routing system that matches paths exactly,
// and allows to register a handler function for each path

// 1. Handlers
// Handler processes request string and returns a (possibly modified) string
type Handler func (string) string

// DouleHandler doubles its input
func DoubleHandler (in string) string {
	return in + in
}

// 1.1 Simple handlers
// Implement the following handlers:

// a constant handler that ignores input and always returns a constant
// constant("a") -> "kurwa"
// constant("b") -> "kurwa"

// an identity handler that returns input as output
// identity("a") -> "a"
// identity("b") -> "b"

// a handler that appends some data to input:
// h("a") -> "a!"
// h("b") -> "b!"

// 1.2 Advanced handlers
// Implement the following handlers using function definitions:
// captHandler -> capitalizes input
// captBangHandler -> capitalizes input and adds "!" to the end
// revHandler -> reverses input
// revBangHandler -> reverses input and adds "!" to the end
// revCapHandler -> reverses order of letters in every word of the input


// 2. Middlewares

// Middleware takes a handler h1 and returns another handler h2.
// To extend h1 functionality, call h1 internally and use its result,
// to replace h1 completely just do not use it in h2
type Middleware func (Handler) Handler

// Double middleware returns a handler, that doubles its input and then
// calls given handler on the result
func DoubleMiddleware (h Handler) Handler {
	return Handler(func (in string) string {
		return h(in + in)
	}
}

// 2.1 Defining middlewares
// Implement the following middlewares:

// const middleware that returns a handler that ignores its input and always
// returns some constant string

// capitalize middleware that returns a handler that capitalizes input and then calls given handler on the result

// bangify middleware that returns a handler that adds a "!" to the end of its input and then calls given handler on the result

// reverse middleware that returns a handler that reverses its input and then calls given handler on the result

// reverseWords middleware that returns a handler that reverses each word in the input and then calls given handler on the result
// strings with number of words <= 1 are not modified


// middleware factory: come up with a function that returns a middleware. For example,
// a function that takes an integer, and if an integer is between 1 and 10 it resulting middleware
// will call the original handler and then add ! to the result, but if the integer is bigger it will add !!!


// 2.2 Using middlewares

// Use DoubleMiddleware and identity handler to define quad handler: a handler that repeats its input
// four times

// Reimplement handlers from 1.2 using only middlewares from 2.1 and identity handler from 1.1. Do not define
// any new handlers with func and do not use handlers from 1.2

// Since you have to use function calls to define them, you cannot do that in global scope, so put your
// definitions in this function
func UseMiddleware() {
	var capt, captBang, rev, revBang, revCap Handler
	// capt = ...
	// captBang = ...
}

// 2.3 Pre and post middlewares
// All middlewares that we defined before were "pre" middlewares. They first modify the input, and then call
// a handler they were given on the result. This allows building a chain, such that for some middlewares
// p, q, r and some initial handler h, we can apply them in order and produce a final handler h' = p(q(r(h))).
// This handler will perform the following effects on its input:
// (input) -> p -> q -> r -> h -> (output)
// There is another place to add your computation in middleware: after running provided handler. This changes order
// in which the computation will be applied, and it's the only way to run middleware code _after_ provided handler.
// So, for some handler h, if p, q and r are post handlers, and we apply them again in the same order
// h' = p(q(r(h))), the order of computation will be the following:
// (input) -> h -> r -> q -> p -> (output)


// define orNot middleware that returns a handler that first calls provided handler, and then appends
// "or not?" string to the end
// create a handler that adds "..." to the end of its input, and then apply orNot middleware to it to
// obtain a handler that adds "...or not?". Observe how you have to use a "post" middleware in this case


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

// 3.2 Using router
// Use router together with middlewares to check how it all works together
func RouterExample() {
	// todo: initialize router as your concrete implementation
	var router Router
	// todo: define rootHandler as a function
	var rootHandler Handler
	router.RegisterHandler(rootHandler)
	// using root path with some data
	router.Match("/", "some text")
	router.Match("/", "some other text")
	// todo: register more paths and use middlewares
	// for example define /revCapBangify path that
	// reverses, capitalizes and adds "!" to the end of input strings
}



