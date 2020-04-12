package tracing

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"golang.org/x/crypto/hkdf"
)

type Person struct {
	TracingKey []byte
}

func GeneratePerson() (Person, error) {
	tracingKey := make([]byte, 32)
	_, err := rand.Read(tracingKey)
	if err != nil {
		return Person{}, fmt.Errorf("generating tracing key: %w", err)
	}

	p := Person{
		tracingKey,
	}

	return p, nil
}

func (p Person) DailyTracingKey(dailyNumber uint32) (DailyTracingKey, error) {
	header := []byte("CT-DTK")
	dailyNumberBytes := make([]byte, 32)
	binary.LittleEndian.PutUint32(dailyNumberBytes, dailyNumber)

	hash := sha256.New
	info := bytes.Join([][]byte{header, dailyNumberBytes}, nil)
	hkdf := hkdf.New(hash, p.TracingKey, nil, info)

	dailyTracingKey := make([]byte, 16)
	_, err := hkdf.Read(dailyTracingKey)
	if err != nil {
		return DailyTracingKey{}, fmt.Errorf("deriving daily tracing key: %w", err)
	}

	return DailyTracingKey{Key: dailyTracingKey}, nil
}

type DailyTracingKey struct {
	Key  []byte
	hash hash.Hash
}

func (k DailyTracingKey) ProximityIdentifier(timeIntervalNumber uint8) []byte {
	if k.hash == nil {
		k.hash = hmac.New(sha256.New, k.Key)
	}

	header := []byte("CT-RPI")
	timeIntervalNumberBytes := byte(timeIntervalNumber)
	data := append(header, timeIntervalNumberBytes)

	k.hash.Reset()
	k.hash.Write(data)
	proximityIdentifier := k.hash.Sum(nil)

	return proximityIdentifier[:16]
}
