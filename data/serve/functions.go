package main

import (
	"encoding/json"
	"fmt"
)

func print_data(data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(data)
		return
	}
	if len(b) < 300 {
		fmt.Println(string(b))
		return
	}
	fmt.Println(string(b[:8130]) + "\n...")
}
