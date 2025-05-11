declare -a platforms=("linux")
declare -a archs=("arm" "arm64" "386" "amd64")

mkdir -p release

for platform in "${platforms[@]}"
do
    for arch in "${archs[@]}"
    do
        echo "Building $platform - $arch"
        env GOOS=$platform GOARCH=$arch go build -o release/sendip-client-$platform-$arch -ldflags "-s -w" client/main.go
        env GOOS=$platform GOARCH=$arch go build -o release/sendip-server-$platform-$arch -ldflags "-s -w" server/main.go
    done
done
