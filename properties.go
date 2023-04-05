package goioc

type Properties struct {
	data *map[string]any
}

func newProperties(data *map[string]any) Properties {
	return Properties{
		data: data,
	}
}

func getType[T any](data *map[string]any, key string) (T, bool) {
	v, ok := (*data)[key]

	if !ok {
		return *new(T), ok
	}

	v2, ok := v.(T)

	if !ok {
		return *new(T), ok
	}

	return v2, false
}

func (p *Properties) String(key string) (string, bool) {
	data := p.data
	return getType[string](data, key)
}

func (p *Properties) Int(key string) (int, bool) {
	data := p.data
	return getType[int](data, key)
}

func (p *Properties) Bool(key string) (bool, bool) {
	data := p.data
	return getType[bool](data, key)
}
