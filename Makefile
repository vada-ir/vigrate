#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED= 0
export GOOS=linux
export ENV=development

export GLIDE_HOME=$(HOME)/.glide

export APP=migrate
export LDFLAGS="-w -s"

export DEBUG= 1

all: lint build

fetch: glide-install

contributors:
	git log --all --format='%aN <%cE>' | sort -u  > CONTRIBUTORS

#######
# Build
#######

build: fetch
	go build -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) *.go

install: fetch
	go install -v -a -installsuffix cgo -ldflags $(LDFLAGS) *.go

run:
	go run *.go

######
# Lint
######

check-gometalinter:
	which gometalinter || (go get -u -v github.com/alecthomas/gometalinter && gometalinter --install)

lint: fetch check-gometalinter
	gometalinter \
	--vendor --skip=vendor/ --exclude=vendor \
	--disable-all \
	--enable=gofmt \
	--enable=vet --enable=vetshadow \
	--enable=gocyclo \
	--cyclo-over=128 \
	--enable=golint \
	--enable=ineffassign \
	--enable=misspell \
	--concurrency=1 \
	--deadline=5m \
	./...

format:
	which gometalinter || go get -u -v golang.org/x/tools/cmd/goimports
	find $(ROOT)/ -type f -name "*.go" | grep -v $(ROOT)/vendor | xargs --max-args=1 --replace=R goimports -w R
	find $(ROOT)/ -type f -name "*.go" | grep -v $(ROOT)/vendor | xargs --max-args=1 --replace=R gofmt -s -w R

#######
# Glide
#######

check-glide: check-glide
	which glide || curl https://glide.sh/get | sh

check-glide-init:
	@[ -f $(ROOT)/glide.yaml ] || make -f $(ROOT)/Makefile glide-init

# Scan a codebase and create a glide.yaml file containing the dependencies.
glide-init: check-glide
	glide init

# Install the latest dependencies into the vendor directory matching the version resolution information.
# The complete dependency tree is installed, importing Glide, Godep, GB, and GPM configuration along the way.
# A lock file is created from the final output.
glide-update: check-glide check-glide-init
	glide update

# Install the dependencies and revisions listed in the lock file into the vendor directory.
# If no lock file exists an update is run.
glide-install: check-glide check-glide-init
	glide install


