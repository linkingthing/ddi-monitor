GOSRC = $(shell find . -type f -name '*.go')

build: ddi_monitor

ddi_monitor: $(GOSRC) 
	CGO_ENABLED=0 GOOS=linux go build -o ddi_monitor cmd/main/main.go

clean:
	rm -rf ddi_monitor

.PHONY: clean install
