package tx

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	error1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type Builder struct {
	ID            string `validate:"uuid4,required"`
	Headers       map[string]string
	ContextLabels map[string]string
	Method
	Tx             `mapstructure:",squash"`
	Chain          `mapstructure:",squash"`
	Contract       `mapstructure:",squash"`
	Private        `mapstructure:",squash"`
	Receipt        *ethereum.Receipt
	Errors         []*error1.Error
	InternalLabels map[string]string
}

func NewBuilder() *Builder {
	return &Builder{
		Headers:        make(map[string]string),
		ContextLabels:  make(map[string]string),
		Errors:         make([]*error1.Error, 0),
		InternalLabels: make(map[string]string),
	}
}

func (b *Builder) GetID() string {
	return b.ID
}

func (b *Builder) SetID(id string) *Builder {
	b.ID = id
	return b
}

func (b *Builder) GetErrors() []*error1.Error {
	return b.Errors
}

// Error returns string representation of errors encountered by envelope
func (b *Builder) Error() string {
	if len(b.GetErrors()) == 0 {
		return ""
	}
	return fmt.Sprintf("%q", b.GetErrors())
}

func (b *Builder) AppendError(err *error1.Error) *Builder {
	b.Errors = append(b.Errors, err)
	return b
}

func (b *Builder) AppendErrors(errs []*error1.Error) *Builder {
	b.Errors = append(b.Errors, errs...)
	return b
}
func (b *Builder) SetReceipt(receipt *ethereum.Receipt) *Builder {
	b.Receipt = receipt
	return b
}
func (b *Builder) GetReceipt() *ethereum.Receipt {
	return b.Receipt
}

func (b *Builder) GetMethod() Method {
	return b.Method
}

func (b *Builder) SetMethod(method Method) *Builder {
	b.Method = method
	return b
}

// IsEthSendRawTransaction for a classic Ethereum transaction
func (b *Builder) IsEthSendRawTransaction() bool {
	return b.Method == Method_ETH_SENDRAWTRANSACTION
}

// IsEthSendPrivateTransaction for Quorum Constellation
func (b *Builder) IsEthSendPrivateTransaction() bool {
	return b.Method == Method_ETH_SENDPRIVATETRANSACTION
}

// IsEthSendRawPrivateTransaction for Quorum Tessera
func (b *Builder) IsEthSendRawPrivateTransaction() bool {
	return b.Method == Method_ETH_SENDRAWPRIVATETRANSACTION
}

// IsEthSendRawTransaction for Besu Orion
func (b *Builder) IsEeaSendPrivateTransaction() bool {
	return b.Method == Method_EEA_SENDPRIVATETRANSACTION
}

func (b *Builder) Carrier() opentracing.TextMapCarrier {
	return b.ContextLabels
}

func (b *Builder) OnlyWarnings() bool {
	for _, err := range b.GetErrors() {
		if !errors.IsWarning(err) {
			return false
		}
	}
	return true
}

func (b *Builder) GetHeaders() map[string]string {
	return b.Headers
}
func (b *Builder) GetHeadersValue(key string) string {
	return b.Headers[key]
}
func (b *Builder) SetHeadersValue(key, value string) *Builder {
	b.Headers[key] = value
	return b
}

func (b *Builder) GetInternalLabels() map[string]string {
	return b.InternalLabels
}

func (b *Builder) GetInternalLabelsValue(key string) string {
	return b.InternalLabels[key]
}

func (b *Builder) SetInternalLabelsValue(key, value string) *Builder {
	b.InternalLabels[key] = value
	return b
}

func (b *Builder) SetContextLabelsValue(key, value string) *Builder {
	b.ContextLabels[key] = value
	return b
}

func (b *Builder) SetContextLabels(ctxLabels map[string]string) *Builder {
	b.ContextLabels = ctxLabels
	return b
}

func (b *Builder) Validate() []error {
	err := utils.GetValidator().Struct(b)
	if err != nil {
		return utils.HandleValidatorError(err.(validator.ValidationErrors))
	}
	return nil
}
func (b *Builder) GetContextLabelsValue(key string) string {
	return b.ContextLabels[key]
}
func (b *Builder) GetContextLabels() map[string]string {
	return b.ContextLabels
}

