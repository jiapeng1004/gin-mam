package transcode

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignCallbackBody(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyCallbackSign(body []byte, secret, signHeader string) bool {
	if secret == "" || signHeader == "" {
		return false
	}
	expected := SignCallbackBody(body, secret)
	return hmac.Equal([]byte(expected), []byte(signHeader))
}