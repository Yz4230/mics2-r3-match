package main

import "fmt"

func verbose(s string) {
	if args.Verbose {
		fmt.Println(s)
	}
}
