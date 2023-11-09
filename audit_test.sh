#!/bin/bash

commands=("go run main.go ./data/example00" "go run main.go ./data/example01" "go run main.go ./data/example02" "go run main.go ./data/example03" "go run main.go ./data/example04" "go run main.go ./data/example05" "go run main.go ./data/example06" "go run main.go ./data/example07" "go run main.go ./data/badexample00" "go run main.go ./data/badexample01")

for cmd in "${commands[@]}"
do
    echo "Running: $cmd"
    echo
    eval $cmd
    echo
    read -p "Press enter to continue"
    clear
done
