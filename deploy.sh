#!/bin/bash

echo "build start"
$(go env GOPATH)/bin/qtdeploy build desktop
echo "built"
./deploy/linux/qt_test