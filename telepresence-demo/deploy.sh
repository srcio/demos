#! /bin/bash

for i in {a,b,c}; do
    kubectl create deployment service-$i --image srcio/telepresence-demo-service-$i
    kubectl expose deployment service-$i --port 80 --target-port 80
done