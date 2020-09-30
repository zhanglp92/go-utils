package uniq

import (
	"fmt"
	"testing"
)

func Test_UniqData(t *testing.T) {
	c1 := NewUniqCollectionInt(Option{10, func(v interface{}) bool { return v.(int) <= 0 }, true})
	for i := 0; i < 30; i++ {
		c1.Add(int(i % 5))
	}
	fmt.Println(c1.Collection())

	c2 := NewUniqCollectionFloat64(Option{10, func(v interface{}) bool { return v.(float64) <= 0 }, false})
	for i := 0; i < 30; i++ {
		c2.Add(float64((i % 5)))
	}
	fmt.Println(c2.Collection())
}
