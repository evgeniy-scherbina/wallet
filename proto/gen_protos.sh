#!/usr/bin/env sh

if [ "$1" != "INSIDE_CONTAINER" ]; then
    echo "Generating root gRPC server protos"

    PB_WALLET_PATH="../pb/wallet"
    DOCS_PATH="../docs"

    mkdir -p ${PB_WALLET_PATH}
    mkdir -p ${DOCS_PATH}

    PROJECT_PATH=$(dirname $(pwd -P))
    WALLET_PATH="/go/src/github.com/evgeniy-scherbina/peach-secure-server"
    docker run --rm \
        -v ${PROJECT_PATH}:${WALLET_PATH} \
        -e WALLET_PATH=${WALLET_PATH} \
        -e PB_WALLET_PATH=${PB_WALLET_PATH} \
        -e DOCS_PATH=${DOCS_PATH} \
        -w ${WALLET_PATH}/proto \
        lightningnetwork/protoc-alpine:20190903.1237.02 \
        ./gen_protos.sh "INSIDE_CONTAINER"
else
    protoc -I. \
        -I${WALLET_PATH}/vendor \
        -I/usr/local/include \
        -I/usr/local/include/third_party/googleapis \
        --go_out=plugins=grpc,paths=source_relative:${PB_WALLET_PATH} \
        --grpc-gateway_out=logtostderr=true,paths=source_relative:${PB_WALLET_PATH} \
        --swagger_out=logtostderr=true:${DOCS_PATH} \
        --govalidators_out=paths=source_relative:${PB_WALLET_PATH} \
        rpc.proto
fi