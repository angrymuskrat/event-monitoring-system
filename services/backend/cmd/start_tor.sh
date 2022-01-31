#!/bin/bash

socks_port=9150
control_port=8218

# Create data directory if it doesn't exist
if [ ! -d "tor_data" ]; then
    mkdir "tor_data"
fi

if [ ! -d "tor_data/tor$i" ]; then
    echo "Creating directory tor_data/tor$i"
    mkdir "tor_data/tor$i"
fi
# Take into account that authentication for the control port is disabled. Must be used in secure and controlled environments

echo "Running: tor --RunAsDaemon 1 --CookieAuthentication 0 --HashedControlPassword \"\" --ControlPort $control_port --PidFile tor$i.pid --SocksPort $socks_port --DataDirectory tor_data/tor$i"

tor --RunAsDaemon 1 --CookieAuthentication 0 --HashedControlPassword "" --ControlPort $control_port --PidFile /home/muskrat/workplace/tor/tort$i.pid --SocksPort $socks_port --DataDirectory /home/muskrat/workplace/tor/tort$i