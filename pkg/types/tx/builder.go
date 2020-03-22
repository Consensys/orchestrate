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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	error1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type Envelope struct {
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

func NewEnvelope() *Envelope {
	return &Envelope{
		Headers:        make(map[string]string),
		ContextLabels:  make(map[string]string),
		Errors:         make([]*error1.Error, 0),
		InternalLabels: make(map[string]string),
	}
}

func (e *Envelope) GetID() string {
	return e.ID
}

func (e *Envelope) SetID(id string) *Envelope {
	e.ID = id
	return e
}

func (e *Envelope) GetErrors() []*error1.Error {
	return e.Errors
}

// Error returns string representation of errors encountered by envelope
func (e *Envelope) Error() string {
	if len(e.GetErrors()) == 0 {
		return ""
	}
	return fmt.Sprintf("%q", e.GetErrors())
}

func (e *Envelope) AppendError(err *error1.Error) *Envelope {
	e.Errors = append(e.Errors, err)
	return e
}

func (e *Envelope) AppendErrors(errs []*error1.Error) *Envelope {
	e.Errors = append(e.Errors, errs...)
	return e
}
func (e *Envelope) SetReceipt(receipt *ethereum.Receipt) *Envelope {
	e.Receipt = receipt
	return e
}
func (e *Envelope) GetReceipt() *ethereum.Receipt {
	return e.Receipt
}

func (e *Envelope) GetMethod() Method {
	return e.Method
}

func (e *Envelope) SetMethod(method Method) *Envelope {
	e.Method = method
	return e
}

// IsEthSendRawTransaction for a classic Ethereum transaction
func (e *Envelope) IsEthSendRawTransaction() bool {
	return e.Method == Method_ETH_SENDRAWTRANSACTION
}

// IsEthSendPrivateTransaction for Quorum Constellation
func (e *Envelope) IsEthSendPrivateTransaction() bool {
	return e.Method == Method_ETH_SENDPRIVATETRANSACTION
}

// IsEthSendRawPrivateTransaction for Quorum Tessera
func (e *Envelope) IsEthSendRawPrivateTransaction() bool {
	return e.Method == Method_ETH_SENDRAWPRIVATETRANSACTION
}

// IsEthSendRawTransaction for Besu Orion
func (e *Envelope) IsEeaSendPrivateTransaction() bool {
	return e.Method == Method_EEA_SENDPRIVATETRANSACTION
}

func (e *Envelope) Carrier() opentracing.TextMapCarrier {
	return e.ContextLabels
}

func (e *Envelope) OnlyWarnings() bool {
	for _, err := range e.GetErrors() {
		if !errors.IsWarning(err) {
			return false
		}
	}
	return true
}

func (e *Envelope) GetHeaders() map[string]string {
	return e.Headers
}

func (e *Envelope) SetHeaders(headers map[string]string) *Envelope {
	if headers != nil {
		e.Headers = headers
	}
	return e
}

func (e *Envelope) GetHeadersValue(key string) string {
	return e.Headers[key]
}
func (e *Envelope) SetHeadersValue(key, value string) *Envelope {
	e.Headers[key] = value
	return e
}

func (e *Envelope) GetInternalLabels() map[string]string {
	return e.InternalLabels
}

func (e *Envelope) GetInternalLabelsValue(key string) string {
	return e.InternalLabels[key]
}

func (e *Envelope) SetInternalLabels(internalLabels map[string]string) *Envelope {
	if internalLabels != nil {
		e.InternalLabels = internalLabels
	}
	return e
}
func (e *Envelope) SetInternalLabelsValue(key, value string) *Envelope {
	e.InternalLabels[key] = value
	return e
}

func (e *Envelope) SetContextLabelsValue(key, value string) *Envelope {
	e.ContextLabels[key] = value
	return e
}

func (e *Envelope) SetContextLabels(ctxLabels map[string]string) *Envelope {
	if ctxLabels != nil {
		e.ContextLabels = ctxLabels
	}
	return e
}

func (e *Envelope) Validate() []error {
	err := utils.GetValidator().Struct(e)
	if err != nil {
		return utils.HandleValidatorError(err.(validator.ValidationErrors))
	}
	return nil
}
func (e *Envelope) GetContextLabelsValue(key string) string {
	return e.ContextLabels[key]
}
func (e *Envelope) GetContextLabels() map[string]string {
	return e.ContextLabels
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

func (e *Envelope) GetTransaction() (*ethtypes.Transaction, error) {
	// TODO: Use custom validation with https://godoc.org/gopkg.in/go-playground/validator.v10#Validate.StructFiltered
	nonce, err := e.GetNonceUint64()
	if err != nil {
		return nil, err
	}
	value, err := e.GetValueBig()
	if value == nil || err != nil {
		_ = e.SetValue(big.NewInt(0))
	}
	gas, err := e.GetGasUint64()
	if err != nil {
		return nil, err
	}
	gasPrice, err := e.GetGasPriceBig()
	if err != nil {
		return nil, err
	}

	if e.IsConstructor() {
		// Create contract deployment transaction
		return ethtypes.NewContractCreation(
			nonce,
			value,
			gas,
			gasPrice,
			e.MustGetDataBytes(),
		), nil
	}

	to, err := e.GetToAddress()
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
		e.MustGetDataBytes(),
	), nil
}

// FROM

func (e *Envelope) GetFrom() *ethcommon.Address {
	return e.From
}

func (e *Envelope) GetFromAddress() (ethcommon.Address, error) {
	if e.From == nil {
		return ethcommon.Address{}, errors.DataError("no from is filled")
	}
	return *e.From, nil
}

func (e *Envelope) MustGetFromAddress() ethcommon.Address {
	if e.From == nil {
		return ethcommon.Address{}
	}
	return *e.From
}

func (e *Envelope) GetFromString() string {
	if e.From == nil {
		return ""
	}
	return e.From.Hex()
}

func (e *Envelope) SetFromString(from string) error {
	if from != "" {
		if !ethcommon.IsHexAddress(from) {
			return errors.DataError("invalid from - got %s", from)
		}
		_ = e.SetFrom(ethcommon.HexToAddress(from))
	}
	return nil
}

func (e *Envelope) MustSetFromString(from string) *Envelope {
	_ = e.SetFrom(ethcommon.HexToAddress(from))
	return e
}

func (e *Envelope) SetFrom(from ethcommon.Address) *Envelope {
	e.From = &from
	return e
}

// TO

func (e *Envelope) GetTo() *ethcommon.Address {
	return e.To
}

func (e *Envelope) GetToAddress() (ethcommon.Address, error) {
	if e.To == nil {
		return ethcommon.Address{}, errors.DataError("no to is filled")
	}
	return *e.To, nil
}

func (e *Envelope) MustGetToAddress() ethcommon.Address {
	if e.To == nil {
		return ethcommon.Address{}
	}
	return *e.To
}

func (e *Envelope) GetToString() string {
	if e.To == nil {
		return ""
	}
	return e.To.Hex()
}

func (e *Envelope) MustSetToString(to string) *Envelope {
	_ = e.SetTo(ethcommon.HexToAddress(to))
	return e
}

func (e *Envelope) SetToString(to string) error {
	if to != "" {
		if !ethcommon.IsHexAddress(to) {
			return errors.DataError("invalid to - got %s", to)
		}
		_ = e.SetTo(ethcommon.HexToAddress(to))
	}
	return nil
}

func (e *Envelope) SetTo(to ethcommon.Address) *Envelope {
	e.To = &to
	return e
}

// GAS

func (e *Envelope) GetGas() *uint64 {
	return e.Gas
}
func (e *Envelope) GetGasUint64() (uint64, error) {
	if e.Gas == nil {
		return 0, errors.DataError("no gas is filled")
	}
	return *e.Gas, nil
}
func (e *Envelope) MustGetGasUint64() uint64 {
	if e.Gas == nil {
		return 0
	}
	return *e.Gas
}
func (e *Envelope) GetGasString() string {
	if e.Gas == nil {
		return ""
	}
	return strconv.FormatUint(*e.Gas, 10)
}
func (e *Envelope) SetGasString(gas string) error {
	if gas != "" {
		g, err := strconv.ParseUint(gas, 10, 32)
		if err != nil {
			return errors.DataError("invalid gasPrice - got %s", gas)
		}
		_ = e.SetGas(g)
	}
	return nil
}
func (e *Envelope) SetGas(gas uint64) *Envelope {
	e.Gas = &(&struct{ x uint64 }{gas}).x
	return e
}

// NONCE

func (e *Envelope) GetNonce() *uint64 {
	return e.Nonce
}
func (e *Envelope) GetNonceUint64() (uint64, error) {
	if e.Nonce == nil {
		return 0, errors.DataError("no nonce is filled")
	}
	return *e.Nonce, nil
}
func (e *Envelope) MustGetNonceUint64() uint64 {
	if e.Nonce == nil {
		return 0
	}
	return *e.Nonce
}
func (e *Envelope) GetNonceString() string {
	if e.Nonce == nil {
		return ""
	}
	return strconv.FormatUint(*e.Nonce, 10)
}
func (e *Envelope) SetNonceString(nonce string) error {
	if nonce != "" {
		g, err := strconv.ParseUint(nonce, 10, 32)
		if err != nil {
			return errors.DataError("invalid nonce - got %s", nonce)
		}
		_ = e.SetNonce(g)
	}
	return nil
}
func (e *Envelope) SetNonce(nonce uint64) *Envelope {
	e.Nonce = &(&struct{ x uint64 }{nonce}).x
	return e
}

// GASPRICE

func (e *Envelope) GetGasPrice() *big.Int {
	return e.GasPrice
}

func (e *Envelope) GetGasPriceBig() (*big.Int, error) {
	if e.GasPrice == nil {
		return nil, errors.DataError("no gasPrice is filled")
	}
	return e.GasPrice, nil
}

func (e *Envelope) GetGasPriceString() string {
	if e.GasPrice == nil {
		return ""
	}
	return e.GasPrice.String()
}

func (e *Envelope) SetGasPriceString(gasPrice string) error {
	if gasPrice != "" {
		g, ok := new(big.Int).SetString(gasPrice, 10)
		if !ok {
			return errors.DataError("invalid gasPrice - got %s", gasPrice)
		}
		_ = e.SetGasPrice(g)
	}
	return nil
}

func (e *Envelope) SetGasPrice(gasPrice *big.Int) *Envelope {
	e.GasPrice = gasPrice
	return e
}

// VALUE

func (e *Envelope) GetValue() *big.Int {
	return e.Value
}
func (e *Envelope) GetValueBig() (*big.Int, error) {
	if e.Value == nil {
		return nil, errors.DataError("no value is filled")
	}
	return e.Value, nil
}

func (e *Envelope) GetValueString() string {
	if e.Value == nil {
		return ""
	}
	return e.Value.String()
}

func (e *Envelope) SetValueString(value string) error {
	if value != "" {
		v, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return errors.DataError("invalid value - got %s", value)
		}
		_ = e.SetValue(v)
	}
	return nil
}

