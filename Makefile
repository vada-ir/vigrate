#@IgnoreInspection BashAddShebang
export CGO_ENABLED=0
export ENV=development

export APP=migrate
export LDFLAGS="-w -s -X main.buildTime=`date -u +%Y/%m/%d_%H:%M:%S`"

build:
	go get -v ./...
	go build -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) *.go

install:
	go get -v ./...
	go install -v -a -installsuffix cgo -ldflags $(LDFLAGS) *.go

lint:
	go get -v ./...
#	gofmt -s - Checks if the code is properly formatted and could not be further simplified.
	gofmt -s -w .
	go get -v github.com/GeertJohan/fgt
#	go vet - Reports potential errors that otherwise compile.
#	go get -v golang.org/x/tools/cmd/vet
	find . -type d | grep -v '\.\/\.' | xargs fgt go vet
#	go vet --shadow - Reports variables that may have been unintentionally shadowed.
#	find . -type d | grep -v '\.\/\.' | xargs fgt vet --shadow
#	gotype - Syntactic and semantic analysis similar to the Go compiler.
	go get -v golang.org/x/tools/cmd/gotype
	gotype .
#	deadcode - Finds unused code.
	go get -v github.com/tsenart/deadcode
	fgt deadcode .
#	gocyclo - Computes the cyclomatic complexity of functions.
	go get -v github.com/fzipp/gocyclo
	find ./*.go -type f | xargs fgt gocyclo -over 15
#	golint - Google's (mostly stylistic) linter.
	go get -v github.com/golang/lint/golint
	find ./*.go -type f | xargs fgt golint
#	varcheck - Find unused global variables and constants.
	go get -v github.com/opennota/check/cmd/varcheck
	varcheck ./...
#	structcheck - Find unused struct fields.
	go get -v github.com/opennota/check/cmd/structcheck
	structcheck ./...
#	aligncheck - Warn about un-optimally aligned structures.
	go get -v github.com/opennota/check/cmd/aligncheck
	aligncheck ./...
#	errcheck - Check that error return values are used.
	go get -v github.com/kisielk/errcheck
	find . -type d | grep -v '\.\/\.' | xargs fgt errcheck
#	dupl - Reports potentially duplicated code.
	go get -v github.com/mibk/dupl
	find ./*.go -type f | xargs fgt dupl -t 100 -plumbing
#	ineffassign - Detect when assignments to existing variables are not used.
	go get -v github.com/gordonklaus/ineffassign
	ineffassign .
#	interfacer - Suggest narrower interfaces that can be used.
	go get -v github.com/mvdan/interfacer/cmd/interfacer
	find . -type d | grep -v '\.\/\.' | xargs fgt interfacer
#	unconvert - Detect redundant type conversions.
	go get -v github.com/mdempsky/unconvert
	unconvert -v .
#	goconst - Finds repeated strings that could be replaced by a constant.
	go get -v github.com/jgautheron/goconst/cmd/goconst
	goconst ./...
#	gosimple - Report simplifications in code.
	go get -v honnef.co/go/simple/cmd/gosimple
	gosimple ./...
#	staticcheck - Check inputs to functions for correctness
	go get -v honnef.co/go/staticcheck/cmd/staticcheck
	staticcheck ./...
#	goimports - Checks missing or unreferenced package imports.
	go get -v golang.org/x/tools/cmd/goimports
	goimports -w .
#	lll - Report long lines (see --line-length=N).
	go get -v github.com/walle/lll/...
	find ./*.go -type f | xargs fgt lll --maxlength 120
#	misspell - Finds commonly misspelled English words.
	go get -v github.com/client9/misspell/cmd/misspell
	find ./*.go -type f | xargs misspell -error
#	unused - Find unused variables.
	go get -v honnef.co/go/unused/cmd/unused
	unused ./...
