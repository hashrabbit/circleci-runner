# Copyright 2016-2017 HashRabbit, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO ?= go
GOX ?= gox

VERSION = $(shell git describe --tags 2>/dev/null || echo "0.0.0-dev")

GO_LDFLAGS = -X main.Version=$(VERSION)

# List of os/arch targets to build binaries for.
GOX_OSARCH = \
	darwin/386 \
	darwin/amd64 \
	linux/386 \
	linux/amd64 \
	linux/arm \
	windows/386 \
	windows/amd64

.PHONY: build
build:
	$(GOX) -ldflags="$(GO_LDFLAGS)" -osarch="$(GOX_OSARCH)"

.PHONY: test
	$(GO) test ./...

.PHONY: clean
clean:
	$(RM) circleci-runner*
