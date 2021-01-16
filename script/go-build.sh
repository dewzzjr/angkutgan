#!/usr/bin/env bash

build_package () {
    package=$1
    package_split=(${package//\// })
    package_name="../executable/${package_split[-1]}"
    platforms=("windows/amd64" "windows/386" "linux/amd64")

    for platform in "${platforms[@]}"
    do
        platform_split=(${platform//\// })
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        output_name=$package_name'-'$GOOS'-'$GOARCH
        if [ $GOOS = "windows" ]; then
            output_name+='.exe'
        fi

        env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
        fi
    done
}

cd script
build_package "github.com/dewzzjr/angkutgan/backend"
build_package "github.com/dewzzjr/angkutgan/fileserver"