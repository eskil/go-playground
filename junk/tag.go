package main

import (
    "fmt"
    "reflect"
)

type T struct {
    f string `one:"1"`
}
func main() {
    t := reflect.TypeOf(T{})
    f, _ := t.FieldByName("f")
    fmt.Println(f.Tag) // one:"1"
    v, ok := f.Tag.Lookup("one")
    fmt.Printf("%s, %t\n", v, ok) // 1, true
}

