# QuickStart Guide

## How to install

Download the repository:
```bash
go get -u github.com/consensys/orchestrate
```

Go to downloaded repository:
```bash
cd ${GOPATH}/src/github.com/consensys/orchestrate
```

Install vendors locally:
```bash
make mod-vendor
```

Compile orchestrate binary
```bash
make gobuild
```

### Troubleshooting

**Missing `include/secp256k1.h`**

Package manager `go mod` prunes no *.go files therefore in order to be able to compile using local vendors you have to copy 
over from $GO_PATH missing  `go-ethereum/crypto/secp256k1/**/**.h`. Run the following command:
```bash
$> GO111MODULE=off go get github.com/ethereum/go-ethereum
$> cp -r ${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
    ./vendor/github.com/ethereum/go-ethereum/crypto/secp256k1
```

## Run Orchestrate

### Step 1: Launch deps

Orchestrate requires a bunch of services to persist data and exchange message such as:
- Kafka
- Postgres
- HashiCorp vault
- Redis

To launch those service execute the next command:
```bash
make deps
```

### Step 2: Blockchain deps

Orchestrate also requires at least one blockchain network for instance `Quorum` or `Besu`, and a EVM client such as `geth`. 
For instance we could run besu and geth using the next make commands:
```bash
make geth
make quorum
``` 

### Step 3: Orchestrate

Finally we can run Orchestrate. Firstly, we validate if we have every required service running:
```bash
make bootstrap-deps
``` 

If that exits without errors we can proceed spawning orchestrate:
```bash
make orchestrate
```

Above command is going to initialize one of each of the orchestrate service using the latest compile version located at `./build/bin/orchestrate`. By default
orchestrate apis are launched in the following ports:
- CHAIN_API_HOST: `localhost:8011`, `localhost:8012`
- CONTRACT_REGISTRY_API_HOST: `localhost:8020`, `localhost:8021`, `localhost:8022`
- TRANSACTION_SCHEDULER_API_HOST: `localhost:8031`, `localhost:8032`

**Logging**
To monitor orchestrate logs we can do as follow:

`docker-composer logs -f`

## Orchestrate API

