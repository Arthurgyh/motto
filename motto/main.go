// Copyright 2014 dong<ddliuhb@gmail.com>.
// Licensed under the MIT license.
//
// The Motto command line tool
package main

import (
	"fmt"
	"os"

	"github.com/Arthurgyh/motto"
)

func usage() {
	fmt.Println("Usage: otto file.js")
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	vm, value, err := motto.Run(os.Args[1])
	if err != nil {
		fmt.Printf("error ocup %v", err)
	} else {
		if value.IsNull() {
			str, err := value.ToString()
			fmt.Printf("result: %s", str)
			_ = err
		}
	}
	_ = vm
}
