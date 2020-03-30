package main

import "fmt"

func sandbox(f func() error) (err error) {
	defer func() {
		if r := recover(); err != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	err = f()

	return
}
