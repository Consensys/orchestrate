set -e

for p in `find . -name *.proto`; do
    echo "# $p"
    protoc -I. --go_out=plugins=grpc+protoc-gen-grpc-gateway+protoc-gen-swagger,paths=source_relative:. $p
done
