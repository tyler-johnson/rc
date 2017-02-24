package rc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ini "gopkg.in/ini.v1"
)

var jsonRegex = regexp.MustCompile("^\\s*{")

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseBool(str string) (value bool, err error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return false, nil
	}
	return false, fmt.Errorf("parsing \"%s\": invalid syntax", str)
}

func parseValue(val string) interface{} {
	if b, err := parseBool(val); err == nil {
		return b
	} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	} else {
		return val
	}
}

func deepSet(data map[string]interface{}, keys []string, val interface{}) {
	head := keys[:len(keys)-1]
	last := keys[len(keys)-1]
	res := data

	for _, k := range head {
		existing := res[k]
		if ex, ok := existing.(map[string]interface{}); ok {
			res = ex
		} else {
			child := make(map[string]interface{})
			res[k] = child
			res = child
		}
	}

	res[last] = val
}

func parseEnv(name string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		value := parts[1]

		if !strings.HasPrefix(key, name+"_") {
			continue
		}

		deepSet(data, strings.Split(strings.TrimPrefix(key, name+"_"), "__"), parseValue(value))
	}

	return data
}

func parseIni(content []byte) (data map[string]interface{}, err error) {
	file, err := ini.Load(content)
	if err != nil {
		return
	}

	data = make(map[string]interface{})

	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			name := make([]string, 0)
			if section.Name() != ini.DEFAULT_SECTION {
				name = append(name, strings.Split(section.Name(), ".")...)
			}
			name = append(name, key.Name())
			deepSet(data, name, parseValue(key.Value()))
		}
	}

	return
}

func parse(content []byte) (data map[string]interface{}, err error) {
	if content == nil {
		return
	}

	if jsonRegex.Match(content) {
		err = json.Unmarshal(content, &data)
		return
	}

	return parseIni(content)
}

func readfile(file string) []byte {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	return d
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func find(name string) (string, bool) {
	var findR func(base string) (string, bool)
	findR = func(base string) (string, bool) {
		full := filepath.Join(base, name)
		_, err := os.Stat(full)
		if err == nil {
			return full, true
		}

		dir := filepath.Dir(base)
		if dir == base {
			return "", false
		}

		return findR(dir)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	return findR(cwd)
}

func merge(a map[string]interface{}, b ...map[string]interface{}) map[string]interface{} {
	if a == nil {
		a = make(map[string]interface{})
	}

	for _, obj := range b {
		if obj == nil {
			continue
		}

		for k, v := range obj {
			aobj, aok := a[k].(map[string]interface{})
			bobj, bok := v.(map[string]interface{})

			if !aok || !bok {
				a[k] = v
			} else {
				a[k] = merge(aobj, bobj)
			}
		}
	}

	return a
}
