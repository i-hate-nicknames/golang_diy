package routing

import "testing"

func identity(s string) string {
	return s
}

func double(s string) string {
	return s + s
}

type test struct {
	input, expected string
}

func runHandlerTests(t, t *testing.T, name string, h Handler, tests []test) {
	t.Run(name, func(t *testing.T) {
		for _, test := range tests {
			if res := h(test.input); res != test.expected {
				t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
			}
		}
	})
}

func runMWTests(t, t *testing.T, name string, initial Handler, mw Middleware, tests []test) {
	h := mw(initial)
	t.Run(name, func(t *testing.T) {
		for _, test := range tests {
			if res := h(test.input); res != test.expected {
				t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
			}
		}
	})
}

func runRouterTests(t, t *testing.T, name, path string, tests []test) {
	t.Run(name, func(t *testing.T) {
		for _, test := range tests {
			if res := router.Match(path, test.input); res != test.expected {
				t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
			}
		}
	})
}

func TestHandlers(t *testing.T) {
	constTests := []test{
		{"a", "kurwa"},
		{"", "kurwa"},
	}
	runHandlerTests(t, "constantHandler", constantHandler, constTests)

	idTests := []test{
		{"abc", "abc"},
		{"", ""},
	}
	runHandlerTests(t, "identityHandler", identityHandler, idTests)

	abangTests := []test{
		{"a", "a!"},
		{"", "!"},
	}
	runHandlerTests(t, "appendBangHandler", appendBangHandler, abangTests)

	captTests := []test{
		{"a", "A"},
		{"test", "Test"},
		{"!!!", "!!!"},
	}
	runHandlerTests(t, "captHandler", captHandler, captTests)

	captBangTests := []test{
		{"a", "A!"},
		{"test", "Test!"},
		{"!!!", "!!!!"},
	}
	runHandlerTests(t, "captBangHandler", captBangHandler, captBangTests)

	revTests := []test{
		{"a", "a"},
		{"", ""},
		{"kurwa", "awruk"},
	}
	runHandlerTests(t, "revHandler", revHandler, revTests)

	revBangTests := []test{
		{"a", "a!"},
		{"", "!"},
		{"kurwa", "awruk!"},
	}
	runHandlerTests(t, "revBangHandler", revBangHandler, revBangTests)

	revCaptTests := []test{
		{"a", "A"},
		{"", ""},
		{"abc", "Cba"},
	}
	runHandlerTests(t, "revCaptHandler", revCaptHandler, revCaptTests)

	captRevBangTests := []test{
		{"a", "A!"},
		{"", "!"},
		{"abc", "cbA!"},
	}
	runHandlerTests(t, "captRevBangHandler", captRevBangHandler, captRevBangTests)
}

func TestMiddlewares(t *testing.T) {
	constTests := []test{
		{"abc", "kurwa"},
		{"", "kurwa"},
	}

	runMWTests(t, "constMiddleware", identity, constMiddleware, constTests)

	capitalizeTestsId := []test{
		{"abc", "Abc"},
		{"", ""},
	}

	capitalizeTestsDbl := []test{
		{"", ""},
		{"abc", "AbcAbc"},
	}

	runMWTests(t, "capitalizeMiddleware", identity, capitalizeMiddleware, capitalizeTestsId)
	runMWTests(t, "capitalizeMiddleware", double, capitalizeMiddleware, capitalizeTestsDbl)

	bangifyTestsId := []test{
		{"abc", "abc!"},
		{"", "!"},
	}

	bangifyTestsDbl := []test{
		{"", ""},
		{"abc", "abc!abc!"},
	}

	runMWTests(t, "bangifyMiddleware", identity, bangifyMiddleware, bangifyTestsId)
	runMWTests(t, "bangifyMiddleware", double, bangifyMiddleware, bangifyTestsDbl)

	reverseTestsId := []test{
		{"abc", "bca"},
		{"", ""},
	}

	reverseTestsDbl := []test{
		{"", ""},
		{"abc", "bcabca"},
	}

	runMWTests(t, "reverseMiddleware", identity, reverseMiddleware, reverseTestsId)
	runMWTests(t, "reverseMiddleware", double, reverseMiddleware, reverseTestsDbl)

	composedIdsHandler := composeMiddleware(identity, identity, identity, identity)

	composeTestsId := []test{
		{"", ""},
		{"abc", "abc"},
	}

	composeDoubles := composeMiddleware(double, double)

	composeTestsDbl := []test{
		{"", ""},
		{"abc", "abcabcabcabc"},
	}

	// since we used compose middleware to create a handler directly, use identity mw
	// for runMWTests
	identityMiddleware := func(h Handler) Handler { return h }

	runMWTests(t, "composeMiddleware", composedIdsHandler, identityMiddleware, composeTestsId)
	runMWTests(t, "composeMiddleware", composeDoubles, identityMiddleware, composeTestsDbl)

	appender := makeAppender("-a")

	appenderTestsId := []test{
		{"abc", "abc-a"},
		{"", "-a"},
	}

	appenderTestsDbl := []test{
		{"", "-a-a"},
		{"abc", "abc-aabc-a"},
	}

	runMWTests(t, "makeAppender", identity, appender, appenderTestsId)
	runMWTests(t, "makeAppender", double, appender, appenderTestsDbl)

}

