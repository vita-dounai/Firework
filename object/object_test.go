package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	str1 := &String{Value: "1"}
	int1 := &Integer{Value: 1}
	int2 := &Integer{Value: 1}

	if hello1.Hash() != hello2.Hash() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.Hash() != diff2.Hash() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.Hash() == diff1.Hash() {
		t.Errorf("strings with different content have same hash keys")
	}

	if int1.Hash() != int2.Hash() {
		t.Errorf("integers with same value have different hash keys")
	}

	if str1.Hash() == int1.Hash() {
		t.Errorf("objects with different type have same hash keys")
	}
}
