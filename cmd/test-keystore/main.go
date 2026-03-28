package main

import (
	"fmt"
	"os"

	"github.com/vishal2098govind/service/foundations/keystore"
)

func main() {
	ks := keystore.New()
	ks.LoadKeys(os.DirFS("zarf/keys/"))
	pb, err := ks.PublicKeyPEM("8234c8c5-0508-4301-bfdf-da1d515399a1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pb)
}
