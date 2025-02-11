package environment

type Scope struct {
	Variables map[string]interface{}
}

type Environment struct {
	Scopes  []Scope
	Returns []interface{}
}

var Global = &Environment{
	Scopes: []Scope{{
		Variables: make(map[string]interface{}),
	}},
	Returns: []interface{}{},
}

func (e *Environment) PushScope() {
	e.Scopes = append(e.Scopes, Scope{
		Variables: make(map[string]interface{}),
	})
}

func (e *Environment) PopScope() {
	if len(e.Scopes) > 1 {
		e.Scopes = e.Scopes[:len(e.Scopes)-1]
	}
}

func (e *Environment) Set(variableName string, value interface{}) {
	currentScope := &e.Scopes[len(e.Scopes)-1]
	currentScope.Variables[variableName] = value
}

func (e *Environment) Reset(variableName string, value interface{}) {
	for i := len(e.Scopes) - 1; i >= 0; i-- {
		if _, ok := e.Scopes[i].Variables[variableName]; ok {
			e.Scopes[i].Variables[variableName] = value
			return
		}
	}
	e.Scopes[len(e.Scopes)-1].Variables[variableName] = value
}

func (e *Environment) Get(variableName string) (interface{}, bool) {
	for i := len(e.Scopes) - 1; i >= 0; i-- {
		if val, ok := e.Scopes[i].Variables[variableName]; ok {
			return val, ok
		}
	}
	return nil, false
}

func (e *Environment) SetReturn(value interface{}) {
	e.Returns = append(e.Returns, value)
}

func (e *Environment) GetReturn() (interface{}, bool) {
	if len(e.Returns) > 0 {
		returnVal := e.Returns[len(e.Returns)-1]
		e.Returns = e.Returns[:len(e.Returns)-1]
		return returnVal, true
	}

	return nil, false
}
