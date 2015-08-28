set -e

mkdir -p /cmd/build/linux

go get -d github.com/netlify/netlify-git-api

cd /deps/libgit2
mkdir -p build/linux
cd build/linux
cmake ../.. -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release -DBUILD_CLAR=OFF -DTHREADSAFE=ON
cmake --build . --target install

FLAGS=$(pkg-config --static --libs /deps/libgit2/build/linux/libgit2.pc)
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="/usr/local/lib/libgit2.a -L/usr/local/lib ${FLAGS}"
cd /go/src/github.com/netlify/netlify-git-api
go build -v -a -i -x -tags netgo --ldflags='-extldflags "-lgpg-error -static"' -o /cmd/build/linux/netlify-git-api