type Tx struct {
	From     *ethcommon.Address
	To       *ethcommon.Address
	Gas      *uint64
	GasPrice *big.Int
	Value    *big.Int
	Nonce    *uint64
	Data     string          `validate:"omitempty,isHex"`
	Raw      string          `validate:"omitempty,isHex,required_with_all=TxHash"`
	TxHash   *ethcommon.Hash `validate:"omitempty,required_with_all=Raw"`
}

func (b *Builder) GetTransaction() (*ethtypes.Transaction, error) {
	nonce, err := b.GetNonceUint64()
	if err != nil {
		return nil, err
	}
	value, err := b.GetValueBig()
	if value == nil || err != nil {
		_ = b.SetValue(big.NewInt(0))
	}
	gas, err := b.GetGasUint64()
	if err != nil {
		return nil, err
	}
	gasPrice, err := b.GetGasPriceBig()
	if err != nil {
		return nil, err
	}

	if b.IsConstructor() {
		// Create contract deployment transaction
		return ethtypes.NewContractCreation(
			nonce,
			value,
			gas,
			gasPrice,
			b.MustGetDataBytes(),
		), nil
	}

	to, err := b.GetToAddress()
	if err != nil {
		return nil, err
	}

	// Create transaction
	return ethtypes.NewTransaction(
		nonce,
		to,
		value,
		gas,
		gasPrice,
		b.MustGetDataBytes(),
	), nil
}

// FROM

func (b *Builder) GetFrom() *ethcommon.Address {
	return b.From
}

func (b *Builder) GetFromAddress() (ethcommon.Address, error) {
	if b.From == nil {
		return ethcommon.Address{}, errors.DataError("no from is filled")
	}
	return *b.From, nil
}

func (b *Builder) MustGetFromAddress() ethcommon.Address {
	if b.From == nil {
		return ethcommon.Address{}
	}
	return *b.From
}

func (b *Builder) GetFromString() string {
	if b.From == nil {
		return ""
	}
	return b.From.Hex()
}

func (b *Builder) SetFromString(from string) error {
	if from != "" {
		if !ethcommon.IsHexAddress(from) {
			return errors.DataError("invalid from - got %s", from)
		}
		_ = b.SetFrom(ethcommon.HexToAddress(from))
	}
	return nil
}

func (b *Builder) MustSetFromString(from string) *Builder {
	_ = b.SetFrom(ethcommon.HexToAddress(from))
	return b
}

func (b *Builder) SetFrom(from ethcommon.Address) *Builder {
	b.From = &from
	return b
}

// TO

func (b *Builder) GetTo() *ethcommon.Address {
	return b.To
}

func (b *Builder) GetToAddress() (ethcommon.Address, error) {
	if b.To == nil {
		return ethcommon.Address{}, errors.DataError("no to is filled")
	}
	return *b.To, nil
}

func (b *Builder) MustGetToAddress() ethcommon.Address {
	if b.To == nil {
		return ethcommon.Address{}
	}
	return *b.To
}

func (b *Builder) GetToString() string {
	if b.To == nil {
		return ""
	}
	return b.To.Hex()
}

func (b *Builder) MustSetToString(to string) *Builder {
	_ = b.SetTo(ethcommon.HexToAddress(to))
	return b
}

func (b *Builder) SetToString(to string) error {
	if to != "" {
		if !ethcommon.IsHexAddress(to) {
			return errors.DataError("invalid to - got %s", to)
		}
		_ = b.SetTo(ethcommon.HexToAddress(to))
	}
	return nil
}

func (b *Builder) SetTo(to ethcommon.Address) *Builder {
	b.To = &to
	return b
}

// GAS

func (b *Builder) GetGas() *uint64 {
	return b.Gas
}
func (b *Builder) GetGasUint64() (uint64, error) {
	if b.Gas == nil {
		return 0, errors.DataError("no gas is filled")
	}
	return *b.Gas, nil
}
func (b *Builder) MustGetGasUint64() uint64 {
	if b.Gas == nil {
		return 0
	}
	return *b.Gas
}
func (b *Builder) GetGasString() string {
	if b.Gas == nil {
		return ""
	}
	return strconv.FormatUint(*b.Gas, 10)
}
func (b *Builder) SetGasString(gas string) error {
	if gas != "" {
		g, err := strconv.ParseUint(gas, 10, 32)
		if err != nil {
			return errors.DataError("invalid gasPrice - got %s", gas)
		}
		_ = b.SetGas(g)
	}
	return nil
}
func (b *Builder) SetGas(gas uint64) *Builder {
	b.Gas = &(&struct{ x uint64 }{gas}).x
	return b
}

