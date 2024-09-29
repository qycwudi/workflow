package rolego

/*
		function Filter(msg, metadata, msgType) {
	        ${jsScript}
	     }
*/

const (
	Start string = "start"

	JsFilter string = "jsFilter"
)

var RoleModel = map[string]*JsFilterModule{
	JsFilter: new(JsFilterModule),
}

type JsFilterModule struct {
	Configuration struct {
		JsScript string `json:"jsScript"`
	} `json:"configuration"`
}
