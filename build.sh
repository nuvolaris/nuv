#!/bin/bash
go build -o bin/nuv
cd bin
cp nuv nuv.exe
./nuv build
