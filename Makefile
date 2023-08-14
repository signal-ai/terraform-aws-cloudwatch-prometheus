.PHONY: LOCAL_BIN
LOCAL_BIN := $(CURDIR)/bin
$(LOCAL_BIN):
	mkdir -p $@
.PHONY: LOCAL_TMP
LOCAL_TMP := $(CURDIR)/tmp
$(LOCAL_TMP):
	mkdir -p $@

SYSTEM := $(shell uname -s | tr A-Z a-z)
ARCH := $(shell uname -m)
ifeq ($(ARCH), x86_64)
	ARCH := amd64
endif

TERRAFORM_DOCS_VERSION ?= v0.16.0
.PHONY: terraform-docs-cmd
terraform-docs-cmd = $(LOCAL_BIN)/terraform-docs
$(terraform-docs-cmd): | $(LOCAL_TMP) $(LOCAL_BIN)
	curl -sSL -o $(LOCAL_TMP)/terraform-docs.tar.gz https://github.com/terraform-docs/terraform-docs/releases/download/$(TERRAFORM_DOCS_VERSION)/terraform-docs-$(TERRAFORM_DOCS_VERSION)-$(SYSTEM)-$(ARCH).tar.gz
	tar -C $(LOCAL_TMP) -xzf $(LOCAL_TMP)/terraform-docs.tar.gz
	mv $(LOCAL_TMP)/terraform-docs $(terraform-docs-cmd)

docs: | $(terraform-docs-cmd)
	$(terraform-docs-cmd) .

build:
	cd lambda_code && env GOOS=linux GOARCH=amd64 go build -o main
	cd lambda_code && zip -r payload.zip main && rm main
