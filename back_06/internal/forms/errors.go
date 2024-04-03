package forms

type customError map[string][]string

func (c customError) Add(field, message string) {
	c[field] = append(c[field], message)
}

func (c customError) Get(field string) string {
	es := c[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
