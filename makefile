GOFILE = main.go
DIR 		= ./build
NAME  	= library-manager


.PHONY: all
all: linux-arm64 linux-arm linux-riscv64 linux-386 linux-amd64 darwin-amd64 darwin-arm64 windows-386 windows-amd64 windows-arm64

.PHONY: clean
clean:
	rm -rf $(DIR)/*

.PHONY: build
build:
	go build -o $(DIR)/$(NAME) $(GOFILE)

.PHONY: linux-arm64
linux-arm64:
	$(eval GOOS := linux)
	$(eval GOARCH := arm64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: linux-arm
linux-arm:
	$(eval GOOS := linux)
	$(eval GOARCH := arm)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: linux-riscv64
linux-riscv64:
	$(eval GOOS := linux)
	$(eval GOARCH := riscv64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: linux-386
linux-386:
	$(eval GOOS := linux)
	$(eval GOARCH := 386)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: linux-amd64
linux-amd64:
	$(eval GOOS := linux)
	$(eval GOARCH := amd64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: darwin-amd64
darwin-amd64:
	$(eval GOOS := darwin)
	$(eval GOARCH := amd64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: darwin-arm64
darwin-arm64:
	$(eval GOOS := darwin)
	$(eval GOARCH := arm64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH) $(GOFILE)

.PHONY: windows-386
windows-386:
	$(eval GOOS := windows)
	$(eval GOARCH := 386)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH).exe $(GOFILE)

.PHONY: windows-amd64
windows-amd64:
	$(eval GOOS := windows)
	$(eval GOARCH := amd64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH).exe $(GOFILE)

.PHONY: windows-arm64
windows-arm64:
	$(eval GOOS := windows)
	$(eval GOARCH := arm64)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(DIR)/$(NAME)_$(GOOS)_$(GOARCH).exe $(GOFILE)
