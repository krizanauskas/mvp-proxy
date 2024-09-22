PROXY_APP_DIR := cmd/proxy/main.go

include .env
export $(shell sed 's/=.*//' .env)

proxy-run:
	@echo "Running Proxy app.."
	go run ${PROXY_APP_DIR}

proxy-debug:
	@echo "Debugging proxy app.."
	dlv debug cmd/proxy/main.go --headless --listen=:2345

proxy-race:
	@echo "Run with Race Detector"
	go run -race ${PROXY_APP_DIR}