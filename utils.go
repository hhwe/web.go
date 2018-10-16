package main

import (
	"crypto/sha256"
	"io"
)

var (
	// todo: add this to config file or os environment
	salt = "123456"
)

// Hash string to string, add salt to hash
func HashSha256(p string) (d string) {
	h := sha256.New()
	io.WriteString(h, p)
	d = string(h.Sum(nil))
	io.WriteString(h, salt)
	io.WriteString(h, d)
	d = string(h.Sum(nil))
	return
}
