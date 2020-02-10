package parser

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	generator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt/generator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

const tenantIDHeader = "tenantid"

type Parser struct {
	Aliases      *AliasRegistry
	JWTGenerator *generator.JWTGenerator
}

func New() *Parser {
	return &Parser{
		Aliases: NewAliasRegistry(),
	}
}

func (p *Parser) ParseEnvelopes(scenario string, table *gherkin.DataTable) ([]*envelope.Envelope, error) {
	var envelopes []*envelope.Envelope
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		e := &envelope.Envelope{}
		err := p.ParseEnvelope(scenario, headers, row, e)
		if err != nil {
			return nil, err
		}
		envelopes = append(envelopes, e)
	}
	return envelopes, nil
}

func (p *Parser) ParseEnvelope(scenario string, headers, row *gherkin.TableRow, e *envelope.Envelope) error {
	for i, cell := range row.Cells {
		header := headers.Cells[i].Value

		// Retrieves alias (first from scenario local namespace then if not found from global namespace)
		value, ok := p.Aliases.Get(scenario, cell.Value)
		if !ok {
			value, ok = p.Aliases.Get("global", cell.Value)
			if !ok {
				value = cell.Value
			}
		}

		err := p.ParseEnvelopeCell(header, value, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) ParseTxCell(header, cell string, tx *ethereum.Transaction) error {
	switch header {
	case "raw":
		tx.Raw = cell
	case "hash":
		tx.Hash = cell
	case "to":
		GetInitTxData(tx).To = cell
	case "gas":
		gas, err := strconv.ParseUint(cell, 10, 32)
		if err != nil {
			return err
		}
		GetInitTxData(tx).Gas = gas
	case "gasPrice":
		gasPrice, ok := (new(big.Int)).SetString(cell, 10)
		if !ok {
			return fmt.Errorf("invalid gas price")
		}
		GetInitTxData(tx).GasPrice = gasPrice.String()
	case "value":
		value, ok := (new(big.Int)).SetString(cell, 10)
		if !ok {
			return fmt.Errorf("invalid value")
		}
		GetInitTxData(tx).Value = value.String()
	case "nonce":
		nonce, err := strconv.ParseUint(cell, 10, 64)
		if err != nil {
			return err
		}
		GetInitTxData(tx).SetNonce(nonce)
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParseMethodCell(header, cell string, method *abi.Method) error {
	switch header {
	case "sig":
		method.Signature = cell
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParseTxChainCell(header, cell string, chn *chain.Chain) error {
	switch header {
	case "chainID":
		// Retrieve chain id
		chainID, ok := (new(big.Int)).SetString(cell, 10)
		if !ok {
			return fmt.Errorf("invalid chainID")
		}
		chn.ChainId = chainID.String()
	case "name":
		chn.Name = cell
	case "uuid":
		chn.Uuid = cell
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParsePrivateArgCell(header, cell string, private *args.Private) error {
	switch header {
	case "privateFrom":
		private.PrivateFrom = cell
	case "privateFor":
		private.PrivateFor = strings.Split(cell, ",")
	case "privateTxType":
		private.PrivateTxType = cell
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParseCallCell(header, cell string, call *args.Call) error {
	switch {
	case header == "args":
		call.Args = strings.Split(cell, ",")
	case strings.HasPrefix(header, "contract."):
		if call.Contract == nil {
			call.Contract = &abi.Contract{}
		}
		contractSpec := &ContractSpec{
			Contract: call.Contract,
		}

		err := p.ParseContractCell(
			strings.TrimPrefix(header, "contract."),
			cell,
			contractSpec,
		)
		if err != nil {
			return err
		}
	case strings.HasPrefix(header, "method."):
		if call.Method == nil {
			call.Method = &abi.Method{}
		}
		err := p.ParseMethodCell(
			strings.TrimPrefix(header, "method."),
			cell,
			call.Method,
		)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParseEnvelopeCell(header, cell string, e *envelope.Envelope) error {
	switch {
	case header == "from":
		e.From = cell
	case strings.HasPrefix(header, "chain."):
		err := p.ParseTxChainCell(
			strings.TrimPrefix(header, "chain."),
			cell,
			GetInitTxChain(e),
		)
		if err != nil {
			return err
		}
	case strings.HasPrefix(header, "tx."):
		err := p.ParseTxCell(
			strings.TrimPrefix(header, "tx."),
			cell,
			GetInitTx(e),
		)
		if err != nil {
			return err
		}
	case header == "args" || strings.HasPrefix(header, "contract.") || strings.HasPrefix(header, "method."):
		err := p.ParseCallCell(
			header,
			cell,
			GetInitCall(GetInitArgs(e)),
		)
		if err != nil {
			return err
		}
	case strings.HasPrefix(header, "private"):
		err := p.ParsePrivateArgCell(
			header,
			cell,
			GetInitPrivate(GetInitArgs(e)),
		)
		if err != nil {
			return err
		}
	case header == "protocol":
		protocol, err := strconv.Atoi(cell)
		if err != nil {
			return err
		}
		e.Protocol = &chain.Protocol{
			Type: chain.ProtocolType(
				int64(protocol),
			),
		}
	case header == tenantIDHeader:
		auth, err := p.JWTGenerator.GenerateAccessTokenWithTenantID(cell, 24*time.Hour)
		if err != nil {
			return err
		}
		// Add authorization header
		e.SetMetadataValue("Authorization", "Bearer "+auth)
	default:
		return fmt.Errorf("got unknown header %q", header)
	}
	return nil
}

func GetInitTxChain(e *envelope.Envelope) *chain.Chain {
	if e.Chain == nil {
		e.Chain = &chain.Chain{}
	}
	return e.Chain
}

func GetInitArgs(e *envelope.Envelope) *envelope.Args {
	if e.Args == nil {
		e.Args = &envelope.Args{}
	}
	return e.Args
}

func GetInitCall(a *envelope.Args) *args.Call {
	if a.Call == nil {
		a.Call = &args.Call{}
	}
	return a.Call
}

func GetInitPrivate(a *envelope.Args) *args.Private {
	if a.Private == nil {
		a.Private = &args.Private{}
	}
	return a.Private
}

func GetInitTx(e *envelope.Envelope) *ethereum.Transaction {
	if e.Tx == nil {
		e.Tx = &ethereum.Transaction{}
	}
	return e.Tx
}

func GetInitTxData(tx *ethereum.Transaction) *ethereum.TxData {
	if tx.TxData == nil {
		tx.TxData = &ethereum.TxData{}
	}
	return tx.TxData
}