func (e *Envelope) SetValue(value *big.Int) *Envelope {
	e.Value = value
	return e
}

// DATA

func (e *Envelope) GetData() string {
	return e.Data
}

func (e *Envelope) MustGetDataBytes() []byte {
	if e.Data == "" {
		return []byte{}
	}
	data, _ := hexutil.Decode(e.Data)
	return data
}

func (e *Envelope) SetData(data []byte) *Envelope {
	e.Data = hexutil.Encode(data)
	return e
}

func (e *Envelope) SetDataString(data string) error {
	_, err := hexutil.Decode(data)
	if err != nil {
		return errors.DataError("invalid data")
	}
	e.Data = data
	return nil
}

func (e *Envelope) MustSetDataString(data string) *Envelope {
	e.Data = data
	return e
}

// RAW

func (e *Envelope) GetShortRaw() string {
	return utils.ShortString(e.Raw, 30)
}

func (e *Envelope) GetRaw() string {
	return e.Raw
}

func (e *Envelope) MustGetRawBytes() []byte {
	if e.Raw == "" {
		return []byte{}
	}
	raw, _ := hexutil.Decode(e.Raw)
	return raw
}

func (e *Envelope) SetRaw(raw []byte) *Envelope {
	e.Raw = hexutil.Encode(raw)
	return e
}

