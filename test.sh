#!/bin/bash

for i in {1..30}; do
    curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/ping
done
