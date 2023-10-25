package arr

import (
	"log"
	"testing"
)

func TestCheckArrayType(t *testing.T) {
	t.Run("check array type test with string", func(t *testing.T) {
		arrayType := NewCheckArrayType(10)
		arrayType.Add("1")
		if !arrayType.Check("1") {
			t.Errorf("1 should be in array type,but it is not in.")
		}
		if arrayType.Check("2") {
			t.Errorf("2 should not be in array type,but it is in.")
		}
		if arrayType.Check(1) {
			t.Errorf("int 1 should not be in array type,but it is in.")
		}
		if arrayType.Len() != 1 {
			t.Errorf("array type len should be 1 ,but get array type len is %d.", arrayType.Len())
		}
		arrayType.AddMutil("2", "3", "4")
		if arrayType.Len() != 4 {
			t.Errorf("array type len should be 4 ,but get array type len is %d.", arrayType.Len())
		}
		for i, v := range arrayType.Values() {
			log.Println(v)
			if !arrayType.Check(i) {
				t.Errorf("%v should be in array type,but it is not in", i)
			}
		}
	})

	t.Run("check array type test with uint", func(t *testing.T) {
		arrayType := NewCheckArrayType(10)
		var one uint = 1
		var two uint = 2
		var three uint = 3
		var four uint = 4
		arrayType.Add(one)
		if !arrayType.Check(one) {
			t.Errorf("1 should be in array type,but it is not in.")
		}
		if arrayType.Check(two) {
			t.Errorf("2 should not be in array type,but it is in.")
		}
		if arrayType.Len() != 1 {
			t.Errorf("array type len should be 1 ,but get array type len is %d.", arrayType.Len())
		}
		arrayType.AddMutil(two, three, four)
		if arrayType.Len() != 4 {
			t.Errorf("array type len should be 4 ,but get array type len is %d.", arrayType.Len())
		}
		for i, v := range arrayType.Values() {
			log.Println(v)
			if !arrayType.Check(i) {
				t.Errorf("%v should be in array type,but it is not in", i)
			}
		}
	})

	t.Run("check array type test with int", func(t *testing.T) {
		arrayType := NewCheckArrayType(10)
		arrayType.Add(1)
		if !arrayType.Check(1) {
			t.Errorf("1 should be in array type,but it is not in.")
		}
		if arrayType.Check(2) {
			t.Errorf("2 should not be in array type,but it is in.")
		}
		if arrayType.Check("1") {
			t.Errorf("string 1 should not be in array type,but it is in.")
		}
		if arrayType.Len() != 1 {
			t.Errorf("array type len should be 1 ,but get array type len is %d.", arrayType.Len())
		}
		arrayType.AddMutil(2, 3, 4)
		if arrayType.Len() != 4 {
			t.Errorf("array type len should be 4 ,but get array type len is %d.", arrayType.Len())
		}
		for i, v := range arrayType.Values() {
			log.Println(v)
			if !arrayType.Check(i) {
				t.Errorf("%v should be in array type,but it is not in", i)
			}
		}
	})
}
