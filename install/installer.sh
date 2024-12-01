#!/bin/bash

# Check if PostgreSQL is installed
if command -v psql &> /dev/null
then
    echo "PostgreSQL is installed."
    
else
    echo "PostgreSQL is not installed."
    echo "Please install PostgreSQL"
    exit 1
fi

# Get the installed Go version (strip "go" prefix and split into major and minor version)
installed_version=$(go version | awk '{print $3}' | cut -d 'o' -f 2)


# Extract the major and minor versions from the installed version
major_version=$(echo $installed_version | cut -d '.' -f 1)
minor_version=$(echo $installed_version | cut -d '.' -f 2)

# Desired version to compare against (1.22)
required_major=1
required_minor=22

if [[ $major_version -gt $required_major ]]; then
    echo "go version :" $installed_version "OK"
elif [[ $major_version -eq $required_major && $minor_version -ge $required_minor ]]; then
    echo "go version :" $installed_version "OK"
else
    echo "go version:" $installed_version ",required version greater than 1.22"
    exit 1
fi
cd install
go run install.go
cd ..