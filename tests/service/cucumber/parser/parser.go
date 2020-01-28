package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	generator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt/generator"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

const tenantIDHeader = "tenantid"

type artifact struct {
	Abi              json.RawMessage
	Bytecode         string
	DeployedBytecode string
}

type Parser struct {
	Aliases      *AliasRegistry
	JWTGenerator *generator.JWTGenerator
}

type ContractSpec struct {
	Contract *abi.Contract
	JWTToken string
}

func New() *Parser {
	return &Parser{
		Aliases: NewAliasRegistry(),
	}
}

func (p *Parser) ParseContracts(scenario string, table *gherkin.DataTable) ([]*ContractSpec, error) {
	var contractSpecs []*ContractSpec
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		contractSpec := &ContractSpec{Contract: &abi.Contract{}}
		err := p.ParseContract(scenario, headers, row, contractSpec)
		if err != nil {
			return nil, err
		}
		contractSpecs = append(contractSpecs, contractSpec)
	}
	return contractSpecs, nil
}

func (p *Parser) ParseContract(scenario string, headers, row *gherkin.TableRow, contractSpec *ContractSpec) error {
	for i, cell := range row.Cells {
		err := p.ParseContractCell(headers.Cells[i].Value, cell.Value, contractSpec)
		if err != nil {
			return err
		}
	}
	return nil
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

func (p *Parser) ParseContractCell(header, cell string, contractSpec *ContractSpec) error {
	switch header {
	case "artifacts":
		raw, err := p.openArtifact(cell)
		if err != nil {
			return err
		}

		var a artifact
		err = json.Unmarshal(raw, &a)
		if err != nil {
			return err
		}

		// Abi is a UTF-8 encoded string. Therefore, we can make the straightforward transition
		contractSpec.Contract.Abi = a.Abi
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.Bytecode = hexutil.MustDecode(a.Bytecode)
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.DeployedBytecode = hexutil.MustDecode(a.DeployedBytecode)
	case "name":
		if contractSpec.Contract.Id == nil {
			contractSpec.Contract.Id = &abi.ContractId{}
		}
		contractSpec.Contract.Id.Name = cell
	case "tag":
		if contractSpec.Contract.Id == nil {
			contractSpec.Contract.Id = &abi.ContractId{}
		}
		contractSpec.Contract.Id.Tag = cell
	case tenantIDHeader:
		var err error
		contractSpec.JWTToken, err = p.JWTGenerator.GenerateAccessTokenWithTenantID(cell, 24*time.Hour)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func (p *Parser) ParseTxCell(header, cell string, tx *ethereum.Transaction) error {
	switch header {
	case "raw":
		tx.Raw = ethereum.HexToData(cell)
	case "hash":
		tx.Hash = ethereum.HexToHash(cell)
	case "to":
		GetInitTxData(tx).To = ethereum.HexToAccount(cell)
	case "gas":
		gas, err := strconv.ParseUint(cell, 10, 32)
		if err != nil {
			return err
		}
		GetInitTxData(tx).Gas = gas
	case "gasPrice":
		gasPrice, err := strconv.Atoi(cell)
		if err != nil {
			return err
		}
		GetInitTxData(tx).GasPrice = ethereum.IntToQuantity(int64(gasPrice))
	case "value":
		value, err := strconv.Atoi(cell)
		if err != nil {
			return err
		}
		GetInitTxData(tx).Value = ethereum.IntToQuantity(int64(value))
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

func (p *Parser) ParseChainCell(header, cell string, chn *chain.Chain) error {
	switch header {
	case "chainID":
		// Retrieve chain id
		raw, err := strconv.ParseInt(cell, 10, 64)
		if err != nil {
			return err
		}
		chn.ChainId = big.NewInt(raw).Bytes()
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
		e.From = ethereum.HexToAccount(cell)
	case strings.HasPrefix(header, "chain."):
		err := p.ParseChainCell(
			strings.TrimPrefix(header, "chain."),
			cell,
			GetInitChain(e),
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

func GetInitChain(e *envelope.Envelope) *chain.Chain {
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

func (p *Parser) openArtifact(fileName string) ([]byte, error) {
	// Loop over all cucumber folders to possibly find file
	// <cucumber_folder>/artifacts/<fileName>
	for _, v := range viper.GetStringSlice("cucumber.paths") {
		f, err := os.Open(path.Join(v, "artifacts", fileName))
		if err != nil {
			continue
		}

		bytes, readErr := ioutil.ReadAll(f)

		err = f.Close()
		if err != nil {
			log.Error(err)
		}

		return bytes, readErr
	}
	return nil, os.ErrNotExist
}
