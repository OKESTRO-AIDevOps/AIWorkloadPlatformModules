#!/bin/bash

echo 'Please enter the 1st arguments(duration): '
read duration

echo 'Please enter the 1st arguments(mode): '
read mode

echo 'Please enter the 1st arguments(net_url): '
read net_url

echo 'Please enter the 1st arguments(net_port): '
read net_port

echo 'Please enter the 1st arguments(network_mode): '
read network_mode

curl "http://10.0.2.193:5000/network_stress?duration=$duration&mode=$mode&net_url=$net_url&net_port=$net_port&network_mode=$network_mode" &
