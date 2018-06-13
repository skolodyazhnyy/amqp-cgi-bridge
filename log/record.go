package log

type R map[string]interface{}

func (r R) With(w map[string]interface{}) R {
	m := make(R)

	if r != nil {
		for k, v := range r {
			m[k] = v
		}
	}

	if w != nil {
		for k, v := range w {
			m[k] = v
		}
	}

	return m
}
