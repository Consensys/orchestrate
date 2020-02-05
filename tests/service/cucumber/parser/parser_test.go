package parser

import (
	"testing"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

func TestParseTxCell(t *testing.T) {
	p := &Parser{}
	tx := &ethereum.Transaction{}
	err := p.ParseTxCell("raw", "0xabcd", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting raw")
	assert.Equal(t, "0xabcd", tx.GetRaw().Hex(), "Raw should have been set")

	err = p.ParseTxCell("hash", "0xabcd000000000000000000000000000000000000000000000000000000000000", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting hash")
	assert.Equal(t, "0xabcd000000000000000000000000000000000000000000000000000000000000", tx.GetHash().Hex(), "Hash should have been set")

	err = p.ParseTxCell("to", "0xabcd000000abcd000000abcd000000abcd000000", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting to")
	assert.Equal(t, ethereum.HexToAccount("0xabcd000000abcd000000abcd000000abcd000000").Hex(), tx.GetTxData().GetTo().Hex(), "To should have been set")

	err = p.ParseTxCell("gas", "1000", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting gas")
	assert.Equal(t, uint64(1000), tx.GetTxData().GetGas(), "Gas should have been set")

	err = p.ParseTxCell("gasPrice", "1000000000", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting gas price")
	assert.Equal(t, "1000000000", tx.GetTxData().GetGasPrice().Value().String(), "GasPrice should have been set")

	err = p.ParseTxCell("nonce", "17", tx)
	assert.NoError(t, err, "ParseTxCell should not error when setting nonce")
	assert.Equal(t, uint64(17), tx.GetTxData().GetNonce(), "Nonce should have been set")

	err = p.ParseTxCell("unknown", "17", tx)
	assert.Error(t, err, "ParseTxCell should error when setting unknonw")
}

func TestParseMethodCell(t *testing.T) {
	p := &Parser{}
	mthd := &abi.Method{}
	err := p.ParseMethodCell("sig", "transfer()", mthd)
	assert.NoError(t, err, "ParseMethodCell should not error when setting signature")
	assert.Equal(t, "transfer()", mthd.GetSignature(), "Signature should have been set")

	err = p.ParseMethodCell("unknown", "17", mthd)
	assert.Error(t, err, "ParseMethodCell should error when setting unknonw")
}

func TestParseTxChainCell(t *testing.T) {
	p := &Parser{}
	chn := &chain.Chain{}
	err := p.ParseTxChainCell("chainID", "17", chn)
	assert.NoError(t, err, "ParseChainCell should not error when setting id")
	assert.Equal(t, "17", chn.GetBigChainID().String(), "UUID should have been set")

	err = p.ParseTxChainCell("unknown", "17", chn)
	assert.Error(t, err, "ParseChainCell should error when setting unknonw")
}

func TestParsePrivateArgCell(t *testing.T) {
	p := &Parser{}
	priv := &args.Private{}
	err := p.ParsePrivateArgCell("privateFrom", "foo", priv)
	assert.NoError(t, err, "ParsePrivateArgCell should not error when setting PrivateFrom")
	assert.Equal(t, "foo", priv.PrivateFrom, "PrivateFrom should have been set")

	err = p.ParsePrivateArgCell("privateFor", "foo,bar", priv)
	assert.NoError(t, err, "ParsePrivateArgCell should not error when setting privateFor")
	assert.Equal(t, []string{"foo", "bar"}, priv.PrivateFor, "PrivateFor should have been set")

	err = p.ParsePrivateArgCell("privateTxType", "test", priv)
	assert.NoError(t, err, "ParsePrivateArgCell should not error when setting privateTxType")
	assert.Equal(t, "test", priv.PrivateTxType, "PrivateTxType should have been set")

	err = p.ParsePrivateArgCell("unknown", "17", priv)
	assert.Error(t, err, "ParsePrivateArgCell should error when setting unknonw")
}

func TestParseEnvelopes(t *testing.T) {
	p := &Parser{
		Aliases: NewAliasRegistry(),
	}

	// Set a chain alias in global namespace
	p.Aliases.Set("global", "chain.primary", "888")

	// Set Contract alias in local namespaces
	p.Aliases.Set("test-1", "Contract.my-token", "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	p.Aliases.Set("test-2", "Contract.my-token", "0x77F888CC34a3E6EC4935eF27a83a48fAe548fa4d")

	headers := &gherkin.TableRow{
		Cells: []*gherkin.TableCell{
			{Value: "chain.chainID"},
			{Value: "from"},
			{Value: "tx.to"},
			{Value: "protocol"},
			{Value: "args"},
		},
	}

	row1 := &gherkin.TableRow{
		Cells: []*gherkin.TableCell{
			{Value: "17"},
			{Value: "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0"},
			{Value: "Contract.my-token"},
			{Value: "2"},
			{Value: "1,2"},
		},
	}

	row2 := &gherkin.TableRow{
		Cells: []*gherkin.TableCell{
			{Value: "chain.primary"},
			{Value: "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0"},
			{Value: "Contract.my-token"},
			{Value: "1"},
			{Value: "1,2,3"},
		},
	}

	table := &gherkin.DataTable{
		Rows: []*gherkin.TableRow{headers, row1, row2},
	}

	evlps, err := p.ParseEnvelopes("test-1", table)
	assert.NoError(t, err, "ParseEnvelopes should not error")
	assert.Equal(t, "17", evlps[0].GetChain().GetBigChainID().String(), "#1 chain id should be correct")
	assert.Equal(t, "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0", evlps[0].GetFrom().Hex(), "#1 chain id should be correct")
	assert.Equal(t, "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8", evlps[0].GetTx().GetTxData().GetTo().Hex(), "#1 chain id should be correct")
	assert.Equal(t, "type:QUORUM_TESSERA ", evlps[0].GetProtocol().String(), "#1 chain id should be correct")
	assert.Equal(t, []string{"1", "2"}, evlps[0].GetArgs().GetCall().GetArgs(), "#1 args should be correct")

	assert.Equal(t, "888", evlps[1].GetChain().GetBigChainID().String(), "#2 chain id should be correct")
	assert.Equal(t, "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0", evlps[1].GetFrom().Hex(), "#2 chain id should be correct")
	assert.Equal(t, "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8", evlps[1].GetTx().GetTxData().GetTo().Hex(), "#2 chain id should be correct")
	assert.Equal(t, "type:QUORUM_CONSTELLATION ", evlps[1].GetProtocol().String(), "#2 chain id should be correct")

	evlps, err = p.ParseEnvelopes("test-2", table)
	assert.NoError(t, err, "ParseEnvelopes should not error")
	assert.Equal(t, "888", evlps[1].GetChain().GetBigChainID().String(), "#3 chain id should be correct")
	assert.Equal(t, "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0", evlps[0].GetFrom().Hex(), "#3 chain id should be correct")
	assert.Equal(t, "0x77F888CC34a3E6EC4935eF27a83a48fAe548fa4d", evlps[0].GetTx().GetTxData().GetTo().Hex(), "#3 chain id should be correct")
}
