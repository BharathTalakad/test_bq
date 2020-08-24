.PHONY: run-emu-spanner
run-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down
	export SPANNER_EMULATOR_HOST=localhost:9010
	@docker-compose -f spandb/docker-compose-spanner.yaml up -d
	echo "Created spanner emulator"

.PHONY : stop-emu-spanner
stop-emu-spanner:
	@docker-compose -f spandb/docker-compose-spanner.yaml down
