package token

import (
	"apotscan/logger"
	"apotscan/types"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ARWEAVE = "arweave"
	IPFS    = "ipfs"
	UNKNOWN = "unknown"
)

func getType(uri string) string {
	if strings.Contains(uri, "IPFS/") {
		return IPFS
	} else if strings.Contains(uri, "arweave.net/") {
		return ARWEAVE
	} else {
		return UNKNOWN
	}
}

func getMetadata(c *http.Client, tokenId string, uri string, logger *logger.Logger) *TokenMetaFromURI {
	switch getType(uri) {
	case ARWEAVE:
		resp, err := c.Get(uri)
		if err != nil {
			logger.WithFields(log.Fields{
				"error":   err,
				"uri":     uri,
				"tokenId": tokenId,
				"type":    ARWEAVE,
			}).Error("can not get uri content")
			return nil
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.WithFields(log.Fields{
				"error":   err,
				"uri":     uri,
				"tokenId": tokenId,
				"type":    ARWEAVE,
			}).Error("can not read resp body")
			return nil
		}
		var metadata TokenMetaFromURI
		if err = json.Unmarshal(data, &metadata); err != nil {
			logger.WithFields(log.Fields{
				"error":   err,
				"uri":     uri,
				"tokenId": tokenId,
			}).Error("can not unmarshal ")
			return nil
		}
		return &metadata
	case IPFS:
	case UNKNOWN:
	}
	return nil
}

type TokenMetaFromURI struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	SellerFeeBasisPoints int64       `json:"seller_fee_basis_points"`
	Description          string      `json:"description"`
	Image                string      `json:"image"`
	ExternalUrl          string      `json:"external_url"`
	AnimationUrl         string      `json:"animation_url"`
	Attributes           types.Value `json:"attributes"`
	Properties           types.Value `json:"properties"`
}
