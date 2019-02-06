package cmd

import (
	"encoding/json"
	"fmt"
)

// PPrint - takes a blank interface and attempts to output json
func PPrint(o interface{}) {
	j, e := json.Marshal(o)
	if e != nil {
		fmt.Println(e.Error())
	}
	if string(j) == "null" {
		return
	}
	fmt.Println(string(j))
}
