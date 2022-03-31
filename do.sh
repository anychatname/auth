#!/bin/bash

main() {
	case $1 in
		run)
			run "$2"
			return
			;;
		test)
			test
			return
			;;
		cover)
			cover
			return
			;;
	esac
}

run() {
	echo "Executing $1 run..."
}

test() {
	echo "Executing tests..."
}

cover() {
	echo "Executing coverage..."
}

main "$@"
