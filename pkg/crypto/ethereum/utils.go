package ethereum

import (
	"encoding/base64"

	"github.com/consensys/orchestrate/pkg/errors"
	quorumtypes "github.com/consensys/quorum/core/types"
)

func GetEncodedPrivateFrom(privateFrom string) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(privateFrom)
	if err != nil {
		return nil, errors.EncodingError("invalid base64 value for 'privateFrom'").AppendReason(err.Error())
	}

	return privateFromEncoded, nil
}

func GetEncodedPrivateRecipient(privacyGroupID string, privateFor []string) (interface{}, error) {
	var privateRecipientEncoded interface{}
	var err error
	if privacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privacyGroupID)
		if err != nil {
			return nil, errors.EncodingError("invalid base64 value for 'privacyGroupId'").AppendReason(err.Error())
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				return nil, errors.EncodingError("invalid base64 value for 'privateFor'").AppendReason(der.Error())
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	return privateRecipientEncoded, nil
}

func GetQuorumPrivateTxSigner() quorumtypes.Signer {
	return quorumtypes.QuorumPrivateTxSigner{}
}
