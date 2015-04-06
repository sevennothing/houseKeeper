package main

import (
	"../models/tools"
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

	var tool *tools.HkTool
	randStr, _ := tool.GenerateRandNumber(6)
	fmt.Println(randStr)

	id := tool.GenerateID16()
	fmt.Println(id)

	id = tool.GenerateID32()
	fmt.Println(id)
}
