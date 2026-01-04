ESBUILD ?= npx esbuild
ESBUILD_OPTIONS ?= --target=es6 --minify
NODE ?= node
MINIFY_HTML ?= npx html-minifier --collapse-whitespace

.PHONY: all build build-challenge build-blocked

all: build

build: build-challenge build-blocked

build-blocked:
	mkdir -p tmp
	$(ESBUILD) --bundle assets/blocked.css --outfile=tmp/blocked.css $(ESBUILD_OPTIONS)

	$(NODE) -e "\
	const fs = require('fs');\
	const html = fs.readFileSync('assets/blocked.html', 'utf8');\
	const css = fs.readFileSync('tmp/blocked.css', 'utf8');\
	const result = html.replace('<!--HEAD-->', '<style>' + css + '</style><!--HEAD-->');\
	fs.writeFileSync('tmp/blocked.html', result);\
	"
	$(MINIFY_HTML) -o tmp/blocked.html tmp/blocked.html
	cp tmp/blocked.html ../internal/blocked.html
	rm -rf tmp

build-challenge:
	mkdir -p tmp
	$(ESBUILD) --bundle assets/challenge.css --outfile=tmp/challenge.css $(ESBUILD_OPTIONS)
	$(ESBUILD) --bundle assets/challenge.js --outfile=tmp/challenge.js $(ESBUILD_OPTIONS)

	$(NODE) -e "\
	const fs = require('fs');\
	const html = fs.readFileSync('assets/challenge.html', 'utf8');\
	const css = fs.readFileSync('tmp/challenge.css', 'utf8');\
	const result = html.replace('<!--HEAD-->', '<style>' + css + '</style><!--HEAD-->');\
	fs.writeFileSync('tmp/challenge.html', result);\
	"
	$(MINIFY_HTML) -o tmp/challenge.html tmp/challenge.html
	cp tmp/challenge.html ../internal/challenge.html
	cp tmp/challenge.js ../internal/challenge.js
	rm -rf tmp
