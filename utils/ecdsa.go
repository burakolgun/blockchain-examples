package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func StringToBigIntTuple(s string) (big.Int, big.Int) {
	bx, err := hex.DecodeString(s[:64])

	if err != nil {
		panic(err)
	}

	by, err := hex.DecodeString(s[64:])

	if err != nil {
		panic(err)
	}

	var bix big.Int
	var biy big.Int

	bix.SetBytes(bx)
	biy.SetBytes(by)
	return bix, biy
}

func PublicKeyFromString(publicKeyStr string) *ecdsa.PublicKey {
	x, y := StringToBigIntTuple(publicKeyStr)

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
}

func PrivateKeyFromString(privateKeyStr string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, err := hex.DecodeString(privateKeyStr[:])

	if err != nil {
		panic(err)
	}

	var bi big.Int
	bi.SetBytes(b)

	return &ecdsa.PrivateKey{
		PublicKey: *publicKey,
		D:         &bi,
	}
}

func SignatureFromString(signatureStr string) *Signature {
	x, y := StringToBigIntTuple(signatureStr)

	return &Signature{
		R: &x,
		S: &y,
	}
}
