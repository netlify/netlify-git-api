set -e

mkdir -p /cmd/build/osx

if [ ! -f /usr/local/bin/clang ]; then
  ln -s $(which clang-3.6) /usr/local/bin/clang
  ln -s $(which clang++-3.6) /usr/local/bin/clang++
fi
#
if [ ! -f /osxcross/tarballs/MacOSX10.10.sdk.tar.xz ]; then
  rm -f /osxcross/tarballs/MacOSX10.6.sdk.tar.bz2
  cp /osx-sdk/MacOSX10.10.sdk.tar.xz /osxcross/tarballs/
  cd /osxcross
  ./build.sh
fi

export MACOSX_DEPLOYMENT_TARGET=10.10

# CMake will blow up if there's no /Applications folder
mkdir -p /Applications

# Need libssh2 for OS X
osxcross-macports install -s libssh2 libgcrypt

cd /deps/libgit2
mkdir -p build/osx
cd build/osx
set +e
OSXCROSS_MP_INC=1 VERBOSE=1 cmake ../.. -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release -DBUILD_CLAR=OFF -DTHREADSAFE=ON \
            -DCMAKE_TOOLCHAIN_FILE=/cmd/toolchains/osx.cmake -DCMAKE_INSTALL_PREFIX=/usr/local/osx \
            -DCMAKE_VERBOSE_MAKEFILE=ON
set -e
# For some reason cmake finds v 0.0 of libssl the first time and errors out, works second time
OSXCROSS_MP_INC=1 VERBOSE=1 cmake ../.. -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release -DBUILD_CLAR=OFF -DTHREADSAFE=ON \
            -DCMAKE_TOOLCHAIN_FILE=/cmd/toolchains/osx.cmake -DCMAKE_INSTALL_PREFIX=/usr/local/osx \
            -DCMAKE_VERBOSE_MAKEFILE=ON
OSXCROSS_MP_INC=1 cmake --build . --target install


# Switch to Go1.4 or OS X because of issues combining dwarfs (did I just write that?!)
# See: https://groups.google.com/forum/#!msg/golang-codereviews/ZBP6jU-v0aQ/q0DYDWHndb0J
GO15_PATH=$PATH
export PATH=/osxcross/target/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# For some reason the osxcross pkg-config messes up the directories. Ideally we would just use the following line:
# FLAGS=$(/osxcross/target/bin/x86_64h-apple-darwin14-pkg-config --static --libs /usr/local/osx/lib/pkgconfig/libgit2.pc)
# But instead we manually fix the output
FLAGS="-Wl,-headerpad_max_install_names -arch x86_64 -L/osxcross/target/macports/pkgs/opt/local/lib -L/usr/local/osx/lib -L/osxcross/target/macports/pkgs/opt/local/lib -liconv -lgit2 -lssh2 -lssl -lcrypto -lz"

export CGO_CFLAGS="-I/usr/local/osx/include"
export CGO_LDFLAGS="/usr/local/osx/lib/libgit2.a -L/usr/local/osx/lib ${FLAGS}"
export PKG_CONFIG_PATH=/usr/local/osx/lib/pkgconfig:/osxcross/target/macports/pkgs/opt/local/lib/pkgconfig

CC=o64-clang GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go get -d github.com/netlify/netlify-git-api
cd /go/src/github.com/netlify/netlify-git-api
CC=o64-clang GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -v -a -i -x -tags netgo -o /cmd/build/osx/netlify-git-api
export PATH=$PATH
