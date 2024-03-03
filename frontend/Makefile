NODE_PACKAGE_MANAGER = yarn

regit:
	(rm -rf .git && cd ../ && rm -rf frontend2 && git clone git@github.com:dinnerdonebetter/frontend frontend2 && cp -rf frontend2/.git frontend/.git && rm -rf frontend2)

clean:
	rm -rf node_modules \
	apps/web/node_modules apps/web/.turbo apps/web/.next \
	packages/eslint-custom-config/node_modules \
	packages/models/node_modules packages/models/.turbo

.PHONY: install
install:
	$(NODE_PACKAGE_MANAGER) install

.PHONY: lint
lint:
	$(NODE_PACKAGE_MANAGER) lint

.PHONY: terraformat
terraformat:
	@(cd environments/dev/webapp/terraform && terraform fmt)
	@(cd environments/dev/landing/terraform && terraform fmt)
	@(cd environments/dev/admin/terraform && terraform fmt)

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd environments/dev/admin/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd environments/dev/landing/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)
	@(cd environments/dev/webapp/terraform && terraform init -upgrade && terraform validate && terraform fmt && terraform fmt -check)

.PHONY: dev
dev:
	$(NODE_PACKAGE_MANAGER) dev

.PHONY: stop
stop:
	kill -KILL `lsof -i tcp:9000 -i tcp:7000 | tail -n +2 | awk '{print $$2}'`

.PHONY: build
build:
	$(NODE_PACKAGE_MANAGER) build

.PHONY: test
test:
	$(NODE_PACKAGE_MANAGER) test

.PHONY: format
format:
	$(NODE_PACKAGE_MANAGER) format

.PHONY: format-check
format-check:
	$(NODE_PACKAGE_MANAGER) format-check
