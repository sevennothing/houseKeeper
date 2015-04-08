/*
 * Copyright 2015
 * Filename:readJson.go
 * Author:caijun.Li  Date:2015-04-08
 * Description:
 * History:
 *         <author>      <time>       <desc>
 *
 */
package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var outJson = map[string]string{}

func ReadFile(filename string) (map[string]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &outJson); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return outJson, nil
}
