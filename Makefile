build:
	cd swagger-ui && rm -rf node_modules && rm -rf dist && npm install && mv node_modules/swagger-ui-dist dist && rm -r node_modules

build-ci: build
	cd swagger-ui && rm -rf node_modules
