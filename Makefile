run-example-app:
	$(MAKE) -C docker run-example-compose-app

test:
	@echo "Running dashboard testing...\n"
	$(MAKE) -C dashboard test
	@echo "Running server unit tests...\n"
	$(MAKE) -C server DASHBOARD_PATH=../dashboard/dist test
generate:
	$(MAKE) -C dashboard generate
	$(MAKE) -C server generate