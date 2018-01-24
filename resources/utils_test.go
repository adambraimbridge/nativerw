package resources

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObtainTxID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesnt/matter", nil)
	req.Header.Add("X-Request-Id", "tid_blahblah")
	txid := obtainTxID(req)
	assert.Equal(t, "tid_blahblah", txid)
}

func TestObtainTxIDGeneratesANewOneIfNoneAvailable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesnt/matter", nil)
	txid := obtainTxID(req)
	assert.Contains(t, txid, "tid_")
}
