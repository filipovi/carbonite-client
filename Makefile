#!make
# include .env.local
# export $(shell sed 's/=.*//' .env.local)

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

server/build: templ/generate ## Build the server
	@echo 'Build the server...'
	go build -buildvcs=false -o ./bin ./cmd/client

server/run: server/build ## Run the project
	@echo 'Launch the server...'
	./bin/client

server/audit: tidy fmt vet staticcheck test gosec ## audit

tidy: ## tidy dependencies
	@echo 'Tidying and verifying module dependencies...' 
	go mod tidy
	go mod verify

fmt: ## Formatting the code
	@echo 'Formatting code...'
	go fmt ./...

vet: ## Vetting the code
	@echo 'Vetting code...'
	go vet ./...

staticcheck: ## Static check the code
	@echo 'Static check...'
	staticcheck ./...

test: ## Launch the tests
	@echo 'Running tests...'
	go test -race -vet=off ./...

gosec: ## Launch gosec
	@echo 'Security check...'
	gosec ./...

templ/generate: ## generate the templates
	@echo 'Generate the templates...'
	~/go/bin/templ generate

env/test: ## Show the env variables
	env

.PHONY: help audit

HELP_FUNCTION = \
    %help; \
    while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-\/]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
    print "usage: make [target]\n\n"; \
    for (sort keys %help) { \
    for (@{$$help{$$_}}) { \
    $$sep = " " x (32 - length $$_->[0]); \
    print "  \033[0;33m$$_->[0]\033[0;37m$$sep\033[0;32m$$_->[1]\033[0;37m\n"; \
    }; \
    print "\n"; }

.DEFAULT_GOAL := help

help: ##@other Show this help.
	@perl -e '$(HELP_FUNCTION)' $(MAKEFILE_LIST)

ERR = $(error found an error!)

err: ; $(ERR)
