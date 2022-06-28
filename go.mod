module github.com/consensys/orchestrate

go 1.16

require (
	github.com/Shopify/sarama v1.27.2
	github.com/alicebob/gopher-json v0.0.0-20180125190556-5a6b3ba71ee6 // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/auth0/go-jwt-middleware/v2 v2.0.0-beta.1
	github.com/c0va23/go-proxyprotocol v0.9.1
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/consensys/quorum v2.7.0+incompatible
	github.com/consensys/quorum-key-manager v0.0.0-20211216092116-3e1719d4b69d
	github.com/cucumber/godog v0.11.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dgraph-io/ristretto v0.1.0
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/eapache/channels v1.1.0
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/ethereum/go-ethereum v1.10.8
	github.com/go-kit/kit v0.10.1-0.20200915143503-439c4d2ed3ea
	github.com/go-pg/migrations/v7 v7.1.9
	github.com/go-pg/pg/v9 v9.1.5
	github.com/go-playground/validator/v10 v10.10.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/justinas/alice v1.2.0
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nmvalera/striped-mutex v0.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.26.0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/traefik/paerser v0.1.4
	github.com/traefik/traefik/v2 v2.5.5
	github.com/umbracle/go-web3 v0.0.0-20220214174057-1d971ac3be78
	github.com/unrolled/secure v1.0.9
	github.com/vulcand/oxy v1.3.0
	github.com/yuin/gopher-lua v0.0.0-20191220021717-ab39c6098bdb // indirect
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6
	google.golang.org/protobuf v1.28.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.11
	github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203
	github.com/spf13/viper => github.com/spf13/viper v1.8.1
	google.golang.org/api => google.golang.org/api v0.10.0
	google.golang.org/grpc => google.golang.org/grpc v1.28.0 // indirect
	hashicorp/consul => hashicorp/consul v1.10.1 // indirect
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
)
