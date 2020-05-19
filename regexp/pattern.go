package main

import (
	"exercise/regexp/test"
	"fmt"
	"golang.org/x/text/number"
	"regexp"
	"strconv"
)

func init() {
	fmt.Println("main init")
}

func main() {
	test.Out()
	/*var b []byte
	b = []byte("ABC")
	fmt.Println(b)
	ok, _ := regexp.Match("ABC", b)

	fmt.Println(ok)
	str :="hello"
	fmt.Println(str[4])

	re,_:=regexp.MatchString("123","1223456")
	fmt.Println(re)*/

	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	pat := "[0-9]+.[0-9]+"

	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
		fmt.Println("match found!")
	}

	re, _ := regexp.Compile(pat)
	result := re.FindAllString(searchIn, -1)
	fmt.Println(result)
	str := re.ReplaceAllString(searchIn, "**.*")
	fmt.Println(str)

	f := func(s string) string {
		v, _ := strconv.ParseFloat(s, 64)
		return strconv.FormatFloat(v, 'f', 10, 64)
	}
	str2 := re.ReplaceAllStringFunc(searchIn, f)
	fmt.Println(str2)
	x := 143.66 - 14.55
	fmt.Println(x)
	str3 := number.Decimal(x)
	fmt.Println(str3)
	d := 1129.6

	fmt.Println(d * 100)
	c := 67.6
	fmt.Println(c * 100)
	//fmt.Printf("%")

}
