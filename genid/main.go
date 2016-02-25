package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"hash"
)

func secureId(key []byte, size int, digest func() hash.Hash) (res []byte, err error) {
	data := make([]byte, size)
	_, err = rand.Read(data)
	if err != nil {
		return
	}
	mac := hmac.New(digest, key)
	mac.Write(data)
	res = mac.Sum(data)
	err = nil
	return
}

func main() {
	var key = flag.String("key", "change-me", "Secret key")
	var n = flag.Int("n", 1, "Number of ID to generate")
	var short = flag.Bool("short", false, "Use SHA1 instead of SHA256")

	flag.Parse()
	var dataSize int
	var algo func() hash.Hash

	if *short {
		dataSize = 20
		algo = sha1.New
	} else {
		dataSize = 32
		algo = sha256.New
	}

	for i := 0; i < *n; i++ {
		res, err := secureId([]byte(*key), dataSize, algo)
		if err != nil {
			panic(err)
		}
		fmt.Println(base64.RawURLEncoding.EncodeToString(res))
	}
}
