package variables

import "strings"

type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type Variables map[string]string

func (v Variables) Define(key, val string) {
	v["<"+key+">"] = val
}

func (v Variables) Replace(src string) string {
	for key, val := range v {
		src = strings.ReplaceAll(src, key, val)
	}
	return src
}
func (v Variables) ReplaceSlice(src []string) []string {
	replaced := make([]string, len(src))
	for i, str := range src {
		replaced[i] = v.Replace(str)
	}
	return replaced
}

func (v Variables) Clone() Variables {
	newV := make(map[string]string, len(v))
	for key, value := range v {
		newV[key] = value
	}
	return newV
}

func (v Variables) Update(newV Variables) (changed bool) {
	for k := range v {
		if _, ok := newV[k]; !ok {
			changed = true
			delete(v, k)
		}
	}

	for key, value := range newV {
		if v[key] != value {
			changed = true
			v[key] = value
		}
	}

	return
}
