package convoy_go

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
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
	DefaultSigHeader              = "X-Convoy-Signature"
)

type signedHeader struct {
	timestamp  time.Time
	signatures [][]byte
	isAdvanced bool
}

type Webhook struct {
	opts *WebhookOpts
}

type WebhookOpts struct {
	SigHeader string
	Secret    string
	Encoding  EncodingType
	Hash      string
	Tolerance time.Duration
}

type EncodingType string

const (
	Base64Encoding EncodingType = "base64"
	HexEncoding    EncodingType = "hex"
)

func NewWebhook(opts *WebhookOpts) *Webhook {

	if isStringEmpty(opts.Hash) {
		opts.Hash = DefaultHash
	}

	if isStringEmpty(string(opts.Encoding)) {
		opts.Encoding = DefaultEncoding
	}

	if opts.Tolerance == 0 {
		opts.Tolerance = DefaultTolerance
	}

	if isStringEmpty(opts.SigHeader) {
		opts.SigHeader = DefaultSigHeader
	}

	return &Webhook{opts}
}

func (w *Webhook) VerifyRequest(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	header := r.Header.Get(w.opts.SigHeader)
	if isStringEmpty(header) {
		return ErrInvalidHeader
	}

	return w.verify(body, header)
}

func (w *Webhook) VerifyPayload(b []byte, header string) error {
	return w.verify(b, header)
}

func (w *Webhook) verify(body []byte, header string) error {
	sh, err := w.parseSignatureHeader(header)
	if err != nil {
		return err
	}

	expectedSignature, err := w.generateSignature(sh, body)
	if err != nil {
		return err
	}

	for _, sig := range sh.signatures {
		// Check all signatures for a match
		if hmac.Equal(expectedSignature, sig) {
			return nil
		}
	}

	return ErrInvalidSignature
}

func (w *Webhook) parseSignatureHeader(header string) (*signedHeader, error) {
	var err error
	sh := &signedHeader{}

	if isStringEmpty(header) {
		return sh, ErrInvalidSignatureHeader
	}

	pairs := strings.Split(header, ",")
	if len(pairs) > 1 {
		sh, err = w.decodeAdvanced(sh, pairs)
		if err != nil {
			return sh, err
		}
	} else {
		sh, err = w.decodeSimple(sh, header)
		if err != nil {
			return sh, err
		}
	}

	if len(sh.signatures) == 0 {
		return sh, ErrInvalidSignature
	}

	return sh, nil
}

func (w *Webhook) generateSignature(sh *signedHeader, body []byte) ([]byte, error) {
	fn, err := w.getHashFunction(w.opts.Hash)
	if err != nil {
		return nil, err
	}

	h := hmac.New(fn, []byte(w.opts.Secret))

	if sh.isAdvanced {
		h.Write([]byte(fmt.Sprintf("%d", sh.timestamp.Unix())))
		h.Write([]byte(","))
	}

	h.Write(body)
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
			continue
		}

		if strings.Contains(item, "v") {
			sig, err := w.decodeString(parts[1])
			if err != nil {
				continue
			}

			sh.signatures = append(sh.signatures, sig)
		}
	}

	expiredTimestamp := time.Since(sh.timestamp) > w.opts.Tolerance
	if expiredTimestamp {
		return nil, ErrTimestampExpired
	}

	sh.isAdvanced = true
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
	switch w.opts.Encoding {
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
