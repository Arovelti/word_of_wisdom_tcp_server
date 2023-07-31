package hashcash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

const (
	maxAttempts            = 1000000
	maxTargetBits          = 64
	defaultBufferBytesSize = 16
	defaultNonceSize       = 8
)

var (
	ErrLowComplexity     = errors.New("the complexity must be more then 0")
	ErrHighComplexity    = errors.New("too large (should be at most 64) complexity, please, make it lower")
	ErrEncodeTargetValue = errors.New("failed to encode target value")
	ErrBytesBufferSize   = errors.New("invalid data bytes buffer size")
	ErrNonceSize         = errors.New("wrong nonce result size")
	ErrWritingDataToHash = errors.New("error writing data to hash")
	ErrCompare           = errors.New("error while comparing data and hash")
	ErrGenerateNonce     = errors.New("failed to generate nonce")
)

type ProofOfWork struct {
	l          *logrus.Logger
	complexity uint8
}

func New(l *logrus.Logger, difficulty uint8) (*ProofOfWork, error) {
	switch {
	case difficulty == 0:
		return nil, ErrLowComplexity
	case difficulty > maxTargetBits:
		return nil, ErrHighComplexity
	default:
		return &ProofOfWork{l: l, complexity: difficulty}, nil
	}
}

func (p *ProofOfWork) GenerateBytesBuffer(difficulty uint8) ([]byte, error) {
	target := uint64(math.Pow(2, float64(64-difficulty)))

	b := make([]byte, defaultBufferBytesSize)
	n := binary.PutUvarint(b[:8], target)
	if n <= 0 {
		return nil, ErrEncodeTargetValue
	}

	_, err := rand.Read(b[8:])
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *ProofOfWork) Prove(data, nonce []byte) error {
	switch {
	case len(data) != defaultBufferBytesSize:
		return ErrBytesBufferSize
	case len(nonce) != defaultNonceSize:
		return ErrNonceSize
	default:
		if err := compare(data, nonce); err != nil {
			return err
		}
	}

	return nil
}

func (p *ProofOfWork) Work(data []byte) ([]byte, error) {
	nonce := make([]byte, defaultNonceSize)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if _, err := rand.Read(nonce); err != nil {
			return nil, ErrGenerateNonce
		}

		if err := compare(data, nonce); err == nil {
			return nonce, nil
		}
	}

	return nil, fmt.Errorf("could not find a valid nonce after %d attempts", maxAttempts)
}

func compare(data, nonce []byte) error {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return ErrWritingDataToHash
	}
	if _, err := h.Write(nonce); err != nil {
		return ErrWritingDataToHash
	}

	hash := h.Sum(nil)
	if !bytes.Equal(data, hash) {
		return ErrCompare
	}

	return nil
}
