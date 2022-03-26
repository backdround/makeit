package makeitdata_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/backdround/makeit/internal/makeitdata"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type readerMock struct {
	mock.Mock
}

func (r readerMock) Read(buffer []byte) (int, error) {
	args := r.Called(buffer)
	return args.Int(0), args.Error(1)
}

func TestNormalWork(t *testing.T) {
	yamlData := dedent.Dedent(`
		version: 1

		tasks:
		  test:
		   script: |
		    echo test1
		    echo test2
	`)
	reader := strings.NewReader(yamlData)
	data, err := makeitdata.New(reader)
	assert.Nil(t, err)

	assert.Equal(t, data.Version, "1")

	assert.Equal(t, 1, len(data.Tasks))

	expectedTestScript := "echo test1\necho test2\n"
	task, ok := data.Tasks["test"]
	assert.True(t, ok)
	assert.Equal(t, expectedTestScript, task.Script)
}

func TestInvalidReader(t *testing.T) {
	reader := readerMock{}

	someReaderError := errors.New("Some reader error")
	reader.On("Read", mock.Anything).Return(0, someReaderError)

	_, err := makeitdata.New(reader)

	if err == nil {
		t.Fatal("Reader error must be proxy received")
	}

	assert.True(t, errors.Is(err, makeitdata.ErrReadData))
	assert.True(t, errors.Is(err, someReaderError))
	reader.AssertExpectations(t)
}

func TestInvalidFormat(t *testing.T) {
	invalidYamlReader := strings.NewReader("\t")

	_, err := makeitdata.New(invalidYamlReader)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, makeitdata.ErrYamlFormat))
}

func TestInvalidData(t *testing.T) {
	expecErrorContainsWords := func(t *testing.T, yamlData string,
		words ...string) {

		t.Helper()

		reader := strings.NewReader(yamlData)

		_, err := makeitdata.New(reader)

		if err == nil {
			t.Fatal("Invalid data structure must be cause a error")
		}

		assert.True(t, errors.Is(err, makeitdata.ErrInvalidData),
			"Expect ErrInvalidData")

		for _, word := range words {
			re := regexp.MustCompile("(?i)" + word)
			matched := re.Match([]byte(err.Error()))

			if !matched {
				t.Fail()
				t.Logf("Expect \"%s\" contained in error", word)
				t.Logf("Error message:\n%s", err.Error())
			}
		}
	}

	t.Run("Root level", func(t *testing.T) {
		t.Run("Must be a map", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				- version: 1

				- tasks:
				  test:
				    script: echo test
			`)
			expecErrorContainsWords(t, yamlData, "root", "map")
		})

		t.Run("Error on unknown fields", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				version: 1

				tasks:
				  test:
				    script: echo test

				aoeuhtn:
			`)
			expecErrorContainsWords(t, yamlData, "root", "aoeuhtn")
		})
	})

	t.Run("Version", func(t *testing.T) {
		t.Run("Is required", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				tasks:
				  test:
				    script: ./test.sh
			`)
			expecErrorContainsWords(t, yamlData, "version", "required")
		})

		t.Run("Must be a string", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				version: []

				tasks:
				  test:
				    script: ./test.sh
			`)
			expecErrorContainsWords(t, yamlData, "version", "string")
		})
	})

	t.Run("Tasks field must be a map", func(t *testing.T) {
		yamlData := dedent.Dedent(`
			version: 1

			tasks: []
		`)
		expecErrorContainsWords(t, yamlData, "tasks", "map")
	})

	t.Run("Task level", func(t *testing.T) {
		t.Run("Error on unknown fields", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				version: 1

				tasks:
				  test:
				    aoeu:
			`)
			expecErrorContainsWords(t, yamlData, "task", "aoeu")
		})

		t.Run("Script is required", func(t *testing.T) {
			yamlData := dedent.Dedent(`
				version: 1

				tasks:
				  test:
			`)
			expecErrorContainsWords(t, yamlData, "task", "script")
		})
	})

	t.Run("Task script must be a string", func(t *testing.T) {
		yamlData := dedent.Dedent(`
			version: 1

			tasks:
			  test:
			    script: {a: 3, b:4}
		`)
		expecErrorContainsWords(t, yamlData, "script", "string")
	})
}
