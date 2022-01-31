#!/bin/bash

base_socks_port=9150
base_control_port=8218

# Create data directory if it doesn't exist
if [ ! -d "/home/muskrat/tor/" ]; then
    mkdir -p "/home/muskrat/tor/"
fi

#for i in {0..10}
for i in {0..20}

do
    socks_port=$((base_socks_port+i))
    control_port=$((base_control_port+i))
    if [ ! -d "/home/muskrat/tor/tor$i" ]; then
        echo "Creating directory data/tor$i"
        mkdir "/home/muskrat/tor/tor$i"
    fi
    # Take into account that authentication for the control port is disabled. Must be used in secure and controlled environments

    echo "Running: tor --RunAsDaemon 1 --CookieAuthentication 0 --HashedControlPassword \"\" --ControlPort $control_port --PidFile tor$i.pid --SocksPort $socks_port --DataDirectory /home/muskrat/tor/tor$i"

    tor --RunAsDaemon 1 --CookieAuthentication 0 --HashedControlPassword "" --ControlPort $control_port --PidFile /home/muskrat/tor/tort$i.pid --SocksPort $socks_port --DataDirectory /home/muskrat/tor/tort$i
done