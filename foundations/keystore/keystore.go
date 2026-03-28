package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

type key struct {
	privatePEM string
	publicPEM  string
}

type KeyStore struct {
	store map[string]key
}

func New() KeyStore {
	return KeyStore{
		store: make(map[string]key),
	}
}

func (ks *KeyStore) PrivateKeyPEM(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", fmt.Errorf("kid not found")
	}
	return key.privatePEM, nil
}

func (ks *KeyStore) PublicKeyPEM(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", fmt.Errorf("kid not found")
	}
	return key.publicPEM, nil
}

// this is usually with the auth service
// and done when the auth service boots up
// best practice is to do initial fetch
// from a KMS like AWS KMS, and cache in
// memory using [KeyStore.store]
func (ks *KeyStore) LoadKeys(fsfs fs.FS) {
	fs.WalkDir(fsfs, ".", func(fileName string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsfs.Open(fileName)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("open pem file: %w", err)
		}

		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("read pem: %w", err)
		}

		privatePem := string(pem)
		publicPem, err := toPublicPEM(privatePem)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("private to public pem: %w", err)
		}
		ks.store[strings.Split(fileName, ".pem")[0]] = key{
			privatePEM: privatePem,
			publicPEM:  publicPem,
		}

		return nil
	})
}

func toPublicPEM(privatePEM string) (string, error) {
	pb, _ := pem.Decode([]byte(privatePEM))
	if pb == nil {
		return "", fmt.Errorf("error decoding private pem. it should of PKCS1 or PKCS8")
	}

	var pvtKey any
	pvtKey, err := x509.ParsePKCS1PrivateKey(pb.Bytes)
	if err != nil {
		pvtKey, err = x509.ParsePKCS8PrivateKey(pb.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse private key: %w", err)
		}
	}

	privateKey, ok := pvtKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("invalid private key")
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}
	publicLem := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	var b bytes.Buffer
	err = pem.Encode(&b, &publicLem)
	if err != nil {
		return "", fmt.Errorf("encoding public pem: %w", err)
	}

	return b.String(), nil
}
