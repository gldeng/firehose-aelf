package aelf

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
)

var (
	ZeroHash = Hash{Value: []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}}
)

func (h *Hash) ToHex() string {
	return hex.EncodeToString(h.Value)
}

func (a *Address) ToBase58() string {
	return base58EncodeWithChecksum(a.Value)
}

// Function to calculate double SHA-256 hash and return the first 4 bytes as checksum
func checksum(data []byte) []byte {
	firstHash := sha256.Sum256(data)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:4]
}

// Function to encode data using Base58 with checksum
func base58EncodeWithChecksum(data []byte) string {
	// Calculate the checksum
	checksum := checksum(data)

	// Append the checksum to the data
	dataWithChecksum := append(data, checksum...)

	// Encode using Base58
	encoded := base58.Encode(dataWithChecksum)
	return encoded
}