func TestUsingMiddlewares(t *testing.T) {
	quadTests := []test{
		{"ab", "abababab"},
		{"", ""},
	}
	runHandlerTests(t, "quadHandler", quadHandler, quadTests)

	captTests := []test{
		{"", ""},
		{"ab", "Ab"},
		{"123", "123"},
	}
	runHandlerTests(t, "captH", captH, captTests)

	capthBangTests := []test{
		{"ab", "Ab!"},
		{"", "!"},
	}
	runHandlerTests(t, "capthBangH", capthBangH, capthBangTests)

	revTests := []test{
		{"abc", "cba"},
		{"", ""},
	}
	runHandlerTests(t, "revH", revH, revTests)

	revBangTests := []test{
		{"ab", "ba!"},
		{"", "!"},
	}
	runHandlerTests(t, "revBangH", revBangH, revBangTests)

	revCaptTests := []test{
		{"ab", "Ba"},
		{"", ""},
	}
	runHandlerTests(t, "revCaptH", revCaptH, revCaptTests)

	captRevBangTests := []test{
		{"abc", "cbA!"},
		{"", "!"},
	}
	runHandlerTests(t, "captRevBangH", captRevBangH, captRevBangTests)

	questionizeTests := []test{
		{"abc", "abc?"},
		{"", "?"},
	}

	h := questionizeMiddleware(identity)
	runHandlerTests(t, "questionizeMiddleware", h, questionizeTests)

	ellipsifyTests := []test{
		{"abc", "abc..."},
		{"", "..."},
	}
	runHandlerTests(t, "ellipsifyHandler", ellipsifyHandler, ellipsifyTests)

	orNotIdTests := []test{
		{"abc", "abcor not?"},
		{"", "or not?"},
	}
	h = orNotMiddleware(identity)
	runHandlerTests(t, "orNotMiddleware", h, orNotIdTests)

	orNotDblTests := []test{
		{"abc", "abcabcornot"},
		{"a", "aaornot"},
		{"", "or not?"},
	}
	h = orNotMiddleware(double)
	runHandlerTests(t, "orNotMiddleware", h, orNotDblTests)
}

func TestRouter(t *testing.T) {
	router.RegisterHandler("/identity", identity)
	idTests := []test{
		{"a", "a"},
		{"", ""},
		{"12345", "12345"},
	}
	runRouterTests(t, "identity handler", "/identity", idTests)

	router.RegisterHandler("/double", double)
	doubleTests := []test{
		{"a", "aa"},
		{"", ""},
		{"12345", "1234512345"},
	}
	runRouterTests(t, "double handler", "/double", doubleTests)

	dblMw := func(h Handler) Handler { return func(s string) string { return h(s + s) } }

	router.UseMiddleware("/doubleMW", dblMw)
	router.RegisterHandler("/doubleMW", identity)
	doubleMWTests := []test{
		{"a", "aa"},
		{"", ""},
		{"12345", "1234512345"},
	}
	runRouterTests(t, "identity handler", "/doubleMW", doubleMWTests)

	router.UseMiddleware("/revcap", reverseMiddleware)
	router.UseMiddleware("/revcap", capitalizeMiddleware)
	router.RegisterHandler("/revcap", identity)
	revcapMWTests := []test{
		{"", ""},
		{"123", "321"},
		{"abcd", "Dcba"},
	}
	runRouterTests(t, "Router reverse, then capitalize", "/revcap", revcapMWTests)

	router.UseMiddleware("/caprev", capitalizeMiddleware)
	router.UseMiddleware("/caprev", reverseMiddleware)
	router.RegisterHandler("/caprev", identity)
	capRevMWTests := []test{
		{"", ""},
		{"123", "321"},
		{"abcd", "dcbA"},
	}
	runRouterTests(t, "Router capitalize, then reverse", "/caprev", capRevMWTests)
}
