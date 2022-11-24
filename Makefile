export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=linux:amd64

all: build

build: app

app:
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -ldflags "$(LDFLAGS)" -o ./release/HustWebAuth_$${target_suffix} ./;\
		echo "Build $${os}-$${arch} done";\
	)
