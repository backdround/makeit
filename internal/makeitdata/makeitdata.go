package makeitdata

import (
	"errors"
	"io"

	"github.com/backdround/makeit/pkg/utilityerrors"
	"gopkg.in/yaml.v3"
)

var (
	ErrReadData = utilityerrors.NewWrapped("Unable read data to parse")
	ErrYamlFormat = utilityerrors.NewRefined("Unable parse yaml data")
	ErrInvalidData = utilityerrors.NewRefined("Invalid data structure")
)

type Task struct {
	Script string
}

func (t *Task) UnmarshalYAML(node *yaml.Node) error {

	// Get task fields.
	taskNodes := make(map[string]yaml.Node)
	err := node.Decode(&taskNodes)
	if err != nil {
		return errors.New("Task must be a map")
	}

	// Parse task fields.
	for fieldName, fieldNode := range taskNodes {
		switch fieldName {
		case "script":
			if fieldNode.Value == "" {
				return errors.New("Script must be a string")
			}
			t.Script = fieldNode.Value

		default:
			return errors.New(`Unknown task field: "` + fieldName + `"`)
		}
	}

	//Check parsed data.
	if t.Script == "" {
		return errors.New("Script field is required in task")
	}

	return nil
}

type MakeItData struct {
	Version string
	Tasks map[string]Task
}

func (d *MakeItData) UnmarshalYAML(node *yaml.Node) error {
	// Get root level fields.
	rootNodes := make(map[string]yaml.Node)
	err := node.Decode(&rootNodes)
	if err != nil {
		finalError := ErrInvalidData.Clarify("Root level must be a map")
		finalError.Clarify(err.Error())
		return finalError
	}

	// Parse root level fileds.
	for rootFieldName, rootNode := range rootNodes {
		err := d.parseRootNode(rootFieldName, rootNode)
		if err != nil {
			return err
		}
	}

	// Check parsed data.
	if d.Version == "" {
		return ErrInvalidData.Clarify("Version is required")
	}

	return nil
}

func (d *MakeItData) parseRootNode(fieldName string, node yaml.Node) error {
	switch fieldName {
	case "version":
		if node.Kind != yaml.ScalarNode {
			return ErrInvalidData.Clarify("Version must be a string")
		}
		d.Version = node.Value

	case "tasks":
		// Get tasks map.
		tasksNodes := make(map[string]yaml.Node)
		err := node.Decode(&tasksNodes)
		if err != nil {
			finalError := ErrInvalidData.Clarify("Tasks must be a map")
			finalError = finalError.Clarify(err.Error())
			return finalError
		}

		// Parse each task in map.
		for taskName, taskNode := range tasksNodes {
			var task Task
			err := task.UnmarshalYAML(&taskNode)
			if err != nil {
				clarification := `Unable to parse task ` + taskName + `"`
				finalError := ErrInvalidData.Clarify(clarification)
				finalError = finalError.Clarify(err.Error())
				return finalError
			}

			d.Tasks[taskName] = task
		}

	default:
		return ErrInvalidData.Clarify(`Unknown root field: "` + fieldName + `"`)
	}
	return nil
}

func New(reader io.Reader) (*MakeItData, error) {

	yamlData, err := io.ReadAll(reader)
	if err != nil {
		return nil, ErrReadData.Wrap(err)
	}
	
	makeItData := &MakeItData{
		Tasks: make(map[string]Task),
	}

	err = yaml.Unmarshal(yamlData, makeItData)

	// Wrap only yaml error itself.
	if err == nil || errors.Is(err, ErrInvalidData) {
		return makeItData, err
	} else {
		return makeItData, ErrYamlFormat.Clarify(err.Error())
	}
}
