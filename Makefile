
run:
	go run ./cmd/rendercv $(filter-out $@,$(MAKECMDGOALS))

%:
    @:
