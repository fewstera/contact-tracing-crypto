package tracing

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/minio/sha256-simd"

	"golang.org/x/crypto/hkdf"
)

// Person represents a single user of contact tracing.
type Person struct {
	TracingKey []byte
}

// GeneratePerson generates a new Person with a random TracingKey
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

//DailyTracingKey returns the daily tracing key for the given dailyNumber
func (p Person) DailyTracingKey(dailyNumber uint32) (DailyTracingKey, error) {
	header := []byte("CT-DTK")
	dailyNumberBytes := make([]byte, 32)
	binary.LittleEndian.PutUint32(dailyNumberBytes, dailyNumber)

	hash := sha256.New
	info := bytes.Join([][]byte{header, dailyNumberBytes}, nil)
	hkdf := hkdf.New(hash, p.TracingKey, nil, info)

	dailyTracingKey := make(DailyTracingKey, 16)
	_, err := hkdf.Read(dailyTracingKey)
	if err != nil {
		return nil, fmt.Errorf("deriving daily tracing key: %w", err)
	}

	return dailyTracingKey, nil
}

// DailyTracingKey is the daily tracing key
type DailyTracingKey []byte

// ProximityIdentifier returns a proximity identifier for a given time interval number
func (key DailyTracingKey) ProximityIdentifier(timeIntervalNumber uint8) []byte {
	header := []byte("CT-RPI")
	timeIntervalNumberBytes := byte(timeIntervalNumber)

	data := append(header, timeIntervalNumberBytes)
	h := hmac.New(sha256.New, key)
	h.Write(data)

	proximityIdentifier := h.Sum(nil)

	return proximityIdentifier[:16]
}

// AllProximityIdentifiers returns all proximity identifiers for a given daily tracing key
func (key DailyTracingKey) AllProximityIdentifiers() [][]byte {
	h := hmac.New(sha256.New, key)
	header := []byte("CT-RPI")
	proximityIdentifiers := [][]byte{}

	for i := 0; i < 144; i++ {
		timeIntervalNumberBytes := []byte{uint8(i)}
		h.Write(bytes.Join([][]byte{header, timeIntervalNumberBytes}, nil))

		proximityIdentifiers = append(proximityIdentifiers, h.Sum(nil)[:16])
		h.Reset()
	}

	return proximityIdentifiers
}
