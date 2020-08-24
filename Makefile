.PHONY: start-emu-spanner
start-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down
	export SPANNER_EMULATOR_HOST=localhost:9010
	@docker-compose -f spandb/docker-compose-spanner.yaml up -d
	@echo "Created spanner emulator"
	gcloud spanner instances create testIns --config=emulator-config --description="Test Instance" --nodes=1


.PHONY : stop-emu-spanner
stop-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down

.PHONY: run-all
run-all:
	@go run main.go