func (e *Envelope) SetRawString(raw string) error {
	_, err := hexutil.Decode(raw)
	if err != nil {
		return errors.DataError("invalid raw")
	}
	e.Raw = raw
	return nil
}

func (e *Envelope) MustSetRawString(raw string) *Envelope {
	e.Raw = raw
	return e
}

// TXHASH

func (e *Envelope) GetTxHash() *ethcommon.Hash {
	return e.TxHash
}

func (e *Envelope) GetTxHashValue() (ethcommon.Hash, error) {
	if e.TxHash == nil {
		return ethcommon.Hash{}, errors.DataError("no tx hash is filled")
	}
	return *e.TxHash, nil
}

func (e *Envelope) MustGetTxHashValue() ethcommon.Hash {
	if e.TxHash == nil {
		return ethcommon.Hash{}
	}
	return *e.TxHash
}

func (e *Envelope) GetTxHashString() string {
	if e.TxHash == nil {
		return ""
	}
	return e.TxHash.Hex()
}

func (e *Envelope) SetTxHash(hash ethcommon.Hash) *Envelope {
	e.TxHash = &hash
	return e
}

func (e *Envelope) SetTxHashString(txHash string) error {
	if txHash != "" {
		h, err := hexutil.Decode(txHash)
		if err != nil || len(h) != ethcommon.HashLength {
			return errors.DataError("invalid txHash - got %s", txHash)
		}
		_ = e.SetTxHash(ethcommon.BytesToHash(h))
	}
	return nil
}

