#!/usr/bin/env bash

if go build
then
    echo "Build ok"
    ./insights-operator-web-ui
else
    echo "Build failed"
fi
