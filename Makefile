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

GOX ?= gox

VERSION = $(shell git describe --tags 2>/dev/null || echo "0.0.0-dev")

GO_LDFLAGS = -X main.Version=$(VERSION)

.PHONY: build
build:
	$(GOX) -ldflags "$(GO_LDFLAGS)"
