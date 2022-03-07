package main

import "testing"

func TestHelloName(t *testing.T) {
	want := "https://www.example.com : 1256"
	got := responseSizeHelper("https://www.example.com")
	println(got)
	if want != got {
		t.Fatalf("wanted %s but got %s", want, got)
	}

}

func TestGoroutines_naive(t *testing.T) {
	goroutines_naive()
}

func TestGoroutines_channels(t *testing.T) {
	goroutines_channels()
}

func TestWaitGroups(t *testing.T) {
	goroutines_waitGroups()
}

func TestSyncDemoUsingAtomic(t *testing.T) {
	goroutines_syncDemoUsingAtomic()
}

func TestSyncDemoUsingMutex(t *testing.T) {
	goroutines_syncDemoUsingMutex()
}
