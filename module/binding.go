package goioc

import (
	"fmt"
	"reflect"
)

type Binding struct {
	targetType reflect.Type
	ctor       any
	resolved   any
}

func newBinding(targetType reflect.Type, value any, isCtor bool) Binding {

	if isCtor {
		return Binding{
			ctor:       value,
			resolved:   nil,
			targetType: targetType,
		}
	} else {
		return Binding{
			ctor:       nil,
			resolved:   value,
			targetType: targetType,
		}
	}

}

func findBinding(bindings []Binding, t reflect.Type) int {
	return find(bindings, func(b Binding) bool {
		return b.matches(t)
	})
}

func (b Binding) matches(t reflect.Type) bool {
	if b.targetType == t {
		return true
	} else {
		return false
	}
}

func (b *Binding) resolve(bindings []Binding, forceRebind bool) (any, error) {
	if b.resolved != nil && !forceRebind {
		return b.resolved, nil
	}

	if b.ctor == nil {
		return b.resolved, nil
	}

	injections := make([]reflect.Value, 0)
	for i := 0; i < reflect.TypeOf(b.ctor).NumIn(); i++ {
		idx := findBinding(bindings, reflect.TypeOf(b.ctor).In(i))
		if idx == -1 {
			return nil, fmt.Errorf("Binding for %v not found", reflect.TypeOf(b.ctor).In(i))
		}

		v, err := bindings[idx].resolve(bindings, forceRebind)
		if err != nil {
			return nil, err
		}

		injections = append(injections, reflect.ValueOf(v))
	}

	b.resolved = reflect.ValueOf(b.ctor).Call(injections)[0].Interface()
	return b.resolved, nil
}
