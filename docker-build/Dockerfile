# Go cross compiler (xgo): Wildcard layer to the latest Go release
# Copyright (c) 2014 Péter Szilágyi. All rights reserved.
#
# Released under the MIT license.

FROM karalabe/xgo-1.4.x

MAINTAINER Péter Szilágyi <peterke@gmail.com>

RUN apt-get install -y cmake libssh2-1-dev

RUN mkdir -p /deps && cd /deps && wget https://github.com/libgit2/libgit2/archive/v0.22.3.tar.gz && \
    tar zxvf v0.22.3.tar.gz && mv libgit2-0.22.3 libgit2 && rm v0.22.3.tar.gz

RUN git clone https://go.googlesource.com/go /usr/local/go-1.5 && \
    cd /usr/local/go-1.5/src && \
    GOROOT_BOOTSTRAP=/usr/local/go ./make.bash

RUN apt-get install -y clang-3.6 mingw-w64-tools

ENV PATH /usr/local/go-1.5/bin:$PATH

ENTRYPOINT /bin/bash
