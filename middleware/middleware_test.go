package middleware

import "testing"

func TestHandlers(t *testing.T) {

	type handlerTest struct {
		input, expected string
	}

	runHandlerTests := func(name string, h Handler, tests []handlerTest) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				if res := h(test.input); res != test.expected {
					t.Errorf("input: %s, expected: %s, got: %s", test.input, test.expected, res)
				}
			}
		})
	}

	constTests := []handlerTest{
		{"a", "kurwa"},
		{"", "kurwa"},
	}
	runHandlerTests("constantHandler", constantHandler, constTests)

	idTests := []handlerTest{
		{"abc", "abc"},
		{"", ""},
	}
	runHandlerTests("identityHandler", identityHandler, idTests)

	abangTests := []handlerTest{
		{"a", "a!"},
		{"", "!"},
	}
	runHandlerTests("appendBangHandler", appendBangHandler, abangTests)

	captTests := []handlerTest{
		{"a", "A"},
		{"test", "Test"},
		{"!!!", "!!!"},
	}
	runHandlerTests("captHandler", captHandler, captTests)

	captBangTests := []handlerTest{
		{"a", "A!"},
		{"test", "Test!"},
		{"!!!", "!!!!"},
	}
	runHandlerTests("captBangHandler", captBangHandler, captBangTests)

	revTests := []handlerTest{
		{"a", "a"},
		{"", ""},
		{"kurwa", "awruk"},
	}
	runHandlerTests("revHandler", revHandler, revTests)

	revBangTests := []handlerTest{
		{"a", "a!"},
		{"", "!"},
		{"kurwa", "awruk!"},
	}
	runHandlerTests("revBangHandler", revBangHandler, revBangTests)

	revCaptTests := []handlerTest{
		{"a", "A"},
		{"", ""},
		{"abc", "Cba"},
	}
	runHandlerTests("revCaptHandler", revCaptHandler, revCaptTests)

	captRevBangTests := []handlerTest{
		{"a", "A!"},
		{"", "!"},
		{"abc", "cbA!"},
	}
	runHandlerTests("captRevBangHandler", captRevBangHandler, captRevBangTests)
}
