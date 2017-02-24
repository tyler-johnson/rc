package rc

import "strings"

// Argv is the result of parsing command-line arguments.
type Argv map[string]interface{}

// Get returns the value at deep key
func (am Argv) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	head := keys[:len(keys)-1]
	last := keys[len(keys)-1]
	res := am

	for _, k := range head {
		existing := res[k]
		if ex, ok := existing.(map[string]interface{}); ok {
			res = ex
		} else {
			return nil
		}
	}

	return res[last]
}