[Orchestrate API Documentation](https://consensys.gitlab.io/client/fr/core-stack/orchestrate/latest/#tag/Chain-Registry/paths/~1chains/post)

### Generate a JWT

**Orchestrate initial setup**
In case of having multi tenancy enabled: 
```bash
export MULTI_TENANCY_ENABLED=1
```

We will also need to define a set of ENV VARIABLES. First we indicate the namespace use within the generate token to store the tenant_id:
```bash
export AUTH_JWT_CLAIMS_NAMESPACE="orchestrate.info"
```

and define the server certificates to encode and verify generate token:
```bash
export AUTH_JWT_CERTIFICATE="MIIDYjCCAkoCCQC9pJWk7qdipjANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJGUjEOMAwGA1UEBwwFUGFyaXMxEjAQBgNVBAoMCUNvbnNlblN5czEQMA4GA1UECwwHUGVnYVN5czEuMCwGA1UEAwwlZTJlLXRlc3RzLm9yY2hlc3RyYXRlLmNvbnNlbnN5cy5wYXJpczAeFw0xOTEyMjcxNjI5MTdaFw0yMDEyMjYxNjI5MTdaMHMxCzAJBgNVBAYTAkZSMQ4wDAYDVQQHDAVQYXJpczESMBAGA1UECgwJQ29uc2VuU3lzMRAwDgYDVQQLDAdQZWdhU3lzMS4wLAYDVQQDDCVlMmUtdGVzdHMub3JjaGVzdHJhdGUuY29uc2Vuc3lzLnBhcmlzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo0NqWqI3TSi1uOBvCUquclWo4LcsYT21tNUXQ8YyqVYRSsiBv+ZKZBCjD8XklLPih40kFSe+r6DNca5/LH/okQIdc8nsQg+BLCkXeH2NFv+QYtPczAw4YhS6GVxJk3u9sfp8NavWBcQbD3MMDpehMOvhSl0zoP/ZlH6ErKHNtoQgUpPNVQGysNU21KpClmIDD/L1drsbq+rFiDrcVWaOLwGxr8SBd/0b4ngtcwH16RJaxcIXXT5AVia1CNdzmU5/AIg3OfgzvKn5AGrMZBsmGAiCyn4/P3PnuF81/WHukk5ETLnzOH+vC2elSmZ8y80HCGeqOiQ1rs66L936wX8cDwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCNcTs3n/Ps+yIZDH7utxTOaqpDTCB10MzPmb22UAal89couIT6R0fAu14p/LTkxdb2STDySsQY2/Lv6rPdFToHGUI9ZYOTYW1GOWkt1EAao9BzdsoJVwmTON6QnOBKy/9RxlhWP+XSWVsY0te6KYzS7rQyzQoJQeeBNMpUnjiQji9kKi5j9rbVMdjIb4HlmYrcE95ps+oFkyJoA1HLVytAeOjJPXGToNlv3k2UPJzOFUM0ujWWeBTyHMCmZ4RhlrfzDNffY5dlW82USjc5dBlzRyZalXSjhcVhK4asUodomVntrvCShp/8C9LpbQZ+ugFNE8J6neStWrhpRU9/sBJx"
export AUTH_JWT_PRIVATE_KEY="MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCjQ2paojdNKLW44G8JSq5yVajgtyxhPbW01RdDxjKpVhFKyIG/5kpkEKMPxeSUs+KHjSQVJ76voM1xrn8sf+iRAh1zyexCD4EsKRd4fY0W/5Bi09zMDDhiFLoZXEmTe72x+nw1q9YFxBsPcwwOl6Ew6+FKXTOg/9mUfoSsoc22hCBSk81VAbKw1TbUqkKWYgMP8vV2uxur6sWIOtxVZo4vAbGvxIF3/RvieC1zAfXpElrFwhddPkBWJrUI13OZTn8AiDc5+DO8qfkAasxkGyYYCILKfj8/c+e4XzX9Ye6STkRMufM4f68LZ6VKZnzLzQcIZ6o6JDWuzrov3frBfxwPAgMBAAECggEARNLHg7t8SoeNy4i45hbYYRRhI5G0IK3t6nQl4YkslBvXIEpT//xpgbNNufl3OYR3SyMhgdWGWe0Ujga8T5sABBj7J3OIp/R3RJFx9nYewwIq8K5VFqNUJWyNYuF3lreEKQHp2Io+p6GasrGR9JjQ95mIGFwfxo/0Pdfzv/5ZhMWTmSTcOi504Vger5TaPobPFOnULq4y1A4eX4puiHDtvx09DUAWbAjGHpCYZjDGRdSXQArYQmUOKy7R46qKT/ollGOWivnEOgsFmXuUWs/shmcrDG4cGBkRrkxyIZhpnpNEEF5TYgulMMzwM+314e8W0lj9iiSB2nXzt8JhEwTz8QKBgQDSCouFj2lNSJDg+kz70eWBF9SQLrBTZ8JcMte3Q+CjCL1FpSVYYBRzwJNvWFyNNv7kHhYefqfcxUVSUnQ1eZIqTXtm9BsLXnTY+uEkV92spjVmfzBKZvtN3zzip97sfMT9qeyagHEHwpP+KaR0nyffAK+VPhlwNMKgQ9rzP4je+QKBgQDG/JwVaL2b53vi9CNh2XI8KNUd6rx6NGC6YTZ/xKVIgczGKTVex/w1DRWFTb0neUsdus5ITqaxQJtJDw/pOwoIag7Q0ttlLNpYsurx3mgMxpYY12/wurvp1NoU3Dq6ob7igfowP+ahUBchRwt1tlezn3TYxVoZpu9dZHtoynOtRwKBgB9vFJJYdBns0kHZM8w8DWzUdCtf0WOqE5xYv4/dyLCdjjXuETi4qFbqayYuwysfH+Zj2kuWCOkxXL6FOH8IQqeyENXHkoSRDkuqwCcAP1ynQzajskZwQwvUbPg+x039Hj4YQCCfOEtBA4T2Fnadmwn0wFJFiOkR/E6f2RSuXX2BAoGALvVqODsxk9s7B0IqH2tbZAsW0CqXNBesRA+w9tIHV2caViFfcPCs+jAORhkkbG5ZZbix+apl+CqQ+trNHHNMWNP+jxVTpTrChHAktdOQpoMu5MnipuLKedI7bPTT/zsweu/FhSFvYd4utzG26J6Rb9hPkOBx9N/KWTXfUcmFJv0CgYAUYVUvPe7MHSd5m8MulxRnVirWzUIUL9Pf1RKWOUq7Ue4oMxzE8CZCJstunCPWgyyxYXgj480PdIuL92eTR+LyaUESb6szZQTxaJfu0mEJS0KYWlONz+jKM4oC06dgJcCMvhgjta2KpXCm3qL1pmKwfFbOLWYBe5uMoHIn9FdJFQ=="
```

and finally we indicate to our API to use an API to authenticate the users:
```bash
export AUTH_API_KEY="with-key"
```

**Generate a token**

Run the following command:
```bash
orchestrate utils generate-jwt --tenant {TENANT_ID} --expiration ${TIME_IN_HOURS}h
```

