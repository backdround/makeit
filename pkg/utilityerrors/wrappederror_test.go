package utilityerrors_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/backdround/makeit/pkg/utilityerrors"
)

func TestWrappedErrorMessage(t *testing.T) {
	t.Run("Without wrapped error", func(t *testing.T) {
		e := utilityerrors.NewWrapped("Wrapped error")
		assert.Equal(t, "Wrapped error", e.Error())
	})

	t.Run("With wrapped error", func(t *testing.T) {
		e := utilityerrors.NewWrapped("External")
		e = e.Wrap(errors.New("Internal"))

		expectedMessage := "External:\n  Internal"
		assert.Equal(t, expectedMessage, e.Error())
	})
}

func TestWrappedErrorWrap(t *testing.T) {
	e := utilityerrors.NewWrapped("Wrapped error")

	internalError := errors.New("Internal error")
	e = e.Wrap(internalError)

	assert.Equal(t, internalError.Error(), e.Unwrap().Error())
}

func TestWrappedErrorIs(t *testing.T) {
	t.Run("WrappedError type", func(t *testing.T) {
		t.Run("Same value", func(t *testing.T) {
			initialError := utilityerrors.NewWrapped("Wrapped error")
			wrappedError := initialError.Wrap(errors.New("Some error"))

			assert.True(t, errors.Is(wrappedError, initialError),
				"They are the same error. Must be true!")
		})

		t.Run("Another value", func(t *testing.T) {
			firstError := utilityerrors.NewWrapped("Wrapped error")
			secondError := utilityerrors.NewWrapped("Wrapped error")

			assert.False(t, errors.Is(firstError, secondError),
				"They are the different error. Must be false!")
		})
	})

	t.Run("Not WrappedError type", func(t *testing.T) {
		wrappedError := utilityerrors.NewWrapped("Wrapped error")
		simpleError := errors.New("Some error")
		assert.False(t, errors.Is(wrappedError, simpleError),
			"They are the different error. Must be false!")
	})
}
