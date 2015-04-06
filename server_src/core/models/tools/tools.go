/**
*  Copyright 2015,
*  Filename: tools.go
*  Author: caijun.Li
*  Date: 2015-04-05
*  Description:
*  History:
*     <author>   <time>   <desc>
*
 */

package tools

import (
	"crypto/md5"
	//	"crypto/rand"
	//	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type HkTool struct {
	Version string `V1.0`
}

func (this *HkTool) GenerateRandNumber(strlen int) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var maxNum = 0
	for i := 0; i < strlen; i++ {
		maxNum = (maxNum * 10) + 9
	}
	//fmt.Println(maxNum)
	randNumber := r.Intn(maxNum)
	strout := strconv.Itoa(randNumber)
	randNumberLen := len(strout)
	if randNumberLen < strlen {
		zeros := ""
		for j := 0; j < strlen-randNumberLen; j++ {
			zeros += "0"
		}
		strout = zeros + strout
	}
	//	fmt.Println(strout)
	return strout, nil

}

// 产生唯一ID，
// 方法： 得到当前的时间戳，同时生成6位数的随机数字符串，将二者拼接
//        后计算MD5值。因为时间戳加6位随机数几乎是唯一的，所以其对应
//        的MD5值也是唯一的，从而保证唯一性。对于大规模分布式处理系统，
//        唯一性的保证可以引入服务器ID进行拼接，然后再做MD5校验。
// 返回: 当传入参数指定输出16位的MD5值时，输出的字符串为16位，否着输出的
//       字符串为32位。
func (this *HkTool) GenerateID(flag16 bool) string {
	var now = time.Now().Unix()
	nandNumStr, _ := this.GenerateRandNumber(6)
	strNow := strconv.FormatInt(now, 10) + nandNumStr
	fmt.Println(strNow)
	h := md5.New()
	h.Write([]byte(strNow))
	strout := hex.EncodeToString(h.Sum(nil))
	if flag16 {
		strbytes := []byte(strout)
		return string(strbytes[8:24])
	} else {
		return strout
	}

}

func (this *HkTool) GenerateID16() string {
	return this.GenerateID(true)
}
func (this *HkTool) GenerateID32() string {
	return this.GenerateID(false)
}
