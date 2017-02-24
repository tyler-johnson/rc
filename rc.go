package rc

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tyler-johnson/minimist"
)

var etc = "/etc"
var win = runtime.GOOS == "windows"
var home string

func init() {
	if win {
		home = os.Getenv("USERPROFILE")
	} else {
		home = os.Getenv("HOME")
	}
}

func Config(appname string, defaults Argv) (Argv, error) {
	return ConfigArgv(appname, defaults, Argv(minimist.Parse(nil)))
}

func ConfigArgv(appname string, defaults map[string]interface{}, argv map[string]interface{}) (data Argv, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	configs := []Argv{defaults}
	conffiles := make([]string, 0)
	env := parseEnv(appname)

	addConfigFile := func(file string) {
		if contains(conffiles, file) {
			return
		}

		content := readfile(file)
		if content != nil {
			d, err := parse(content)
			if err != nil {
				panic(err)
			}

			configs = append(configs, d)
			conffiles = append(conffiles, file)
		}
	}

	if !win {
		addConfigFile(filepath.Join(etc, appname, "config"))
		addConfigFile(filepath.Join(etc, appname+"rc"))
	}

	if home != "" {
		addConfigFile(filepath.Join(home, ".config", appname, "config"))
		addConfigFile(filepath.Join(home, ".config", appname))
		addConfigFile(filepath.Join(home, "."+appname, "config"))
		addConfigFile(filepath.Join(home, "."+appname+"rc"))
	}

	if local, ok := find("." + appname + "rc"); ok {
		addConfigFile(local)
	}

	if ec, ok := env["config"].(string); ok {
		addConfigFile(ec)
	}

	if cc, ok := argv["config"].(string); ok {
		addConfigFile(cc)
	}

	for i := 0; i < len(configs); i++ {
		data = merge(data, configs[i])
	}
	data = merge(data, env)
	data = merge(data, argv)

	confinfo := Argv{
		"configs": conffiles,
	}
	if len(conffiles) > 0 {
		confinfo["config"] = conffiles[len(conffiles)-1]
	}
	data = merge(data, confinfo)

	return
}