// NONCE

func (b *Builder) GetNonce() *uint64 {
	return b.Nonce
}
func (b *Builder) GetNonceUint64() (uint64, error) {
	if b.Nonce == nil {
		return 0, errors.DataError("no nonce is filled")
	}
	return *b.Nonce, nil
}
func (b *Builder) MustGetNonceUint64() uint64 {
	if b.Nonce == nil {
		return 0
	}
	return *b.Nonce
}
func (b *Builder) GetNonceString() string {
	if b.Nonce == nil {
		return ""
	}
	return strconv.FormatUint(*b.Nonce, 10)
}
func (b *Builder) SetNonceString(nonce string) error {
	if nonce != "" {
		g, err := strconv.ParseUint(nonce, 10, 32)
		if err != nil {
			return errors.DataError("invalid nonce - got %s", nonce)
		}
		_ = b.SetNonce(g)
	}
	return nil
}
func (b *Builder) SetNonce(nonce uint64) *Builder {
	b.Nonce = &(&struct{ x uint64 }{nonce}).x
	return b
}

// GASPRICE

func (b *Builder) GetGasPrice() *big.Int {
	return b.GasPrice
}

func (b *Builder) GetGasPriceBig() (*big.Int, error) {
	if b.GasPrice == nil {
		return nil, errors.DataError("no gasPrice is filled")
	}
	return b.GasPrice, nil
}

func (b *Builder) GetGasPriceString() string {
	if b.GasPrice == nil {
		return ""
	}
	return b.GasPrice.String()
}

func (b *Builder) SetGasPriceString(gasPrice string) error {
	if gasPrice != "" {
		g, ok := new(big.Int).SetString(gasPrice, 10)
		if !ok {
			return errors.DataError("invalid gasPrice - got %s", gasPrice)
		}
		_ = b.SetGasPrice(g)
	}
	return nil
}

func (b *Builder) SetGasPrice(gasPrice *big.Int) *Builder {
	b.GasPrice = gasPrice
	return b
}

// VALUE

func (b *Builder) GetValue() *big.Int {
	return b.Value
}
func (b *Builder) GetValueBig() (*big.Int, error) {
	if b.Value == nil {
		return nil, errors.DataError("no value is filled")
	}
	return b.Value, nil
}

func (b *Builder) GetValueString() string {
	if b.Value == nil {
		return ""
	}
	return b.Value.String()
}

func (b *Builder) SetValueString(value string) error {
	if value != "" {
		v, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return errors.DataError("invalid value - got %s", value)
		}
		_ = b.SetValue(v)
	}
	return nil
}

func (b *Builder) SetValue(value *big.Int) *Builder {
	b.Value = value
	return b
}

// DATA

func (b *Builder) GetData() string {
	return b.Data
}

func (b *Builder) MustGetDataBytes() []byte {
	if b.Data == "" {
		return []byte{}
	}
	data, _ := hexutil.Decode(b.Data)
	return data
}

func (b *Builder) SetData(data []byte) *Builder {
	b.Data = hexutil.Encode(data)
	return b
}

func (b *Builder) SetDataString(data string) error {
	_, err := hexutil.Decode(data)
	if err != nil {
		return errors.DataError("invalid data")
	}
	b.Data = data
	return nil
}

func (b *Builder) MustSetDataString(data string) *Builder {
	b.Data = data
	return b
}

// RAW

func (b *Builder) GetShortRaw() string {
	return utils.ShortString(b.Raw, 30)
}

func (b *Builder) GetRaw() string {
	return b.Raw
}

func (b *Builder) SetRaw(raw []byte) *Builder {
	b.Raw = hexutil.Encode(raw)
	return b
}

func (b *Builder) SetRawString(raw string) error {
	_, err := hexutil.Decode(raw)
	if err != nil {
		return errors.DataError("invalid raw")
	}
	b.Raw = raw
	return nil
}

