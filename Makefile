run: ## run the application
	go run -ldflags " \
	      -X github.com/sjeandeaux/azure-ad-go/information.Version=$(shell cat VERSION) \
	      -X github.com/sjeandeaux/azure-ad-go/information.BuildTime=$(shell date +"%Y-%m-%dT%H:%M:%S") \
	      -X github.com/sjeandeaux/azure-ad-go/information.GitCommit=$(shell git rev-parse --short HEAD) \
	      -X github.com/sjeandeaux/azure-ad-go/information.GitDescribe=$(shell git describe --tags --always) \
	      -X github.com/sjeandeaux/azure-ad-go/information.GitDirty=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)" \
				main.go


help: ## this help
	@grep -hE '^[a-zA-Z_-]+.*?:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
