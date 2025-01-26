.PHONY: generate icons clean

# Default target
all: icons

# Target to generate icons
icons:
	@echo "Generating icons..."
	@mkdir -p scripts
	@if [ ! -d "scripts/nvim-web-devicons" ]; then \
		git clone https://github.com/nvim-tree/nvim-web-devicons.git scripts/nvim-web-devicons; \
	fi
	@node scripts/generate_icons.js

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf scripts/nvim-web-devicons
