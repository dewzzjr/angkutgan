#!/usr/bin/env bash

build_package () {
    package=$1
    package_split=(${package//\// })
    package_name="../bin/${package_split[-1]}"
    platforms=("windows/amd64" "linux/amd64")

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

rm -rf bin/*
cd scripts
build_package "github.com/dewzzjr/angkutgan/cmd"