func (e *Envelope) MustSetTxHashString(txHash string) *Envelope {
	_ = e.SetTxHash(ethcommon.HexToHash(txHash))
	return e
}

type Chain struct {
	ChainID   *big.Int
	ChainName string
	ChainUUID string `validate:"omitempty,uuid4"`
}

func (e *Envelope) GetChainID() *big.Int {
	return e.ChainID
}

func (e *Envelope) GetChainIDString() string {
	if e.ChainID == nil {
		return ""
	}
	return e.ChainID.String()
}

func (e *Envelope) SetChainID(chainID *big.Int) *Envelope {
	e.ChainID = chainID
	return e
}

func (e *Envelope) SetChainIDUint64(chainID uint64) *Envelope {
	e.ChainID = big.NewInt(int64(chainID))
	return e
}

func (e *Envelope) SetChainIDString(chainID string) error {
	if chainID != "" {
		v, ok := new(big.Int).SetString(chainID, 10)
		if !ok {
			return errors.DataError("invalid chainID - got %s", chainID)
		}
		_ = e.SetChainID(v)
	}
	return nil
}

func (e *Envelope) GetChainName() string {
	return e.ChainName
}

func (e *Envelope) SetChainName(chainName string) *Envelope {
	e.ChainName = chainName
	return e
}

func (e *Envelope) GetChainUUID() string {
	return e.ChainUUID
}

func (e *Envelope) SetChainUUID(chainUUID string) *Envelope {
	e.ChainUUID = chainUUID
	return e
}

type Contract struct {
	ContractName    string `validate:"omitempty,required_with_all=ContractTag"`
	ContractTag     string `validate:"omitempty"`
	MethodSignature string `validate:"omitempty,isValidMethodSig"`
	Args            []string
}

func (e *Envelope) GetContractID() *abi.ContractId {
	return &abi.ContractId{
		Name: e.ContractName,
		Tag:  e.ContractTag,
	}
}

// IsConstructor indicate whether the method refers to a deployment
func (e *Envelope) IsConstructor() bool {
	return e.MustGetMethodName() == "constructor"
}

// Short returns a short string representation of contract information
func (e *Envelope) MustGetMethodName() string {
	return strings.Split(e.MethodSignature, "(")[0]
}

func (e *Envelope) GetMethodSignature() string {
	return e.MethodSignature
}

func (e *Envelope) GetArgs() []string {
	return e.Args
}

func (e *Envelope) SetContractName(contractName string) *Envelope {
	e.ContractName = contractName
	return e
}

func (e *Envelope) SetMethodSignature(methodSignature string) *Envelope {
	e.MethodSignature = methodSignature
	return e
}
func (e *Envelope) SetArgs(args []string) *Envelope {
	e.Args = args
	return e
}

func (e *Envelope) SetContractTag(contractTag string) *Envelope {
	e.ContractTag = contractTag
	return e
}

func (e *Envelope) ShortContract() string {
	if e.ContractName == "" {
		return ""
	}

	if e.ContractTag == "" {
		return e.ContractName
	}

	return fmt.Sprintf("%v[%v]", e.ContractName, e.ContractTag)
}

type Private struct {
	PrivateFor     []string `validate:"dive,base64"`
	PrivateFrom    string   `validate:"omitempty,base64"`
	PrivateTxType  string
	PrivacyGroupID string
}

func (e *Envelope) GetPrivateFor() []string {
	return e.PrivateFor
}
func (e *Envelope) SetPrivateFor(privateFor []string) *Envelope {
	e.PrivateFor = privateFor
	return e
}

