APP_NAME := xm-msa-asset
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call omo.msa.asset AssetService.GetByOwner '{"owner":"hzz", "uid":"5f1022fb6b52c6d205aa8e16"}'

.PHONY: tester
tester:
	go build -o ./bin/ ./tester

.PHONY: dist
dist:
	mkdir -p dist
	rm -f dist/${APP_NAME}-${BUILD_VERSION}.tar.gz
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo.msa.asset:latest

.PHONY: updev
updev:
	scp -P 9700 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@192.168.1.10:/root/

.PHONY: upload2
upload2:
	scp -P 9099 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@47.93.209.105:/root/

.PHONY: upload
upload:
	scp -P 9099 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@101.200.166.80:/root/
