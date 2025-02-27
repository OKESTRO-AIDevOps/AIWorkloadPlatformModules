#!/bin/bash

echo 'Please enter the 1st arguments(duration): '
read duration

echo 'Please enter the 2st arguments(cpu_num): '
read cpu_num

echo 'Please enter the 3st arguments(percentage): '
read percentage

curl "http://10.0.2.193:5000/cpu_stress?duration=$duration&cpu_num=$cpu_num&percentage=$percentage" &
