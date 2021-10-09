package main

import (
	"bytes"
	"testing"
)

func TestHashing(t *testing.T) {
	var bk []byte
	got := hash_password([]byte("mynameisaniket"), "ckdjjekk29i2")
	res :=bytes.Compare(got, bk)
	if res!=0 {
		t.Errorf("Correct")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < 5; i++ {
		hash_password([]byte("mynameisaniket"), "ckdjjekk29i2")
	}
}
