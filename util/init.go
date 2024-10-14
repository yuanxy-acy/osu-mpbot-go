package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func StringToFloat(str string) float32 {
	flo, err := strconv.ParseFloat(str, 32)
	if err != nil {
		fmt.Println("string to float32 fail:", err)
		return -1
	}
	return float32(flo)
}

func FloatToString(val float32) string {
	return strconv.FormatFloat(float64(val), 'f', 1, 32)
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("string to int fail", err)
		return -1
	}
	return i
}

func StringMatch(str, pattern string) bool {
	res, err := regexp.MatchString(pattern, str)
	if err != nil {
		fmt.Println()
		return false
	}
	return res
}

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regStr := "\\s{2,}"                         //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regStr)            //编译正则表达式
	s2 := make([]byte, len(s1))                 //定义字符数组切片
	copy(s2, s1)                                //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spcIndex) > 0 {                     //找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(s2))           //继续在字符串中搜索
	}
	return string(s2)
}

func GetModName(enabledMods string) string {
	modNum, _ := strconv.Atoi(enabledMods)
	if modNum == 0 {
		return "NM "
	}
	var modName = ""
	if modNum >= 536870912 {
		modName += "V2 "
		modNum -= 536870912
	}
	if modNum >= 16416 {
		modName += "PF "
		modNum -= 16416
	}
	if modNum >= 1024 {
		modName += "FL "
		modNum -= 1024
	}
	if modNum >= 576 {
		modName += "NC "
		modNum -= 576
	}
	if modNum >= 256 {
		modName += "HF "
		modNum -= 256
	}
	if modNum >= 64 {
		modName += "DT "
		modNum -= 64
	}
	if modNum >= 32 {
		modName += "SD "
		modNum -= 32
	}
	if modNum >= 16 {
		modName += "HR "
		modNum -= 16
	}
	if modNum >= 8 {
		modName += "HD "
		modNum -= 8
	}
	if modNum >= 2 {
		modName += "EZ "
		modNum -= 2
	}
	if modNum >= 1 {
		modName += "NF "
		modNum -= 1
	}
	if modNum == 0 {
		return modName
	}
	return ""
}
