#!/bin/bash

echo 'Please enter the 1st arguments(duration): '
read duration

echo 'Please enter the 2st arguments(mem_amount): '
read mem_amount

curl "http://10.0.2.193:5000/memory_stress?duration=$duration&mem_amount=$mem_amount" &
