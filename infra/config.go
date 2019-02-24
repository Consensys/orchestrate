package infra

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	faucetAddressFlag     = "faucet-address"
	faucetAddressViperKey = "faucet.addresses"
	faucetAddressDefault  = []string{
		"3:0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
	}
	faucetAddressEnv = "FAUCET_ADDRESS"
)

// InitFlags set flags for Faucet config
func InitFlags(f *pflag.FlagSet) {
	FaucetAddress(f)
	FaucetAmount(f)
	FaucetBlacklist(f)
	FaucetCooldown(f)
	FaucetMax(f)
	FaucetTopic(f)
	ABIs(f)
}

// FaucetAddress register flag for Faucet address
func FaucetAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Faucet address for each chain (format <chainID>:<Address>)
Environment variable: %q`, faucetAddressEnv)
	f.StringSlice(faucetAddressFlag, faucetAddressDefault, desc)
	viper.BindPFlag(faucetAddressViperKey, f.Lookup(faucetAddressFlag))
	viper.BindEnv(faucetAddressViperKey, faucetAddressEnv)
}

var (
	faucetAmountFlag     = "faucet-amount"
	faucetAmountViperKey = "faucet.amount"
	faucetAmountDefault  = "100000000000000000"
	faucetAmountEnv      = "FAUCET_CREDIT_AMOUNT"
)

// FaucetAmount register flag for Faucet Amount
func FaucetAmount(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Amount to credit when calling Faucet (Wei in decimal format)
Environment variable: %q`, faucetAmountEnv)
	f.String(faucetAmountFlag, faucetAmountDefault, desc)
	viper.BindPFlag(faucetAmountViperKey, f.Lookup(faucetAmountFlag))
	viper.BindEnv(faucetAmountViperKey, faucetAmountEnv)
}

var (
	faucetCooldownFlag     = "faucet-cooldown"
	faucetCooldownViperKey = "faucet.cooldown"
	faucetCooldownDefault  = 60 * time.Second
	faucetCooldownEnv      = "FAUCET_COOLDOWN_TIME"
)

// FaucetCooldown register flag for Faucet Cooldown
func FaucetCooldown(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Cooldown time.
Environment variable: %q`, faucetCooldownEnv)
	f.Duration(faucetCooldownFlag, faucetCooldownDefault, desc)
	viper.BindPFlag(faucetCooldownViperKey, f.Lookup(faucetCooldownFlag))
	viper.BindEnv(faucetCooldownViperKey, faucetCooldownEnv)
}

var (
	faucetBlacklistFlag     = "faucet-blacklist"
	faucetBlacklistViperKey = "faucet.blacklist"
	faucetBlacklistDefault  = []string{
		"3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
	}
	faucetBlacklistEnv = "FAUCET_BLACKLIST"
)

// FaucetBlacklist register flag for Faucet address
func FaucetBlacklist(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Blacklisted address (format <chainID>-<Address>)
Environment variable: %q`, faucetBlacklistEnv)
	f.StringSlice(faucetBlacklistFlag, faucetBlacklistDefault, desc)
	viper.BindPFlag(faucetBlacklistViperKey, f.Lookup(faucetBlacklistFlag))
	viper.BindEnv(faucetBlacklistViperKey, faucetBlacklistEnv)
}

var (
	faucetMaxFlag     = "faucet-max"
	faucetMaxViperKey = "faucet.max"
	faucetMaxDefault  = "200000000000000000"
	faucetMaxEnv      = "FAUCET_MAX_BALANCE"
)

// FaucetMax register flag for Faucet Max Balance
func FaucetMax(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Max balance (Wei in decimal format)
Environment variable: %q`, faucetMaxEnv)
	f.String(faucetMaxFlag, faucetMaxDefault, desc)
	viper.BindPFlag(faucetMaxViperKey, f.Lookup(faucetMaxFlag))
	viper.BindEnv(faucetMaxViperKey, faucetMaxEnv)
}

var (
	faucetTopicFlag     = "faucet-topic"
	faucetTopicViperKey = "faucet.topic"
	faucetTopicDefault  = "topic-tx-crafter"
	faucetTopicEnv      = "KAFKA_TOPIC_TX_CRAFTER"
)

// FaucetTopic register flag for Faucet Topic
func FaucetTopic(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Kafka topic to send credit request to
Environment variable: %q`, faucetTopicEnv)
	f.String(faucetTopicFlag, faucetTopicDefault, desc)
	viper.BindPFlag(faucetTopicViperKey, f.Lookup(faucetTopicFlag))
	viper.BindEnv(faucetTopicViperKey, faucetTopicEnv)
}

var (
	abiFlag     = "abi"
	abiViperKey = "abis"
	abiDefault  = []string{
		"ERC1400:[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"authorizeOperatorByPartition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"revokeOperatorByPartition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferWithData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"balanceOfByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"granularity\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"checkCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalPartitions\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"isOperatorForPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"certificateSigners\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"partitionsOf\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"controllers\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"}],\"name\":\"controllersByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"transferFromWithData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorTransferByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"authorizeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"getDefaultPartitions\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"isOperatorFor\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partitions\",\"type\":\"bytes32[]\"}],\"name\":\"setDefaultPartitions\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"revokeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"granularity\",\"type\":\"uint256\"},{\"name\":\"controllers\",\"type\":\"address[]\"},{\"name\":\"certificateSigner\",\"type\":\"address\"},{\"name\":\"tokenDefaultPartitions\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Checked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"TransferWithData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Issued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"AuthorizedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"RevokedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"fromPartition\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"TransferByPartition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"fromPartition\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"toPartition\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"ChangedPartition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"partition\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"AuthorizedOperatorByPartition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"partition\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"RevokedOperatorByPartition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"uri\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"documentHash\",\"type\":\"bytes32\"}],\"name\":\"Document\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"partition\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"IssuedByPartition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"partition\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"RedeemedByPartition\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"name\",\"type\":\"bytes32\"}],\"name\":\"getDocument\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"bytes32\"},{\"name\":\"uri\",\"type\":\"string\"},{\"name\":\"documentHash\",\"type\":\"bytes32\"}],\"name\":\"setDocument\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isControllable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isIssuable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"tokenHolder\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"issueByPartition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"redeemByPartition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"tokenHolder\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorRedeemByPartition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"canTransferByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"canOperatorTransferByPartition\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceControl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceIssuance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"setControllers\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"partition\",\"type\":\"bytes32\"},{\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"setPartitionControllers\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"setCertificateSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenDefaultPartitions\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"defaultPartitions\",\"type\":\"bytes32[]\"}],\"name\":\"setTokenDefaultPartitions\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"redeemFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	}
	abiEnv = "ABI"
)

// ABIs register flag for ABIs
func ABIs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Smart Contract ABIs to register for crafting
Environment variable: %q`, abiEnv)
	f.StringSlice(abiFlag, abiDefault, desc)
	viper.BindPFlag(abiViperKey, f.Lookup(abiFlag))
	viper.BindEnv(abiViperKey, abiEnv)
}
