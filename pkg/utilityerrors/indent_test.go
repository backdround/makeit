package utilityerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	t.Run("Single line input", func(t *testing.T) {
		inputData := "some data"
		indentedData := indent(inputData, "  ", 1)

		expectedLine := "  some data"
		assert.Equal(t, expectedLine, indentedData)
	})

	t.Run("Multiline input", func(t *testing.T) {
		inputData := "some data\nanother data"
		indentedData := indent(inputData, "  ", 1)

		expectedLine := "  some data\n  another data"
		assert.Equal(t, expectedLine, indentedData)
	})

	t.Run("Empty line input", func(t *testing.T) {
		inputData := "some data\n\n"
		indentedData := indent(inputData, "  ", 1)

		expectedLine := "  some data\n\n"
		assert.Equal(t, expectedLine, indentedData)
	})

	t.Run("Count of times", func(t *testing.T) {
		t.Run("Zero", func(t *testing.T) {
			inputData := "some data"
			indentedData := indent(inputData, "  ", 0)

			expectedLine := "some data"
			assert.Equal(t, expectedLine, indentedData)
		})

		t.Run("Three", func(t *testing.T) {
			inputData := "some data"
			indentedData := indent(inputData, " ", 3)

			expectedLine := "   some data"
			assert.Equal(t, expectedLine, indentedData)
		})
	})
}
