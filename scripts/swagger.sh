# Chain registry docs generation
swag init --dir ./services/chain-registry/api --generalInfo exported.go --output ./services/chain-registry/api/docs
mv ./services/chain-registry/api/docs/swagger.json ./public/swagger-specs/types/chain-registry/swagger.json
rm -r ./services/chain-registry/api/docs