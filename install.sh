#!/usr/bin/env bash

go build ./cmd/fastgrep 

if [[ -f /usr/local/bin/fastgrep ]]; then
    echo "command already installed"
    echo "do you want to replace existing command with the current binary? (y/n)"]
    read ans
    if [[ "$ans" = "y" ]]; then
        sudo rm /usr/local/bin/fastgrep
        sudo mv fastgrep /usr/local/bin/fastgrep
    fi

else 
    sudo mv fastgrep /usr/local/bin/fastgrep
fi
