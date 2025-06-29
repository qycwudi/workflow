package components

import (
	"errors"

	"github.com/bytedance/sonic"

	"workflow/pkg/constants"
	"workflow/pkg/core"
)

// ComponentFactoryFixed
func ComponentFactoryFixed(e core.WorkflowEngine, nodeType string, inputs core.NodeData) (Component, error) {
	switch nodeType {
	case constants.ComponentStart:
		return NewStartComponent()
	case constants.ComponentEnd:
		return NewEndComponent()
	case constants.ComponentHTTP:
		return NewHTTPComponent(inputs.Custom)
	case constants.ComponentCode:
		return NewCodeComponent(inputs.Custom)
	case constants.ComponentModel:
		return NewModelComponent(inputs.Custom)
	case constants.ComponentCondition:
		jsonConfig, err := sonic.Marshal(inputs.Conditions)
		if err != nil {
			return nil, errors.New("component configuration serialization failed: " + err.Error())
		}
		return NewConditionComponent(jsonConfig)
	case constants.ComponentLoop:
		return NewIterationComponent(e, inputs.BatchFor, inputs.Custom)
	case constants.ComponentDatabase:
		return NewDatabaseComponent(inputs.Custom)
	}
	return nil, errors.New("Component type not found: " + nodeType)
}
