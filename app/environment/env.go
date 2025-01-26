package environment

type Environment struct {
	Scope []map[string]interface{}
}

var Global = &Environment{
	Scope: []map[string]interface{}{
		make(map[string]interface{}),
	},
}

func (e *Environment) PushScope() {
	e.Scope = append(e.Scope, make(map[string]interface{}))
}

func (e *Environment) PopScope() {
	if len(e.Scope) > 1 {
		e.Scope = e.Scope[:len(e.Scope)-1]
	}
}

func (e *Environment) Set(variableName string, value interface{}) {
	e.Scope[len(e.Scope)-1][variableName] = value
}

func (e *Environment) Get(variableName string) (interface{}, bool) {
	for i := len(e.Scope) - 1; i >= 0; i-- {
		if val, ok := e.Scope[i][variableName]; ok {
			return val, true
		}

	}
	return nil, false
}
