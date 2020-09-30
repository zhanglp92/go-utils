package xmath

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_MaximinInt(t *testing.T) {
	var a int = 10
	MaximinInt(&a, 20, 0, MinBorderLine)
	fmt.Println("a = ", a)

	var b int64 = 10
	MaximinInt64(&b, 20, 0, MinBorderLine)
	fmt.Println("b = ", b)
}

func Test_CallC(t *testing.T) {
	callC()
}

func Test_interface(t *testing.T) {
	var a interface{} = int(1)
	var b interface{} = int(2)
	aType := reflect.TypeOf(a)
	fmt.Println(aType, b)
}
