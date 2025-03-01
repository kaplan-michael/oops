
VERSION := $(shell git describe --tags --dirty 2>/dev/null)
ifeq ($(VERSION),)
    VERSION := dirty
endif
IMAGE_NAME ?= oops
REPOSITORY ?= quay.io/mkaplan
IMAGE=$(REPOSITORY)/$(IMAGE_NAME)

build-container:
	@echo "Building container $(IMAGE) with tag $(VERSION)"
ifneq ($(findstring dirty,$(VERSION)),)
	buildah bud -t $(IMAGE):$(VERSION)
	buildah push $(IMAGE):$(VERSION)
else
	buildah bud -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .
	buildah push $(IMAGE):$(VERSION) $(IMAGE):latest
endif
