package commands

type ParamType string

const (
	ParamTypeBool        ParamType = "bool"
	ParamTypeFilePath    ParamType = "filepath"
	ParamTypeString      ParamType = "string"
	ParamTypeStringArray ParamType = "[string]"
	ParamTypeText        ParamType = "text"
	ParamTypeUsername    ParamType = "username"
)

type ParamDesc struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Type  ParamType `json:"type"`
}

type ParamsDesc []ParamDesc

func (p ParamsDesc) ensureTyped(src map[string]any) {
	if src == nil {
		src = map[string]any{}
	}
	for _, param := range p {
		src[param.ID] = param.extractTyped(src)
	}
}

func (p ParamDesc) extractTyped(src map[string]any) any {
	value, ok := src[p.ID]
	switch p.Type {

	case ParamTypeBool:
		if !ok {
			return false
		}
		b, ok := value.(bool)
		if !ok {
			return false
		}
		return b

	case ParamTypeFilePath, ParamTypeString, ParamTypeText, ParamTypeUsername:
		if !ok {
			return ""
		}
		s, ok := value.(string)
		if !ok {
			return ""
		}
		return s

	case ParamTypeStringArray:
		if !ok {
			return []string{}
		}
		arr, ok := value.([]any)
		if !ok {
			return []string{}
		}
		result := make([]string, len(arr))
		for i, src := range arr {
			val, ok := src.(string)
			if !ok {
				result[i] = ""
				continue
			}
			result[i] = val
		}
		return result
	default:
		return nil
	}
}
