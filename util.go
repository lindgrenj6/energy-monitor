package main

import "log"

// get rid of if err != nil over and over again in some cases, YAY GENERICS
func must[T any](thing T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return thing
}
