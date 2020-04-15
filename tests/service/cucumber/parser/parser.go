package parser

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog/gherkin"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
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

func (p *Parser) ParseEnvelopes(scenario string, table *gherkin.DataTable) ([]*tx.Envelope, error) {
	var envelopes []*tx.Envelope
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		e, err := p.ParseTxRequest(scenario, headers, row)
		if err != nil {
			return nil, err
		}
		envelopes = append(envelopes, e)
	}
	return envelopes, nil
}

func (p *Parser) ParseTxRequest(scenario string, headers, row *gherkin.TableRow) (*tx.Envelope, error) {
	envelope := tx.NewEnvelope()
	gherkinRequest := make(map[string]interface{})

	for i, cell := range row.Cells {
		header := headers.Cells[i].Value

		if cell.Value == "" {
			continue
		}

		// Retrieves alias (first from scenario local namespace then if not found from global namespace)
		value, ok := p.Aliases.Get(scenario, cell.Value)
		if !ok {
			value, ok = p.Aliases.Get("global", cell.Value)
			if !ok {
				value = cell.Value
			}
		}
		switch {
		case header == "gasPrice", header == "value", header == "chainID":
			v, ok := new(big.Int).SetString(value, 10)
			if !ok {
				return nil, errors.DataError("%s invalid big int - got %s", header, value)
			}
			gherkinRequest[header] = v
		case header == "nonce", header == "gas":
			v, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, err
			}
			gherkinRequest[header] = v
		case strings.Contains(header, "."):
			err := parseMappingStringString(gherkinRequest, header, value)
			if err != nil {
				return nil, err
			}
		case header == "args", header == "privateFor":
			gherkinRequest[header] = strings.Split(value, ",")
		case header == "method":
			gherkinRequest[header] = tx.MethodMap[value]
		case header == "to", header == "from":
			gherkinRequest[header] = ethcommon.HexToAddress(value)
		case header == tenantIDHeader:
			auth, err := p.JWTGenerator.GenerateAccessTokenWithTenantID(value, 24*time.Hour)
			if err != nil {
				return nil, err
			}
			_ = parseMappingStringString(gherkinRequest, "headers.Authorization", "Bearer "+auth)
		default:
			gherkinRequest[header] = value
		}
	}

	dec, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: envelope})
	err := dec.Decode(gherkinRequest)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

func parseMappingStringString(gherkinRequest map[string]interface{}, header, value string) error {
	keyValue := strings.Split(header, ".")
	if len(keyValue) != 2 {
		return errors.DataError("invalid header")
	}
	if gherkinRequest[keyValue[0]] == nil {
		gherkinRequest[keyValue[0]] = make(map[string]string)
	}
	gherkinRequest[keyValue[0]].(map[string]string)[keyValue[1]] = value
	return nil
}
