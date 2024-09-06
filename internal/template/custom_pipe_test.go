package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired_PanicCasesProvided_Panic(t *testing.T) {
	expectedPanicMessage := "Required field not provided. Check your configuration!"

	for scenario, input := range map[string]interface{}{
		"null pointer": nil,
		"empty string": "",
		"empty map":    map[string]string{},
		"empty slice":  make([]string, 0),
		"empty array":  [0]string{},
	} {
		t.Run(scenario, func(t *testing.T) {
			defer func() {
				if r := recover(); r != expectedPanicMessage {
					t.Errorf("expected panic message: %s, got: %s", expectedPanicMessage, r)
				}
			}()

			required(input)
		})
	}
}

func TestRequired_ValidCases_ValueReturned(t *testing.T) {
	for scenario, input := range map[string]interface{}{
		"pointer": &struct{}{},
		"string":  "hello",
		"map":     map[string]string{"key": "value"},
		"slice":   make([]string, 1),
		"array":   []string{"hello"},
	} {
		t.Run(scenario, func(t *testing.T) {
			result := required(input)

			assert.Equal(t, input, result)
		})
	}
}
