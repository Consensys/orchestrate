module gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git

go 1.13

require (
	github.com/ConsenSys/golang-utils v0.0.0-20190722185538-95555d181804
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/DATA-DOG/godog v0.7.13
	github.com/Shopify/sarama v1.24.1
	github.com/alicebob/gopher-json v0.0.0-20180125190556-5a6b3ba71ee6 // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20191106175434-873d404c7f40 // indirect
	github.com/aws/aws-sdk-go v1.25.45
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/c0va23/go-proxyprotocol v0.9.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible
	github.com/containous/alice v0.0.0-20181107144136-d83ebdd94cbd
	github.com/containous/traefik/v2 v2.0.5
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eapache/channels v1.1.0
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/ethereum/go-ethereum v1.9.7
	github.com/frankban/quicktest v1.5.0 // indirect
	github.com/go-acme/lego/v3 v3.1.0
	github.com/go-pg/migrations v6.7.3+incompatible
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.4
	github.com/hashicorp/vault v1.2.3
	github.com/hashicorp/vault/api v1.0.5-0.20190909201928-35325e2c3262
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/julien-marchand/healthcheck v0.1.0
	github.com/klauspost/compress v1.9.3 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/nmvalera/striped-mutex v0.1.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.0.0-20191202183732-d1d2010b5bee // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.5.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/vulcand/oxy v1.0.0
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036 // indirect
	go.uber.org/atomic v1.5.1 // indirect
	golang.org/x/crypto v0.0.0-20191202143827-86a70503ff7e
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20191203134012-c197fd4bf371 // indirect
	google.golang.org/genproto v0.0.0-20191203145615-049a07e0debe
	google.golang.org/grpc v1.25.1
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/jcmturner/gokrb5.v7 v7.3.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	mellium.im/sasl v0.2.1 // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.1+incompatible
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190725163905-fa8dd90ceb7b
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20180112153951-65b0cdae8d7f
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
	github.com/rancher/go-rancher-metadata => github.com/containous/go-rancher-metadata v0.0.0-20190402144056-c6a65f8b7a28
)
