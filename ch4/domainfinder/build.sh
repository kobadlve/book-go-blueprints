#!/bin/bash

echo "building domainfinder..."
go build -o domainfinder
echo "building synonyms..."
cd ../synonyms
go build -o ../domainfinder/lib/synonyms
echo "building available..."
cd ../available
go build -o ../domainfinder/lib/available
echo "building sprinkle..."
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle
echo "building coolify..."
cd ../coolify
go build -o ../domainfinder/lib/coolify
echo "building daminify..."
cd ../domainify
go build -o ../domainfinder/lib/domainify
echo "Done"
