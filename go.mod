module github.com/consensys/orchestrate

go 1.16

require (
	github.com/Shopify/sarama v1.27.2
	github.com/alicebob/gopher-json v0.0.0-20180125190556-5a6b3ba71ee6 // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/c0va23/go-proxyprotocol v0.9.1
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/consensys/gnark-crypto v0.5.0
	github.com/consensys/quorum v2.7.0+incompatible
	github.com/consensys/quorum-key-manager v0.0.0-20210922104442-c2cd4d2be057
	github.com/containous/alice v0.0.0-20181107144136-d83ebdd94cbd
	github.com/containous/traefik/v2 v2.2.0
	github.com/cucumber/godog v0.11.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dgraph-io/ristretto v0.0.3
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/eapache/channels v1.1.0
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/ethereum/go-ethereum v1.10.8
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-pg/migrations/v7 v7.1.9
	github.com/go-pg/pg/v9 v9.1.5
	github.com/go-playground/validator/v10 v10.5.0
	github.com/go-test/deep v1.0.2 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.0
	github.com/gomodule/redigo v1.8.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/hashicorp/serf v0.8.3 // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20200117231345-460d63e36490
	github.com/hashicorp/vault/sdk v0.1.14-0.20200305172021-03a3749f220d // indirect
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/justinas/alice v1.2.0
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/mapstructure v1.3.2
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/nmvalera/striped-mutex v0.1.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/fosite v0.34.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/unrolled/secure v1.0.7
	github.com/vulcand/oxy v1.1.0
	github.com/yuin/gopher-lua v0.0.0-20191220021717-ab39c6098bdb // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	google.golang.org/protobuf v1.26.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203
	google.golang.org/api => google.golang.org/api v0.10.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/cloudflare/cloudflare-go => github.com/cloudflare/cloudflare-go v0.10.2
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
	google.golang.org/grpc => google.golang.org/grpc v1.28.0 // indirect
)
