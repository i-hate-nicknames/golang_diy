package middleware

import "testing"

func TestHandlers(t *testing.T) {

	type test struct {
		input, expected string
	}

	runHandlerTests := func(name string, h Handler, tests []test) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				if res := h(test.input); res != test.expected {
					t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
				}
			}
		})
	}

	constTests := []test{
		{"a", "kurwa"},
		{"", "kurwa"},
	}
	runHandlerTests("constantHandler", constantHandler, constTests)

	idTests := []test{
		{"abc", "abc"},
		{"", ""},
	}
	runHandlerTests("identityHandler", identityHandler, idTests)

	abangTests := []test{
		{"a", "a!"},
		{"", "!"},
	}
	runHandlerTests("appendBangHandler", appendBangHandler, abangTests)

	captTests := []test{
		{"a", "A"},
		{"test", "Test"},
		{"!!!", "!!!"},
	}
	runHandlerTests("captHandler", captHandler, captTests)

	captBangTests := []test{
		{"a", "A!"},
		{"test", "Test!"},
		{"!!!", "!!!!"},
	}
	runHandlerTests("captBangHandler", captBangHandler, captBangTests)

	revTests := []test{
		{"a", "a"},
		{"", ""},
		{"kurwa", "awruk"},
	}
	runHandlerTests("revHandler", revHandler, revTests)

	revBangTests := []test{
		{"a", "a!"},
		{"", "!"},
		{"kurwa", "awruk!"},
	}
	runHandlerTests("revBangHandler", revBangHandler, revBangTests)

	revCaptTests := []test{
		{"a", "A"},
		{"", ""},
		{"abc", "Cba"},
	}
	runHandlerTests("revCaptHandler", revCaptHandler, revCaptTests)

	captRevBangTests := []test{
		{"a", "A!"},
		{"", "!"},
		{"abc", "cbA!"},
	}
	runHandlerTests("captRevBangHandler", captRevBangHandler, captRevBangTests)
}

func TestMiddlewares(t *testing.T) {
	type test struct {
		input, expected string
	}

	identity := func(s string) string { return s }
	double := func(s string) string { return s + s }

	runMWTests := func(name string, initial Handler, mw Middleware, tests []test) {
		h := mw(initial)
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				if res := h(test.input); res != test.expected {
					t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
				}
			}
		})
	}

	constTests := []test{
		{"abc", "kurwa"},
		{"", "kurwa"},
	}

	runMWTests("constMiddleware", identity, constMiddleware, constTests)

	capitalizeTestsId := []test{
		{"abc", "Abc"},
		{"", ""},
	}

	capitalizeTestsDbl := []test{
		{"", ""},
		{"abc", "AbcAbc"},
	}

	runMWTests("capitalizeMiddleware", identity, capitalizeMiddleware, capitalizeTestsId)
	runMWTests("capitalizeMiddleware", double, capitalizeMiddleware, capitalizeTestsDbl)

	bangifyTestsId := []test{
		{"abc", "abc!"},
		{"", "!"},
	}

	bangifyTestsDbl := []test{
		{"", ""},
		{"abc", "abc!abc!"},
	}

	runMWTests("bangifyMiddleware", identity, bangifyMiddleware, bangifyTestsId)
	runMWTests("bangifyMiddleware", double, bangifyMiddleware, bangifyTestsDbl)

	reverseTestsId := []test{
		{"abc", "bca"},
		{"", ""},
	}

	reverseTestsDbl := []test{
		{"", ""},
		{"abc", "bcabca"},
	}

	runMWTests("reverseMiddleware", identity, reverseMiddleware, reverseTestsId)
	runMWTests("reverseMiddleware", double, reverseMiddleware, reverseTestsDbl)

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

	runMWTests("composeMiddleware", composedIdsHandler, identityMiddleware, composeTestsId)
	runMWTests("composeMiddleware", composeDoubles, identityMiddleware, composeTestsDbl)

	appender := makeAppender("-a")

	appenderTestsId := []test{
		{"abc", "abc-a"},
		{"", "-a"},
	}

	appenderTestsDbl := []test{
		{"", "-a-a"},
		{"abc", "abc-aabc-a"},
	}

	runMWTests("makeAppender", identity, appender, appenderTestsId)
	runMWTests("makeAppender", double, appender, appenderTestsDbl)

}

func TestUsingMiddlewares(t *testing.T) {

	identity := func(s string) string { return s }
	double := func(s string) string { return s + s }

	type test struct {
		input, expected string
	}

	runHandlerTests := func(name string, h Handler, tests []test) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				if res := h(test.input); res != test.expected {
					t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
				}
			}
		})
	}

	quadTests := []test{
		{"ab", "abababab"},
		{"", ""},
	}
	runHandlerTests("quadHandler", quadHandler, quadTests)

	captTests := []test{
		{"", ""},
		{"ab", "Ab"},
		{"123", "123"},
	}
	runHandlerTests("captH", captH, captTests)

	capthBangTests := []test{
		{"ab", "Ab!"},
		{"", "!"},
	}
	runHandlerTests("capthBangH", capthBangH, capthBangTests)

	revTests := []test{
		{"abc", "cba"},
		{"", ""},
	}
	runHandlerTests("revH", revH, revTests)

	revBangTests := []test{
		{"ab", "ba!"},
		{"", "!"},
	}
	runHandlerTests("revBangH", revBangH, revBangTests)

	revCaptTests := []test{
		{"ab", "Ba"},
		{"", ""},
	}
	runHandlerTests("revCaptH", revCaptH, revCaptTests)

	captRevBangTests := []test{
		{"abc", "cbA!"},
		{"", "!"},
	}
	runHandlerTests("captRevBangH", captRevBangH, captRevBangTests)

	questionizeTests := []test{
		{"abc", "abc?"},
		{"", "?"},
	}

	h := questionizeMiddleware(identity)
	runHandlerTests("questionizeMiddleware", h, questionizeTests)

	ellipsifyTests := []test{
		{"abc", "abc..."},
		{"", "..."},
	}
	runHandlerTests("ellipsifyHandler", ellipsifyHandler, ellipsifyTests)

	orNotIdTests := []test{
		{"abc", "abcor not?"},
		{"", "or not?"},
	}
	h = orNotMiddleware(identity)
	runHandlerTests("orNotMiddleware", h, orNotIdTests)

	orNotDblTests := []test{
		{"abc", "abcabcornot"},
		{"a", "aaornot"},
		{"", "or not?"},
	}
	h = orNotMiddleware(double)
	runHandlerTests("orNotMiddleware", h, orNotDblTests)
}
