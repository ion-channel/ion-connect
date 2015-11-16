package ionconnect

import (
  "fmt"
)


var Debug bool = false

func Debugln(str string) {
  if Debug {
    fmt.Printf("DEBUG: %v\n", str)
  }
}

func Debugf(str string, args ...interface{}) {
  if Debug {
    var format = fmt.Sprintf(str, args...)
    fmt.Printf("DEBUG: %v\n", format)
  }
}
