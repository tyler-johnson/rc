# rc

The non-configurable configuration loader for lazy people.

A Golang port of [rc](https://github.com/dominictarr/rc).

## Basic Usage

```go
package main

import "github.com/tyler-johnson/rc"

func main() {
  conf := rc.Config("myapp", map[string]interface{}{
    "somedefault": "value"
  })

  // conf is a basic map[string]interface{} type
  fmt.Println(conf["foo"])
}
```
