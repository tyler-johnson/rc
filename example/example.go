package main

import (
	"fmt"

	"github.com/tyler-johnson/rc"
)

func main() {
	conf, err := rc.Config("myapp", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
}
