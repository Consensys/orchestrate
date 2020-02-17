package tx

import (
	"math/big"
	"regexp"

	error1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var MethodMap = map[string]Method{
	"ETH_SENDRAWTRANSACTION":        Method_ETH_SENDRAWTRANSACTION,
	"ETH_SENDPRIVATETRANSACTION":    Method_ETH_SENDPRIVATETRANSACTION,
	"ETH_SENDRAWPRIVATETRANSACTION": Method_ETH_SENDRAWPRIVATETRANSACTION,
	"EEA_SENDPRIVATETRANSACTION":    Method_EEA_SENDPRIVATETRANSACTION,
}

func (m *TxEnvelope) Envelope() (*Envelope, error) {
	var b *Envelope

	var err error
	switch x := m.Msg.(type) {
	case *TxEnvelope_TxRequest:
		b, err = x.TxRequest.Envelope()
	case *TxEnvelope_TxResponse:
		b, err = x.TxResponse.Envelope()
	default:
		return nil, errors.DataError("invalid tx envelope")
	}
	if err != nil {
		return nil, err
	}

	b.InternalLabels = m.GetInternalLabels()
	if err := b.internalToFields(); err != nil {
		return nil, err
	}

	return b, nil
}

func (m *TxRequest) Envelope() (*Envelope, error) {
	envelope := &Envelope{
		ID:            m.GetId(),
		Headers:       m.GetHeaders(),
		ContextLabels: m.GetContextLabels(),
		Method:        m.Method,
		Tx: Tx{
			Data: m.GetParams().GetData(),
			Raw:  m.GetParams().GetRaw(),
		},
		Chain: Chain{
			ChainName: m.GetChain(),
		},
		Contract: Contract{
			MethodSignature: m.GetParams().GetMethodSignature(),
			Args:            m.GetParams().GetArgs(),
		},
		Private: Private{
			PrivateFor:     m.GetParams().GetPrivateFor(),
			PrivateFrom:    m.GetParams().GetPrivateFrom(),
			PrivateTxType:  m.GetParams().GetPrivateTxType(),
			PrivacyGroupID: m.GetParams().GetPrivacyGroupId(),
		},
		Errors:         make([]*error1.Error, 0),
		InternalLabels: make(map[string]string),
	}

	if errs := envelope.loadPtrFields(m.GetParams().GetGas(), m.GetParams().GetNonce(), m.GetParams().GetGasPrice(), m.GetParams().GetValue(), m.GetParams().GetFrom(), m.GetParams().GetTo()); len(errs) > 0 {
		return nil, errors.DataError("%v", errs)
	}

	contractName, contractTag, err := m.GetParams().GetParsedContract()
	if err != nil {
		return nil, errors.DataError("%v", err)
	}
	_ = envelope.SetContractName(contractName).SetContractTag(contractTag)

	if err := envelope.Validate(); err != nil {
		return nil, errors.DataError("%v", err)
	}

	return envelope, nil
}

func (m *TxResponse) Envelope() (*Envelope, error) {
	envelope := &Envelope{
		ID:            m.GetId(),
		Headers:       m.GetHeaders(),
		ContextLabels: m.GetContextLabels(),
		Tx: Tx{
			Data: m.GetTransaction().GetData(),
			Raw:  m.GetTransaction().GetRaw(),
		},
		Receipt:        m.GetReceipt(),
		Errors:         m.GetErrors(),
		InternalLabels: make(map[string]string),
	}

	if errs := envelope.loadPtrFields(m.GetTransaction().GetGas(), m.GetTransaction().GetNonce(), m.GetTransaction().GetGasPrice(), m.GetTransaction().GetValue(), m.GetTransaction().GetFrom(), m.GetTransaction().GetTo()); len(errs) > 0 {
		return nil, errors.DataError("%v", errs)
	}

	if err := envelope.Validate(); err != nil {
		return nil, errors.DataError("%v", err)
	}

	return envelope, nil
}

func (p *Params) GetParsedContract() (contractName, contractTag string, err error) {
	if p.GetContract() == "" {
		return "", "", nil
	}

	re := regexp.MustCompile(`^(.*)\[(.*)\]$`)
	t := re.FindStringSubmatch(p.GetContract())

	if len(t) == 3 && t[0] == p.GetContract() {
		return t[1], t[2], nil
	}

	return p.GetContract(), "", nil
}

func (m *TxEnvelope) GetChainID() string {
	return m.InternalLabels["chainID"]
}

func (m *TxEnvelope) SetChainID(chainID *big.Int) *TxEnvelope {
	m.InternalLabels["chainID"] = chainID.String()
	return m
}

func (m *TxEnvelope) SetChainUUID(chainUUID string) *TxEnvelope {
	m.InternalLabels["chainUUID"] = chainUUID
	return m
}

func (m *TxEnvelope) GetChainUUID() string {
	return m.InternalLabels["chainUUID"]
}

func (m *TxEnvelope) GetTxHash() string {
	return m.InternalLabels["txHash"]
}

func (m *TxEnvelope) TxHash() ethcommon.Hash {
	return ethcommon.HexToHash(m.InternalLabels["txHash"])
}

func (m *TxEnvelope) SetTxHash(txHash string) *TxEnvelope {
	m.InternalLabels["txHash"] = txHash
	return m
}

func (m *TxEnvelope) GetID() string {
	switch x := m.Msg.(type) {
	case *TxEnvelope_TxRequest:
		return x.TxRequest.GetId()
	case *TxEnvelope_TxResponse:
		return x.TxResponse.GetId()
	default:
		return ""
	}
}

func (m *TxEnvelope) MustGetTxRequest() *TxRequest {
	return m.Msg.(*TxEnvelope_TxRequest).TxRequest
}

func (m *TxEnvelope) MustGetTxResponse() *TxResponse {
	return m.Msg.(*TxEnvelope_TxResponse).TxResponse
}
