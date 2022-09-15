package token

import "strings"

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

func getMetadata(uri string) TokenMetaFromURI {
	switch getType(uri) {
	case ARWEAVE:
	case IPFS:
	case UNKNOWN:

	}
}

type TokenMetaFromURI struct {
}
