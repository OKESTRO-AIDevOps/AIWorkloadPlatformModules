#!/bin/bash

echo 'Please enter the 1st arguments(duration): '
read duration

echo 'Please enter the 1st arguments(size_mb): '
read size_mb

curl "http://10.0.2.193:5000/disk_stress?duration=$duration&size_mb=$size_mb" &
