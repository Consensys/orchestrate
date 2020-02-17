# Chain registry docs generation
swag init -d ./services/chain-registry/api -g exported.go -o ./services/chain-registry/api/docs
mv ./services/chain-registry/api/docs/swagger.json ./public/swagger-specs/types/chain-registry/swagger.json
rm -r ./services/chain-registry/api/docs