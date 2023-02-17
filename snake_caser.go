package opensearchutil

type SnakeCaser struct{}

func NewSnakeCaser() *SnakeCaser {
	return &SnakeCaser{}
}

func (s SnakeCaser) TransformFieldName(name string) (string, error) {
	return s.toSnakeCase(name), nil
}

func (s SnakeCaser) toSnakeCase(name string) string {
	nameBytes := []byte(name)
	nameLen := len(nameBytes)

	var result []byte
	for i := 0; i < nameLen; i++ {
		c := nameBytes[i]
		prevCharUnderscore := i > 0 && nameBytes[i-1] == '_'
		if c >= 'A' && c <= 'Z' {
			if len(result) > 0 {
				if i > 0 &&
					!(nameBytes[i-1] >= 'A' && nameBytes[i-1] <= 'Z') && // Previous character not uppercase?
					!prevCharUnderscore {
					result = append(result, '_')
				}
			}
			result = append(result, c+32)
		} else if c >= '0' && c <= '9' {
			if len(result) > 0 {
				if i > 0 &&
					!(nameBytes[i-1] >= '0' && nameBytes[i-1] <= '9') && // Previous character not a number?
					!prevCharUnderscore {
					result = append(result, '_')
				}
			}
			result = append(result, c)
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}
