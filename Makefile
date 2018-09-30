OS = $(shell uname)
EXE = email-service
PKG = github.com/gernest/alien github.com/urfave/negroni
build: pkg
	@if [ $(OS) != "Linux" ];then \
		export CGO_ENABLED=0;\
		export GOOS=linux;\
		export GOARCH=amd64;\
		GOPATH=`pwd` go install -ldflags "-extldflags -static" $(EXE); \
	else \
		GOPATH=`pwd` go install -ldflags "-linkmode external -extldflags -static" $(EXE); \
	fi
pkg:
	@for p in $(PKG); do \
		echo "downloading $$p ...";\
		GOPATH=`pwd` go get $$p;\
	done

docker: build
	@if [ $(OS) != "Linux" ];then \
		cp bin/linux_amd64/$(EXE) docker-config/; \
	else \
		cp bin/$(EXE) docker-config/; \
	fi
	@echo "building docker-image $(EXE):0.0.1 ..."
	docker build -t $(EXE):0.0.1 docker-config
clean:
	rm -rf bin pkg docker-config/$(EXE) src/github.com

run: 
	@echo "runing $(EXE)"
	docker-compose -f docker-config/docker-compose.yml up -d
stop: 
	@echo "stopping $(EXE)"
	docker-compose -f docker-config/docker-compose.yml down

help:
	@echo "make           -- build $(EXE)"
	@echo "make clean     -- delete temp files"
	@echo "make docker    -- build docker-images"
	@echo "make run       -- run $(EXE) docker server"
	@echo "make stop      -- stop $(EXE) docker server"
