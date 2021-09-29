GO111MODULE=on go get github.com/kubernetes/code-generator

(
  # To support running this script from anywhere, we have to first cd into this directory
  # so we can install the tools.
  cd $GOPATH/src/github.com/kubernetes/code-generator
  go install ./cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}
)

$GOPATH/bin/deepcopy-gen --input-dirs github.com/consensys/orchestrate/pkg/http/config/dynamic -O deepcopy --output-package pkg/http/config -h $PWD/scripts/deepcopy/boilerplate.go.tmpl -o $PWD
