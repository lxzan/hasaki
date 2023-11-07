test:
	go test -count 1 -timeout 30s -run ^Test ./...

bench:
	go test -benchmem -run=^$$ -bench . github.com/lxzan/hasaki

cover:
	go test -coverprofile=bin/cover.out --cover ./...
