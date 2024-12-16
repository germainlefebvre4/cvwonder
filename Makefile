
run:
	go run ./cmd/cvrender $(filter-out $@,$(MAKECMDGOALS))

%:
    @:
