package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	//now := time.Now().Unix()
	now := time.Now()

	fmt.Println(now)
	fmt.Println(now.Unix())
	fmt.Println(reflect.TypeOf(now.Unix()))
	fmt.Println(now.Year())

}
