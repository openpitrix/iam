// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package base58

import (
	"crypto/sha256"
	"math/big"
	"strings"
)

const (
	base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func base58Hash(ba []byte) []byte {
	sha := sha256.New()
	sha2 := sha256.New() // hash twice
	ba = sha.Sum(ba)
	return sha2.Sum(ba)
}

func EncodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}

	// Expected size increase from base58Table conversion is approximately 137%, use 138% to be safe
	ri := len(ba) * 138 / 100
	ra := make([]byte, ri+1)

	x := new(big.Int).SetBytes(ba) // ba is big-endian
	x.Abs(x)
	y := big.NewInt(58)
	m := new(big.Int)

	for x.Sign() > 0 {
		x, m = x.DivMod(x, y, m)
		ra[ri] = base58Table[int32(m.Int64())]
		ri--
	}

	// Leading zeroes encoded as base58Table zeros
	for i := 0; i < len(ba); i++ {
		if ba[i] != 0 {
			break
		}
		ra[ri] = '1'
		ri--
	}
	return ra[ri+1:]
}

func DecodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}

	x := new(big.Int)
	y := big.NewInt(58)
	z := new(big.Int)
	for _, b := range ba {
		v := strings.IndexRune(base58Table, rune(b))
		z.SetInt64(int64(v))
		x.Mul(x, y)
		x.Add(x, z)
	}
	xa := x.Bytes()

	// Restore leading zeros
	i := 0
	for i < len(ba) && ba[i] == '1' {
		i++
	}
	ra := make([]byte, i+len(xa))
	copy(ra[i:], xa)
	return ra
}

func EncodeBase58Check(ba []byte) []byte {
	// add 4-byte hash check to the end
	hash := base58Hash(ba)
	ba = append(ba, hash[:4]...)
	ba = EncodeBase58(ba)
	return ba
}

func DecodeBase58Check(ba []byte) bool {
	ba = DecodeBase58(ba)
	if len(ba) < 4 || ba == nil {
		return false
	}

	k := len(ba) - 4
	hash := base58Hash(ba[:k])
	for i := 0; i < 4; i++ {
		if hash[i] != ba[k+i] {
			return false
		}
	}
	return true
}
