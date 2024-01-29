package convoy_go

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Webhook_VerifyRequest(t *testing.T) {
	tests := map[string]struct {
		opts          *WebhookOpts
		req           func() *http.Request
		expectedError error
	}{
		"invalid_header": {
			opts: &WebhookOpts{
				SigHeader: "",
				Secret:    "random_secret",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(``))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				return req
			},
			expectedError: ErrInvalidHeader,
		},
		"should_verify_simple_hex_signature": {
			opts: &WebhookOpts{
				Secret:   "8IX9njirDG",
				Hash:     "SHA512",
				Encoding: "hex",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"666060cbe1348bbc7ec98f4e93dda8eaaf11bbf283d6a2dd56e841b2ef12fcd465c846903f709942473e1442604798186746f04848702c44a773f80672de7b21")

				return req
			},
			expectedError: nil,
		},
		"should_verify_simple_base64_signature": {
			opts: &WebhookOpts{
				Secret:   "8IX9njirDG",
				Hash:     "SHA512",
				Encoding: "base64",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"ZmBgy+E0i7x+yY9Ok92o6q8Ru/KD1qLdVuhBsu8S/NRlyEaQP3CZQkc+FEJgR5gYZ0bwSEhwLESnc/gGct57IQ==")

				return req
			},
			expectedError: nil,
		},
		"invalid_signature_header": {
			opts: &WebhookOpts{
				Secret:   "8IX9njirDG",
				Hash:     "SHA512",
				Encoding: "base64",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"d33C9sJXVO4CnE1hisHHQzUf0inr5KWJH7T8+zvgATTWEgAq5vErZR/xihDXqtok5ubv77xGP/RE++NphZnWLg==")

				return req
			},
			expectedError: ErrInvalidSignature,
		},
		"should_verify_advanced_hex_signature": {
			opts: &WebhookOpts{
				Secret:   "Convoy",
				Hash:     "SHA256",
				Encoding: "hex",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com"}`))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"t=2048976161,v1=c6c39e1bd410fc1dc4db90e97039f006d088c950a275296767595d330195088f,v1=6594ee0713f1cc1f54c3f713d06a60718cd10949c7684412f159034d49621e07")

				return req
			},
			expectedError: nil,
		},
		"should_verify_advanced_base64_signature": {
			opts: &WebhookOpts{
				Secret:   "8IX9njirDG",
				Hash:     "SHA256",
				Encoding: "base64",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com"}`))

				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"t=2048976161,v1=afdb90313acfa15a3fc425755ae651a204947710315bb2a90bccaa87fce88998,v1=fLBDCBUiX5iIs0L5zfNq45h23EkX1HAMpFF+2lHrnes=")

				return req
			},
		},
		"invalid_timestamp_header": {
			opts: &WebhookOpts{
				SigHeader: "t=2202-1-1,v1=U5yBiZmFYHiom0A5hEnfLPCoQzndno4ocR45W/zkO+w=",
				Secret:    "8IX9njirDG",
				Hash:      "SHA256",
				Encoding:  "base64",
			},
			req: func() *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"test@gmail.com"}`))
				req, err := http.NewRequest(http.MethodPost, "localhost:5005", body)
				require.NoError(t, err)

				req.Header.Add(DefaultSigHeader,
					"t=2202-1-1,v1=U5yBiZmFYHiom0A5hEnfLPCoQzndno4ocR45W/zkO+w=")

				return req
			},
			expectedError: ErrInvalidHeader,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := NewWebhook(tc.opts)

			err := w.VerifyRequest(tc.req())

			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
