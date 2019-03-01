# Tx-Crafter

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production ready and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **grpc**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to the input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Crafter is a Core-Stack worker responsible to 

- **Craft transaction payload**  
- **Set Gas Price** of the transaction
- **Set Gas Limit** of the transaction
- **Request Faucet crediting**

It consumes message from *tx crafting* Kafka topic and publish to *tx nonce* topic.

## Quick-Start

### Prerequisite

Having *docker* and *docker-compose* install

### Start the application

You can start the application with default configuration by running

```sh
$ docker-compose up
```

### Configure application

Application can be configured through flags or environment variables, you can run the ```help run``` command line


```text
Run application

Usage:
  app run [flags]

Flags:
      --abi strings                Smart Contract ABIs to register for crafting
                                   Environment variable: "ABI" (default ["ERC1400:[{""constant"":true,""inputs"":[],""name"":""name"",""outputs"":[{""name"":"""",""type"":""string""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""operator"",""type"":""address""}],""name"":""authorizeOperatorByPartition"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""operator"",""type"":""address""}],""name"":""revokeOperatorByPartition"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""totalSupply"",""outputs"":[{""name"":"""",""type"":""uint256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""transferWithData"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""tokenHolder"",""type"":""address""}],""name"":""balanceOfByPartition"",""outputs"":[{""name"":"""",""type"":""uint256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""granularity"",""outputs"":[{""name"":"""",""type"":""uint256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""sender"",""type"":""address""}],""name"":""checkCount"",""outputs"":[{""name"":"""",""type"":""uint256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""totalPartitions"",""outputs"":[{""name"":"""",""type"":""bytes32[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""operator"",""type"":""address""},{""name"":""tokenHolder"",""type"":""address""}],""name"":""isOperatorForPartition"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""tokenHolder"",""type"":""address""}],""name"":""balanceOf"",""outputs"":[{""name"":"""",""type"":""uint256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[],""name"":""renounceOwnership"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""operator"",""type"":""address""}],""name"":""certificateSigners"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""tokenHolder"",""type"":""address""}],""name"":""partitionsOf"",""outputs"":[{""name"":"""",""type"":""bytes32[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""controllers"",""outputs"":[{""name"":"""",""type"":""address[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""partition"",""type"":""bytes32""}],""name"":""controllersByPartition"",""outputs"":[{""name"":"""",""type"":""address[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""from"",""type"":""address""},{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""},{""name"":""operatorData"",""type"":""bytes""}],""name"":""transferFromWithData"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""from"",""type"":""address""},{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""},{""name"":""operatorData"",""type"":""bytes""}],""name"":""operatorTransferByPartition"",""outputs"":[{""name"":"""",""type"":""bytes32""}],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""owner"",""outputs"":[{""name"":"""",""type"":""address""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""isOwner"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""operator"",""type"":""address""}],""name"":""authorizeOperator"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""symbol"",""outputs"":[{""name"":"""",""type"":""string""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""account"",""type"":""address""}],""name"":""addMinter"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[],""name"":""renounceMinter"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""account"",""type"":""address""}],""name"":""isMinter"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""tokenHolder"",""type"":""address""}],""name"":""getDefaultPartitions"",""outputs"":[{""name"":"""",""type"":""bytes32[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""operator"",""type"":""address""},{""name"":""tokenHolder"",""type"":""address""}],""name"":""isOperatorFor"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partitions"",""type"":""bytes32[]""}],""name"":""setDefaultPartitions"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""newOwner"",""type"":""address""}],""name"":""transferOwnership"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""transferByPartition"",""outputs"":[{""name"":"""",""type"":""bytes32""}],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""operator"",""type"":""address""}],""name"":""revokeOperator"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""inputs"":[{""name"":""name"",""type"":""string""},{""name"":""symbol"",""type"":""string""},{""name"":""granularity"",""type"":""uint256""},{""name"":""controllers"",""type"":""address[]""},{""name"":""certificateSigner"",""type"":""address""},{""name"":""tokenDefaultPartitions"",""type"":""bytes32[]""}],""payable"":false,""stateMutability"":""nonpayable"",""type"":""constructor""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""account"",""type"":""address""}],""name"":""MinterAdded"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""account"",""type"":""address""}],""name"":""MinterRemoved"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":false,""name"":""sender"",""type"":""address""}],""name"":""Checked"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""previousOwner"",""type"":""address""},{""indexed"":true,""name"":""newOwner"",""type"":""address""}],""name"":""OwnershipTransferred"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""from"",""type"":""address""},{""indexed"":true,""name"":""to"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""TransferWithData"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""to"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""Issued"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""from"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""Redeemed"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""tokenHolder"",""type"":""address""}],""name"":""AuthorizedOperator"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""tokenHolder"",""type"":""address""}],""name"":""RevokedOperator"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""fromPartition"",""type"":""bytes32""},{""indexed"":false,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""from"",""type"":""address""},{""indexed"":true,""name"":""to"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""TransferByPartition"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""fromPartition"",""type"":""bytes32""},{""indexed"":true,""name"":""toPartition"",""type"":""bytes32""},{""indexed"":false,""name"":""value"",""type"":""uint256""}],""name"":""ChangedPartition"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""partition"",""type"":""bytes32""},{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""tokenHolder"",""type"":""address""}],""name"":""AuthorizedOperatorByPartition"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""partition"",""type"":""bytes32""},{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""tokenHolder"",""type"":""address""}],""name"":""RevokedOperatorByPartition"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""name"",""type"":""bytes32""},{""indexed"":false,""name"":""uri"",""type"":""string""},{""indexed"":false,""name"":""documentHash"",""type"":""bytes32""}],""name"":""Document"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""partition"",""type"":""bytes32""},{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""to"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""IssuedByPartition"",""type"":""event""},{""anonymous"":false,""inputs"":[{""indexed"":true,""name"":""partition"",""type"":""bytes32""},{""indexed"":true,""name"":""operator"",""type"":""address""},{""indexed"":true,""name"":""from"",""type"":""address""},{""indexed"":false,""name"":""value"",""type"":""uint256""},{""indexed"":false,""name"":""data"",""type"":""bytes""},{""indexed"":false,""name"":""operatorData"",""type"":""bytes""}],""name"":""RedeemedByPartition"",""type"":""event""},{""constant"":true,""inputs"":[{""name"":""name"",""type"":""bytes32""}],""name"":""getDocument"",""outputs"":[{""name"":"""",""type"":""string""},{""name"":"""",""type"":""bytes32""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""name"",""type"":""bytes32""},{""name"":""uri"",""type"":""string""},{""name"":""documentHash"",""type"":""bytes32""}],""name"":""setDocument"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""isControllable"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""isIssuable"",""outputs"":[{""name"":"""",""type"":""bool""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""tokenHolder"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""issueByPartition"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""redeemByPartition"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""tokenHolder"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""},{""name"":""operatorData"",""type"":""bytes""}],""name"":""operatorRedeemByPartition"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""canTransferByPartition"",""outputs"":[{""name"":"""",""type"":""bytes1""},{""name"":"""",""type"":""bytes32""},{""name"":"""",""type"":""bytes32""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":true,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""from"",""type"":""address""},{""name"":""to"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""},{""name"":""operatorData"",""type"":""bytes""}],""name"":""canOperatorTransferByPartition"",""outputs"":[{""name"":"""",""type"":""bytes1""},{""name"":"""",""type"":""bytes32""},{""name"":"""",""type"":""bytes32""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[],""name"":""renounceControl"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[],""name"":""renounceIssuance"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""operators"",""type"":""address[]""}],""name"":""setControllers"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""partition"",""type"":""bytes32""},{""name"":""operators"",""type"":""address[]""}],""name"":""setPartitionControllers"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""operator"",""type"":""address""},{""name"":""authorized"",""type"":""bool""}],""name"":""setCertificateSigner"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":true,""inputs"":[],""name"":""getTokenDefaultPartitions"",""outputs"":[{""name"":"""",""type"":""bytes32[]""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""defaultPartitions"",""type"":""bytes32[]""}],""name"":""setTokenDefaultPartitions"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""}],""name"":""redeem"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""from"",""type"":""address""},{""name"":""value"",""type"":""uint256""},{""name"":""data"",""type"":""bytes""},{""name"":""operatorData"",""type"":""bytes""}],""name"":""redeemFrom"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""}]"])
      --eth-client strings         Ethereum client URLs.
                                   Environment variable: "ETH_CLIENT_URL" (default [https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7,https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c])
      --faucet-address strings     Faucet address for each chain (format <chainID>:<Address>)
                                   Environment variable: "FAUCET_ADDRESS" (default [3:0x7E654d251Da770A068413677967F6d3Ea2FeA9E4])
      --faucet-amount string       Amount to credit when calling Faucet (Wei in decimal format)
                                   Environment variable: "FAUCET_CREDIT_AMOUNT" (default "100000000000000000")
      --faucet-blacklist strings   Blacklisted address (format <chainID>-<Address>)
                                   Environment variable: "FAUCET_BLACKLIST" (default [3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4])
      --faucet-cooldown duration   Cooldown time.
                                   Environment variable: "FAUCET_COOLDOWN_TIME" (default 1m0s)
      --faucet-max string          Max balance (Wei in decimal format)
                                   Environment variable: "FAUCET_MAX_BALANCE" (default "200000000000000000")
      --faucet-topic string        Kafka topic to send credit request to
                                   Environment variable: "KAFKA_TOPIC_TX_CRAFTER" (default "topic-tx-crafter")
  -h, --help                       help for run
      --http-hostname string       Hostname to expose healthchecks and metrics.
                                   Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --kafka-address strings      Address of Kafka server to connect to.
                                   Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --worker-group string        Kafka consumer group.
                                   Environment variable: "KAFKA_CRAFTER_GROUP" (default "tx-crafter-group")
      --worker-in string           Kafka topic to consume message from.
                                   Environment variable: "KAFKA_TOPIC_TX_CRAFTER" (default "topic-tx-crafter")
      --worker-out string          Kafka topic to send message to after processing.
                                   Environment variable: "KAFKA_TOPIC_TX_NONCE" (default "topic-tx-nonce")
      --worker-slots uint          Maximum number of messages the worker can treat in parallel.
                                   Environment variable: "WORKER_SLOTS" (default 100)

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## High Level Architecture

Tx-Crafter expect all consumed messages to respect [Core-Stack standard protobuf format](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/protobuf)

Consumed messages should have 

- ```Chain``` entry set
- ```Call``` entry set or Tx-Worker will consider the transaction as a basic Ethereum transaction with no payload.

1. **Crafting**

To craft transaction payload Tx-Worker inspects the ```Call``` entry of input protobuf message 
 
- it expect the ```ID``` entry formated as ```<method>@<contract_name>``` (e.g. "transfer@ERC20") (this will evolve to handle versioning of contracts) 
- it expects the ```Args``` entry to be a list of expected ordered arguments in ```string``` format (e.g. ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff", "0xabced"] for *to* and *value* args)

By basing on the ```ID```, Tx-Worker requests method ABI from *ABI registry*, then it casts ```string``` arguments in the expected Solidity type to craft payload.

2. **Gas Price**

Tx-Crafter interogates *Ethereum client* using the identifier in ```Chain``` by calling jsonRPC ```eth_gasPrice```

3. **Gas Cost**

Tx-Crafter interogates *Ethereum client* using the identifier in ```Chain``` by calling jsonRPC ```eth_estimateGas```

4. **Faucet**

It request a credit to *Faucet*

## Technical Architecture

![alt core-stack-architecture](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/blob/master/diagrams/Core_Stack_Architecture.png)
