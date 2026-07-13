package convoy_go

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// signature-vectors.json is generated from the server signing code
// (convoy/pkg/signature) and vendored here so this SDK verifies against the same
// canonical set as every other Convoy SDK. Regenerate upstream with
// CONVOY_WRITE_VECTORS=1 go test ./pkg/signature -run TestGenerateSignatureVectors
type signatureVector struct {
	Name      string `json:"name"`
	Advanced  bool   `json:"advanced"`
	Hash      string `json:"hash"`
	Encoding  string `json:"encoding"`
	Secret    string `json:"secret"`
	Payload   string `json:"payload"`
	Header    string `json:"header"`
	Tolerance int64  `json:"tolerance"`
	Valid     bool   `json:"valid"`
}

func Test_Webhook_SharedVectors(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "signature-vectors.json"))
	require.NoError(t, err)

	var vectors []signatureVector
	require.NoError(t, json.Unmarshal(raw, &vectors))
	require.NotEmpty(t, vectors)

	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			w := NewWebhook(&WebhookOpts{
				Secret:    v.Secret,
				Hash:      v.Hash,
				Encoding:  EncodingType(v.Encoding),
				Tolerance: time.Duration(v.Tolerance) * time.Second,
			})

			err := w.VerifyPayload([]byte(v.Payload), v.Header)
			if v.Valid {
				require.NoError(t, err, v.Name)
			} else {
				require.Error(t, err, v.Name)
			}
		})
	}
}

// Test_Webhook_AdvancedSignatureKeyMustStartWithV guards the parse tightening:
// only keys starting with "v" carry a signature. A real signature smuggled
// under any other key must be rejected, not accepted because the key happens to
// contain "v".
func Test_Webhook_AdvancedSignatureKeyMustStartWithV(t *testing.T) {
	const (
		secret  = "convoy-webhook-secret"
		payload = `{"m":1}`
	)
	ts := time.Now().Unix()

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d,%s", ts, payload)))
	sig := hex.EncodeToString(mac.Sum(nil))

	w := NewWebhook(&WebhookOpts{Secret: secret, Hash: "SHA256", Encoding: "hex"})

	// Control: the same signature under a v-prefixed key verifies.
	require.NoError(t, w.VerifyPayload([]byte(payload), fmt.Sprintf("t=%d,v1=%s", ts, sig)))

	// A non-v key (which still contains "v") must not be treated as a signature.
	require.Error(t, w.VerifyPayload([]byte(payload), fmt.Sprintf("t=%d,nav=%s", ts, sig)))
}
