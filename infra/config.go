package infra

// FaucetConfig is the configuration for a Faucet
type FaucetConfig struct {
	Addresses    map[string]string `long:"faucet-address" env:"FAUCET_ADDRESS" default:"3:0x7E654d251Da770A068413677967F6d3Ea2FeA9E4" description:"Address of the faucet for a chain <chainID>-<Address> (e.g 3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4)"`
	CreditAmount string            `long:"faucet-credit-amount" env:"FAUCET_CREDIT_AMOUNT" default:"100000000000000000" description:"Amount to credit when calling faucet"`

	CoolDownTime string   `long:"faucet-cd-time" env:"FAUCET_COOLDOWN_TIME" default:"60s" description:"Cooldown time. Valid time units are ns, us (or µs), ms, s, m, h"`
	BlackList    []string `long:"faucet-blacklist" env:"FAUCET_BLACKLIST" default:"3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4" description:"Blacklisted <chainID>-<Account> (e.g 3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4)"`
	MaxBalance   string   `long:"faucet-max-balance" env:"FAUCET_MAX_BALANCE" default:"500000000000000000000" description:"Max possible balance in decimal format"`

	Topic string `long:"faucet-topic" env:"KAFKA_TOPIC_TX_CRAFTER" default:"topic-tx-crafter" description:"Kafka topic to send credit request to"`
}
