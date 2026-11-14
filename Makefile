# maulid.id monolith build and verification entry points.
#
# Go 1.27.4+ and Node.js 22+ are required. Override GO or NPM when the tools are
# not on PATH, for example:
#   make test GO=/home/me/.local/share/maulid-toolchains/go1.27.4/go/bin/go

GO ?= go
NPM ?= npm
WEB := web
BIN := bin

.PHONY: all deps build web-build embed test test-go test-web coverage coverage-go coverage-web \
	lint lint-go lint-web vet race security run run-dev staff-bootstrap clean

all: build

## Install locked dependencies.
deps:
	$(GO) mod download
	cd $(WEB) && $(NPM) ci

## Build the single production artefact (binary + embedded React bundle).
build: web-build embed
	mkdir -p $(BIN)
	$(GO) build -trimpath -o $(BIN)/maulid ./cmd/server
	$(GO) build -trimpath -o $(BIN)/staff-bootstrap ./cmd/staff-bootstrap

## Build the React public site and backoffice.
web-build:
	cd $(WEB) && $(NPM) run build

## Copy the React build into the Go embed directory.
embed:
	rm -rf internal/webui/dist
	mkdir -p internal/webui/dist
	cp -R $(WEB)/dist/. internal/webui/dist/
	touch internal/webui/dist/.keep

## All test suites.
test: test-go test-web

test-go:
	$(GO) test ./...

test-web:
	cd $(WEB) && $(NPM) test

## Coverage gates: backend statements and frontend metrics must exceed 80%.
coverage: coverage-go coverage-web

coverage-go:
	./scripts/check-go-coverage.sh

coverage-web:
	cd $(WEB) && $(NPM) run coverage

## Static checks.
lint: lint-go lint-web

lint-go:
	@unformatted="$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi
	$(GO) vet ./...

lint-web:
	cd $(WEB) && $(NPM) run lint

vet:
	$(GO) vet ./...

race:
	$(GO) test -race ./...

## Dependency vulnerabilities plus the frontend dependency audit.
security:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
	cd $(WEB) && $(NPM) run audit

## Run the monolith on http://127.0.0.1:8080 with local data in ./data.
run:
	MAULID_ENV=development MAULID_DB_PATH=$${MAULID_DB_PATH:-./data/maulid.db} \
		$(GO) run ./cmd/server

## Frontend dev server on :5173 proxying /api to the Go server on :8080.
run-dev:
	cd $(WEB) && $(NPM) run dev

## Create the first staff account through the operator-only path.
staff-bootstrap:
	MAULID_ENV=development MAULID_DB_PATH=$${MAULID_DB_PATH:-./data/maulid.db} \
		$(GO) run ./cmd/staff-bootstrap --email "$${EMAIL:?set EMAIL=...}" --name "$${NAME:?set NAME=...}"

clean:
	rm -rf $(BIN) $(WEB)/dist $(WEB)/coverage data
