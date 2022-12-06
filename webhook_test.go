package convoy_go

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Webhook_Verifier(t *testing.T) {
	tests := map[string]struct {
		data          *ConfigOpts
		expectedError error
	}{
		"invalid_signature": {
			data: &ConfigOpts{
				SigHeader: "",
				Payload:   []byte("test payload"),
				Secret:    "random_secret",
			},
			expectedError: ErrInvalidSignatureHeader,
		},
		"should_verify_simple_hex_signature": {
			data: &ConfigOpts{
				SigHeader: "666060cbe1348bbc7ec98f4e93dda8eaaf11bbf283d6a2dd56e841b2ef12fcd465c846903f709942473e1442604798186746f04848702c44a773f80672de7b21",
				Payload:   []byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`),
				Secret:    "8IX9njirDG",
				Hash:      "SHA512",
				Encoding:  "hex",
			},
			expectedError: nil,
		},
		"should_verify_simple_base64_signature": {
			data: &ConfigOpts{
				SigHeader: "ZmBgy+E0i7x+yY9Ok92o6q8Ru/KD1qLdVuhBsu8S/NRlyEaQP3CZQkc+FEJgR5gYZ0bwSEhwLESnc/gGct57IQ==",
				Payload:   []byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`),
				Secret:    "8IX9njirDG",
				Hash:      "SHA512",
				Encoding:  "base64",
			},
			expectedError: nil,
		},
		"invalid_signature_header": {
			data: &ConfigOpts{
				SigHeader: "d33C9sJXVO4CnE1hisHHQzUf0inr5KWJH7T8+zvgATTWEgAq5vErZR/xihDXqtok5ubv77xGP/RE++NphZnWLg==",
				Payload:   []byte(`{"email":"test@gmail.com","first_name":"test","last_name":"test"}`),
				Secret:    "8IX9njirDG",
				Hash:      "SHA512",
				Encoding:  "base64",
			},
			expectedError: ErrInvalidSignature,
		},
		"should_verify_advanced_hex_signature": {
			data: &ConfigOpts{
				SigHeader: "t=2048976161,v1=c6c39e1bd410fc1dc4db90e97039f006d088c950a275296767595d330195088f,v1=6594ee0713f1cc1f54c3f713d06a60718cd10949c7684412f159034d49621e07",
				Payload:   []byte(`{"email":"test@gmail.com"}`),
				Secret:    "Convoy",
				Hash:      "SHA256",
				Encoding:  "hex",
			},
			expectedError: nil,
		},
		"should_verify_advanced_base64_signature": {
			data: &ConfigOpts{
				SigHeader: "t=2048976161,v1=afdb90313acfa15a3fc425755ae651a204947710315bb2a90bccaa87fce88998,v1=fLBDCBUiX5iIs0L5zfNq45h23EkX1HAMpFF+2lHrnes=",
				Payload:   []byte(`{"email":"test@gmail.com"}`),
				Secret:    "8IX9njirDG",
				Hash:      "SHA256",
				Encoding:  "base64",
			},
		},
		"invalid_timestamp_header": {
			data: &ConfigOpts{
				SigHeader: "t=2202-1-1,v1=U5yBiZmFYHiom0A5hEnfLPCoQzndno4ocR45W/zkO+w=",
				Payload:   []byte(`{"email":"test@gmail.com"}`),
				Secret:    "8IX9njirDG",
				Hash:      "SHA256",
				Encoding:  "base64",
			},
			expectedError: ErrInvalidHeader,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := NewWebhook(tc.data)

			err := w.Verify()

			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
