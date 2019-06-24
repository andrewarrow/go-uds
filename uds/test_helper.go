package uds

import "testing"
import "fmt"

func eq(t *testing.T, a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as != bs {
		fmt.Printf("%s: %s != %s\n", t.Name(), as, bs)
		t.Fail()
	}
}
func neq(t *testing.T, a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as == bs {
		fmt.Printf("%s: %s == %s\n", t.Name(), as, bs)
		t.Fail()
	}
}