func (b *Builder) MustSetRawString(raw string) *Builder {
	b.Raw = raw
	return b
}

// TXHASH

func (b *Builder) GetTxHash() *ethcommon.Hash {
	return b.TxHash
}

func (b *Builder) GetTxHashValue() (ethcommon.Hash, error) {
	if b.TxHash == nil {
		return ethcommon.Hash{}, errors.DataError("no tx hash is filled")
	}
	return *b.TxHash, nil
}

func (b *Builder) MustGetTxHashValue() ethcommon.Hash {
	if b.TxHash == nil {
		return ethcommon.Hash{}
	}
	return *b.TxHash
}

func (b *Builder) GetTxHashString() string {
	if b.TxHash == nil {
		return ""
	}
	return b.TxHash.Hex()
}

func (b *Builder) SetTxHash(hash ethcommon.Hash) *Builder {
	b.TxHash = &hash
	return b
}

func (b *Builder) SetTxHashString(txHash string) error {
	if txHash != "" {
		h, err := hexutil.Decode(txHash)
		if err != nil || len(h) != ethcommon.HashLength {
			return errors.DataError("invalid txHash - got %s", txHash)
		}
		_ = b.SetTxHash(ethcommon.BytesToHash(h))
	}
	return nil
}

func (b *Builder) MustSetTxHashString(txHash string) *Builder {
	_ = b.SetTxHash(ethcommon.HexToHash(txHash))
	return b
}

type Chain struct {
	ChainID   *big.Int
	ChainName string
	ChainUUID string `validate:"omitempty,uuid4"`
}

func (b *Builder) GetChainID() *big.Int {
	return b.ChainID
}

func (b *Builder) GetChainIDString() string {
	if b.ChainID == nil {
		return ""
	}
	return b.ChainID.String()
}

func (b *Builder) SetChainID(chainID *big.Int) *Builder {
	b.ChainID = chainID
	return b
}

func (b *Builder) SetChainIDUint64(chainID uint64) *Builder {
	b.ChainID = big.NewInt(int64(chainID))
	return b
}

func (b *Builder) SetChainIDString(chainID string) error {
	if chainID != "" {
		v, ok := new(big.Int).SetString(chainID, 10)
		if !ok {
			return errors.DataError("invalid chainID - got %s", chainID)
		}
		_ = b.SetChainID(v)
	}
	return nil
}

func (b *Builder) GetChainName() string {
	return b.ChainName
}

func (b *Builder) SetChainName(chainName string) *Builder {
	b.ChainName = chainName
	return b
}

func (b *Builder) GetChainUUID() string {
	return b.ChainUUID
}

func (b *Builder) SetChainUUID(chainUUID string) *Builder {
	b.ChainUUID = chainUUID
	return b
}

type Contract struct {
	ContractName    string `validate:"omitempty,required_with_all=ContractTag"`
	ContractTag     string `validate:"omitempty"`
	MethodSignature string `validate:"omitempty,isValidMethodSig"`
	Args            []string
}

func (b *Builder) GetContractID() *abi.ContractId {
	return &abi.ContractId{
		Name: b.ContractName,
		Tag:  b.ContractTag,
	}
}

// IsConstructor indicate whether the method refers to a deployment
func (b *Builder) IsConstructor() bool {
	return b.GetMethodName() == "constructor"
}

// Short returns a short string representation of contract information
func (b *Builder) GetMethodName() string {
	return strings.Split(b.MethodSignature, "(")[0]
}

func (b *Builder) GetMethodSignature() string {
	return b.MethodSignature
}

func (b *Builder) GetArgs() []string {
	return b.Args
}

func (b *Builder) SetContractName(contractName string) *Builder {
	b.ContractName = contractName
	return b
}

func (b *Builder) SetMethodSignature(methodSignature string) *Builder {
	b.MethodSignature = methodSignature
	return b
}
func (b *Builder) SetArgs(args []string) *Builder {
	b.Args = args
	return b
}

func (b *Builder) SetContractTag(contractTag string) *Builder {
	b.ContractTag = contractTag
	return b
}

