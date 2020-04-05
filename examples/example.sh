#!/bin/sh

for n in $(seq 100); do 
    jq -r @json examples/job.json|klystron&
done
