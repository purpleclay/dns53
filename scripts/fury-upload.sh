#!/bin/bash

# Copyright (c) 2022 Purple Clay
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# in the Software without restriction, including without limitation the rights
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/fury-upload.sh
set -e
if [ "${1: -4}" == ".deb" ] || [ "${1: -4}" == ".rpm" ]; then
	cd dist
	echo "uploading $1"
	status="$(curl -s -q -o /dev/null -w "%{http_code}" -F package="@$1" "https://${FURY_TOKEN}@push.fury.io/purpleclay/")"
	echo "got: $status"
	if [ "$status" == "200" ] || [ "$status" == "409" ]; then
		exit 0
	fi
	exit 1
fi