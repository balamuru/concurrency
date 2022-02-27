package main

import "testing"


func TestHelloName(t *testing.T) {
want := "https://www.example.com : 1256"
got := responseSizeHelper("https://www.example.com")
println(got)
    if want != got  {
        t.Fatalf("wanted %s but got %s", want, got)
    }
}
