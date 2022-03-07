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


func TestSequential(t *testing.T) {
	sequential()
}

func TestPrintln(t *testing.T) {
	println()
}

func TestGoroutines_naive(t *testing.T) {
	goroutines_naive()
}

func TestGoroutines_channels(t *testing.T) {
	goroutines_channels()
}

func TestWaitGroups(t *testing.T) {
	waitGroups()
}

func TestSyncDemoUsingAtomic(t *testing.T) {
	syncDemoUsingAtomic()
}

func TestSyncDemoUsingMutex(t *testing.T) {
	syncDemoUsingMutex()
}