func (b *Builder) ShortContract() string {
	if b.ContractName == "" {
		return ""
	}

	if b.ContractTag == "" {
		return b.ContractName
	}

	return fmt.Sprintf("%v[%v]", b.ContractName, b.ContractTag)
}

type Private struct {
	PrivateFor     []string `validate:"dive,base64"`
	PrivateFrom    string   `validate:"omitempty,base64"`
	PrivateTxType  string
	PrivacyGroupID string
}

func (b *Builder) GetPrivateFor() []string {
	return b.PrivateFor
}
func (b *Builder) SetPrivateFor(privateFor []string) *Builder {
	b.PrivateFor = privateFor
	return b
}

func (b *Builder) SetPrivateFrom(privateFrom string) *Builder {
	b.PrivateFrom = privateFrom
	return b
}
func (b *Builder) GetPrivateFrom() string {
	return b.PrivateFrom
}

func (b *Builder) TxRequest() *TxRequest {
	req := &TxRequest{
		Id:      b.ID,
		Headers: b.Headers,
		Chain:   b.ChainName,
		Method:  b.Method,
		Params: &Params{
			From:            b.GetFromString(),
			To:              b.GetToString(),
			Gas:             b.GetGasString(),
			GasPrice:        b.GetGasPriceString(),
			Value:           b.GetValueString(),
			Nonce:           b.GetNonceString(),
			Data:            b.Data,
			Contract:        b.ShortContract(),
			MethodSignature: b.MethodSignature,
			Args:            b.Args,
			Raw:             b.Raw,
			PrivateFor:      b.PrivateFor,
			PrivateFrom:     b.PrivateFrom,
			PrivateTxType:   b.PrivateTxType,
			PrivacyGroupId:  b.PrivacyGroupID,
		},
		ContextLabels: b.ContextLabels,
	}

	return req
}

func (b *Builder) fieldsToInternal() {
	if b.GetChainID() != nil {
		b.InternalLabels["chainID"] = b.GetChainIDString()
	}
	if b.GetTxHash() != nil {
		b.InternalLabels["txHash"] = b.GetTxHashString()
	}
	if b.GetChainUUID() != "" {
		b.InternalLabels["chainUUID"] = b.GetChainUUID()
	}
}

func (b *Builder) internalToFields() error {
	if err := b.SetTxHashString(b.InternalLabels["txHash"]); err != nil {
		return err
	}
	if err := b.SetChainIDString(b.InternalLabels["chainID"]); err != nil {
		return err
	}
	_ = b.SetChainUUID(b.InternalLabels["chainUUID"])
	return nil
}

func (b *Builder) TxEnvelopeAsRequest() *TxEnvelope {
	b.fieldsToInternal()
	return &TxEnvelope{
		InternalLabels: b.InternalLabels,
		Msg:            &TxEnvelope_TxRequest{b.TxRequest()},
	}
}

func (b *Builder) TxEnvelopeAsResponse() *TxEnvelope {
	b.fieldsToInternal()
	return &TxEnvelope{
		InternalLabels: b.InternalLabels,
		Msg:            &TxEnvelope_TxResponse{b.TxResponse()},
	}
}

func (b *Builder) TxResponse() *TxResponse {
	res := &TxResponse{
		Headers:       b.Headers,
		Id:            b.ID,
		ContextLabels: b.ContextLabels,
		Transaction: &ethereum.Transaction{
			From:     b.GetFromString(),
			Nonce:    b.GetNonceString(),
			To:       b.GetToString(),
			Value:    b.GetValueString(),
			Gas:      b.GetGasString(),
			GasPrice: b.GetGasPriceString(),
			Data:     b.GetData(),
			TxHash:   b.GetTxHashString(),
		},
		Receipt: b.Receipt,
		Errors:  b.Errors,
	}

	return res
}

func (b *Builder) loadPtrFields(gas, nonce, gasPrice, value, from, to string) []*error1.Error {
	errs := make([]*error1.Error, 0)
	if err := b.SetGasString(gas); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := b.SetNonceString(nonce); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := b.SetGasPriceString(gasPrice); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := b.SetValueString(value); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := b.SetFromString(from); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := b.SetToString(to); err != nil {
		errs = append(errs, errors.FromError(err))
	}

	return errs
}
