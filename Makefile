GO         := $(shell command -v go 2>/dev/null || echo /usr/local/go/bin/go)
GO_VERSION := 1.25.0
INSTALL_DIR := $(HOME)/.local/bin
BINARY     := hb

.PHONY: all install build go-install

all: install

install: go-ensure build
	@mkdir -p $(INSTALL_DIR)
	@cp habits $(INSTALL_DIR)/$(BINARY)
	@echo "Instalado em $(INSTALL_DIR)/$(BINARY)"
	@grep -qF "$(INSTALL_DIR)" $(HOME)/.bashrc 2>/dev/null || \
		echo 'export PATH=$$PATH:$(INSTALL_DIR)' >> $(HOME)/.bashrc
	@grep -qF "$(INSTALL_DIR)" $(HOME)/.zshrc 2>/dev/null || \
		echo 'export PATH=$$PATH:$(INSTALL_DIR)' >> $(HOME)/.zshrc 2>/dev/null || true
	@echo "PATH atualizado em ~/.bashrc. Rode: source ~/.bashrc"

build: go-ensure
	@$(GO) build -o habits ./cmd/habits
	@echo "Build concluído: ./habits"

go-ensure:
	@if ! command -v go >/dev/null 2>&1 && [ ! -x /usr/local/go/bin/go ]; then \
		$(MAKE) go-install; \
	fi

go-install:
	$(eval OS   := $(shell uname -s | tr A-Z a-z))
	$(eval ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'))
	$(eval PKG  := go$(GO_VERSION).$(OS)-$(ARCH).tar.gz)
	@echo "Instalando Go $(GO_VERSION)..."
	@curl -fsSL "https://go.dev/dl/$(PKG)" -o /tmp/$(PKG)
	@sudo tar -C /usr/local -xzf /tmp/$(PKG)
	@rm /tmp/$(PKG)
	@echo "Go instalado em /usr/local/go"
