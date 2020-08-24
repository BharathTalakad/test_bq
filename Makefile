.PHONY: start-emu-spanner
start-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down
	export SPANNER_EMULATOR_HOST=localhost:9010
	@docker-compose -f spandb/docker-compose-spanner.yaml up -d
	@echo "Started spanner emulator"


.PHONY : stop-emu-spanner
stop-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down

.PHONY: run-all
run-all:
	@go run main.go