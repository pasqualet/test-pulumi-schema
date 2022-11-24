VERSION         := $(shell pulumictl get version)

PACK            := test-pulumi-schema
PROJECT         := github.com/pasqualet/${PACK}

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}

WORKING_DIR     := $(shell pwd)
SCHEMA_PATH     := ${WORKING_DIR}/schema.yaml

gen_nodejs_sdk::
	rm -rf sdk/nodejs
	cd provider/cmd/${CODEGEN} && go run . nodejs ../../../sdk/nodejs ${SCHEMA_PATH}

build_nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
build_nodejs_sdk:: gen_nodejs_sdk
	cd sdk/nodejs/ && \
		npm install && \
		npm run build && \
		cp -R scripts bin/ && \
		cp ../../README.md ../../LICENSE package.json package-lock.json ./bin/ && \
		sed -i.bak -e "s/\$${VERSION}/$(VERSION)/g" ./bin/package.json && \
		rm ./bin/package.json.bak
