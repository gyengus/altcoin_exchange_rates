package main

import "testing"

func TestLoadConfig(t *testing.T) {
	var config = new(Config)
	config = loadConfig()
	if config.RequestTimeout != 20 {
	    t.Fatalf("RequestTimeout is %d, expected 20", config.RequestTimeout)
	}
	//t.Skip()
}

func BenchmarkLoadConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadConfig()
	}
}

