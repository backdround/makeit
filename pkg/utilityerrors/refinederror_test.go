package utilityerrors_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/backdround/makeit/pkg/utilityerrors"
)


func TestRefinedErrorMessage(t *testing.T) {
	t.Run("Without clarification", func(t *testing.T) {
		e := utilityerrors.NewRefined("Refined error")
		assert.Equal(t, "Refined error", e.Error())
	})

	t.Run("With one clarification", func(t *testing.T) {
		e := utilityerrors.NewRefined("Refined error")
		e = e.Clarify("More accurate info")

		expectedMessage := "Refined error:\n"
		expectedMessage += "  More accurate info"
		assert.Equal(t, expectedMessage , e.Error())
	})

	t.Run("With three clarification", func(t *testing.T) {
		e := utilityerrors.NewRefined("Refined error")
		e = e.Clarify("Top")
		e = e.Clarify("Middle")
		e = e.Clarify("Bottom")

		expectedMessage := "Refined error:\n"
		expectedMessage += "  Top:\n"
		expectedMessage += "    Middle:\n"
		expectedMessage += "      Bottom"
		assert.Equal(t, expectedMessage , e.Error())
	})

}

func TestRefinedErrorIs(t *testing.T) {
	t.Run("RefinedError type", func(t *testing.T) {
		t.Run("Same value", func(t *testing.T) {
			initialError := utilityerrors.NewRefined("Refined error")
			preciseError := initialError.Clarify("Some description")

			assert.True(t, errors.Is(preciseError, preciseError),
				"They are the same error. Must be true!")
		})

		t.Run("Another value", func(t *testing.T) {
			firstError := utilityerrors.NewRefined("Some error")
			secondError := utilityerrors.NewRefined("Another error")

			assert.False(t, errors.Is(firstError, secondError),
				"They are the different error. Must be false!")
		})
	})

	t.Run("Not RefinedError type", func(t *testing.T) {
		refinedError := utilityerrors.NewRefined("Refined error")
		simpleError := errors.New("Simple error")
		assert.False(t, errors.Is(refinedError, simpleError),
			"They are the different error. Must be false!")
	})
}
