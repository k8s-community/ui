#!/bin/bash

GO_OS=${GO_OS:-"linux"}

function detect_os {
    # Detect the OS name
    case "$(uname -s)" in
      Darwin)
        host_os=darwin
        ;;
      Linux)
        host_os=linux
        ;;
      *)
        echo "Unsupported host OS. Must be Linux or Mac OS X." >&2
        exit 1
        ;;
    esac

   GO_OS="${host_os}"
}

detect_os

old_version=$(grep RELEASE?= ./Makefile | sed -e 's/RELEASE?=\(.*\)/\1/g')
supposed_version=$(echo $old_version | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')

echo "Current version $old_version."
echo -n "Please enter bumped version [$supposed_version]:"
read new_version
new_version=${new_version:-$supposed_version}

if [[ $new_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Bumped version: $new_version"
else
    echo "Version is incorrect, please use XX.XX.XX format (ie: 4.15.3)"
    exit
fi

if [ "${GO_OS}" == "darwin" ]; then
    sed -i '' -e "s/\(RELEASE?=\).*/\1$new_version/" Makefile
    sed -i '' -e "s/\(version:\)\(\s*\).*/\1 $new_version/" charts/Chart.yaml
    sed -i '' -e "s/\(tag:\)\(\s*\).*/\1 $new_version/" charts/values-stable.yaml
    sed -i '' -e "s/\(tag:\)\(\s*\).*/\1 $new_version/" charts/values-dev.yaml
else
    sed -i -e "s/\(RELEASE?=\).*/\1$new_version/" Makefile
    sed -i -e "s/\(version:\)\(\s*\).*/\1 $new_version/" charts/Chart.yaml
    sed -i -e "s/\(tag:\)\(\s*\).*/\1 $new_version/" charts/values-stable.yaml
    sed -i -e "s/\(tag:\)\(\s*\).*/\1 $new_version/" charts/values-dev.yaml
fi