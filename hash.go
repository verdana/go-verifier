package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// Hasher ...
type Hasher struct {
	hashType string
	handler  hash.Hash
}

var hashes = map[string]crypto.Hash{
	"md5":    crypto.MD5,    // 2
	"sha1":   crypto.SHA1,   // 3
	"sha256": crypto.SHA256, // 5
}

func isValid(algo string) bool {
	for name := range hashes {
		if name == algo {
			return true
		}
	}
	return false
}

// NewHasher Create new instance
func NewHasher() *Hasher {
	return &Hasher{}
}

// SetHashType ...
func (h *Hasher) SetHashType(hashType string) {
	if false == isValid(hashType) {
		fmt.Fprintf(os.Stderr, "Unsupport hash type: %s\n", hashType)
		return
	}
	h.hashType = hashType
}

// CreateHandler ...
func (h *Hasher) CreateHandler() {
	switch h.hashType {
	default:
	case "md5":
		h.handler = md5.New()
	case "sha1":
		h.handler = sha1.New()
	case "sha256":
		h.handler = sha256.New()
	}
}

// HashFile ...
func (h *Hasher) HashFile(fpath string) string {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	defer f.Close()

	hasher.CreateHandler()
	if _, err := io.Copy(h.handler, f); err != nil {
		fmt.Println(err)
		return ""
	}

	format := "%x"
	if *upper == true {
		format = "%X"
	}

	return fmt.Sprintf(format, h.handler.Sum(nil))
}

// Verify ...
func (h *Hasher) Verify(checksum string, fpath string) bool {
	var hashType string
	switch len(checksum) {
	case 32:
		hashType = "md5"
	case 40:
		hashType = "sha1"
	case 64:
		hashType = "sha256"
	}

	if hashType == "" {
		return false
	}

	// Remove the leading character * if exists
	fpath = strings.Trim(fpath, "*")

	hasher.SetHashType(hashType)
	if h.HashFile(fpath) == checksum {
		return true
	}
	return false
}