func (e *Envelope) SetPrivateFrom(privateFrom string) *Envelope {
	e.PrivateFrom = privateFrom
	return e
}
func (e *Envelope) GetPrivateFrom() string {
	return e.PrivateFrom
}

func (e *Envelope) SetPrivateTxType(privateTxType string) *Envelope {
	e.PrivateTxType = privateTxType
	return e
}

func (e *Envelope) GetPrivateTxType() string {
	return e.PrivateTxType
}

func (e *Envelope) SetPrivacyGroupID(privacyGroupID string) *Envelope {
	e.PrivacyGroupID = privacyGroupID
	return e
}

func (e *Envelope) GetPrivacyGroupID() string {
	return e.PrivacyGroupID
}

func (e *Envelope) TxRequest() *TxRequest {
	req := &TxRequest{
		Id:      e.ID,
		Headers: e.Headers,
		Chain:   e.GetChainName(),
		Method:  e.Method,
		Params: &Params{
			From:            e.GetFromString(),
			To:              e.GetToString(),
			Gas:             e.GetGasString(),
			GasPrice:        e.GetGasPriceString(),
			Value:           e.GetValueString(),
			Nonce:           e.GetNonceString(),
			Data:            e.GetData(),
			Contract:        e.ShortContract(),
			MethodSignature: e.GetMethodSignature(),
			Args:            e.GetArgs(),
			Raw:             e.GetRaw(),
			PrivateFor:      e.GetPrivateFor(),
			PrivateFrom:     e.GetPrivateFrom(),
			PrivateTxType:   e.PrivateTxType,
			PrivacyGroupId:  e.PrivacyGroupID,
		},
		ContextLabels: e.ContextLabels,
	}

	return req
}

func (e *Envelope) fieldsToInternal() {
	if e.InternalLabels == nil {
		e.InternalLabels = make(map[string]string)
	}

	if e.GetChainID() != nil {
		e.InternalLabels["chainID"] = e.GetChainIDString()
	}
	if e.GetTxHash() != nil {
		e.InternalLabels["txHash"] = e.GetTxHashString()
	}
	if e.GetChainUUID() != "" {
		e.InternalLabels["chainUUID"] = e.GetChainUUID()
	}
}

func (e *Envelope) internalToFields() error {
	hash, ok := e.InternalLabels["txHash"]
	if err := e.SetTxHashString(hash); err != nil && ok {
		return err
	}
	if err := e.SetChainIDString(e.InternalLabels["chainID"]); err != nil {
		return err
	}
	_ = e.SetChainUUID(e.InternalLabels["chainUUID"])
	return nil
}

func (e *Envelope) TxEnvelopeAsRequest() *TxEnvelope {
	e.fieldsToInternal()
	return &TxEnvelope{
		InternalLabels: e.InternalLabels,
		Msg:            &TxEnvelope_TxRequest{e.TxRequest()},
	}
}

func (e *Envelope) TxEnvelopeAsResponse() *TxEnvelope {
	e.fieldsToInternal()
	return &TxEnvelope{
		InternalLabels: e.InternalLabels,
		Msg:            &TxEnvelope_TxResponse{e.TxResponse()},
	}
}

func (e *Envelope) TxResponse() *TxResponse {
	res := &TxResponse{
		Headers:       e.Headers,
		Id:            e.ID,
		ContextLabels: e.ContextLabels,
		Transaction: &ethereum.Transaction{
			From:     e.GetFromString(),
			Nonce:    e.GetNonceString(),
			To:       e.GetToString(),
			Value:    e.GetValueString(),
			Gas:      e.GetGasString(),
			GasPrice: e.GetGasPriceString(),
			Data:     e.GetData(),
			Raw:      e.GetRaw(),
			TxHash:   e.GetTxHashString(),
		},
		Receipt: e.Receipt,
		Errors:  e.Errors,
	}

	return res
}

func (e *Envelope) loadPtrFields(gas, nonce, gasPrice, value, from, to string) []*error1.Error {
	errs := make([]*error1.Error, 0)
	if err := e.SetGasString(gas); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := e.SetNonceString(nonce); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := e.SetGasPriceString(gasPrice); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := e.SetValueString(value); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := e.SetFromString(from); err != nil {
		errs = append(errs, errors.FromError(err))
	}
	if err := e.SetToString(to); err != nil {
		errs = append(errs, errors.FromError(err))
	}

	return errs
}
