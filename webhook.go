package convoy_go

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidSignatureHeader = errors.New("webhook has no signature header")
	ErrInvalidHeader          = errors.New("webhook has invalid header")
	ErrInvalidEncoding        = errors.New("invalid encoding")
	ErrInvalidSignature       = errors.New("webhook has no valid signature")
	ErrInvalidHashAlgorithm   = errors.New("invalid hash algorithm")
	ErrTimestampExpired       = errors.New("timestamp has expired")
)

var (
	DefaultTolerance              = 300 * time.Second
	DefaultEncoding  EncodingType = HexEncoding
	DefaultHash                   = "SHA256"
)

type signedHeader struct {
	timestamp  time.Time
	signatures [][]byte
}

type Webhook struct {
	Payload    []byte
	SigHeader  string
	Secret     string
	IsAdvanced bool
	Encoding   EncodingType
	Hash       string
	Tolerance  time.Duration
}

type CreateWebhook struct {
	Payload    []byte
	SigHeader  string
	Secret     string
	IsAdvanced bool
	Encoding   EncodingType
	Hash       string
	Tolerance  time.Duration
}

type EncodingType string

const (
	Base64Encoding EncodingType = "base64"
	HexEncoding    EncodingType = "hex"
)

func NewWebhook(data *CreateWebhook) *Webhook {
	w := &Webhook{
		Payload:   data.Payload,
		SigHeader: data.SigHeader,
		Secret:    data.Secret,
		Hash:      data.Hash,
		Encoding:  data.Encoding,
		Tolerance: data.Tolerance,
	}

	if w.Hash == "" {
		w.Hash = DefaultHash
	}

	if w.Encoding == "" {
		w.Encoding = DefaultEncoding
	}

	if w.Tolerance == 0 {
		w.Tolerance = DefaultTolerance
	}

	return w
}

func (w *Webhook) Verify() error {
	header, err := w.parseSignatureHeader()
	if err != nil {
		return err
	}

	expectedSignature, err := w.generateSignature()
	if err != nil {
		return err
	}

	for _, sig := range header.signatures {
		// Check all signatures for a match
		if hmac.Equal(expectedSignature, sig) {
			return nil
		}
	}

	return ErrInvalidSignature
}

func (w *Webhook) parseSignatureHeader() (*signedHeader, error) {
	var err error
	sh := &signedHeader{}

	if w.SigHeader == "" {
		return sh, ErrInvalidSignatureHeader
	}

	pairs := strings.Split(w.SigHeader, ",")
	if len(pairs) > 1 {
		sh, err = w.decodeAdvanced(sh, pairs)
		if err != nil {
			return sh, err
		}
	} else {
		sh, err = w.decodeSimple(sh, w.SigHeader)
		if err != nil {
			return sh, err
		}
	}

	if len(sh.signatures) == 0 {
		return sh, ErrInvalidSignature
	}

	return sh, nil
}

func (w *Webhook) generateSignature() ([]byte, error) {
	fn, err := w.getHashFunction(w.Hash)
	if err != nil {
		return nil, err
	}

	h := hmac.New(fn, []byte(w.Secret))
	h.Write(w.Payload)
	return h.Sum(nil), nil
}

func (w *Webhook) getHashFunction(algorithm string) (func() hash.Hash, error) {
	switch algorithm {
	case "SHA256":
		return sha256.New, nil
	case "SHA512":
		return sha512.New, nil
	}
	return nil, ErrInvalidHashAlgorithm
}

func (w *Webhook) decodeAdvanced(sh *signedHeader, pairs []string) (*signedHeader, error) {
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return sh, ErrInvalidHeader
		}

		item := parts[0]

		if item == "t" {
			timestamp, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return sh, ErrInvalidHeader
			}

			sh.timestamp = time.Unix(timestamp, 0)
		}

		if strings.Contains(item, "v") {
			sig, err := w.decodeString(parts[1])
			if err != nil {
				continue
			}

			sh.signatures = append(sh.signatures, sig)
		}

		continue
	}

	expiredTimestamp := time.Since(sh.timestamp) > w.Tolerance
	if expiredTimestamp {
		return nil, ErrTimestampExpired
	}

	return sh, nil
}

func (w *Webhook) decodeSimple(sh *signedHeader, value string) (*signedHeader, error) {
	sig, err := w.decodeString(value)
	if err != nil {
		return sh, err
	}

	sh.signatures = append(sh.signatures, sig)
	return sh, nil
}

func (w *Webhook) decodeString(value string) ([]byte, error) {
	switch w.Encoding {
	case HexEncoding:
		sig, err := hex.DecodeString(value)
		return sig, err
	case Base64Encoding:
		sig, err := base64.StdEncoding.DecodeString(value)
		return sig, err
	default:
		return nil, ErrInvalidEncoding
	}
}
