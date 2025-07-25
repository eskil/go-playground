package greetings

import (
	"fmt"
	"regexp"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello(name)
	fmt.Println(msg)
	fmt.Println(want.MatchString(msg))
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("`+name+`") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want ""`, msg, err)
	}
}
