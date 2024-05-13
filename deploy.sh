#!/bin/bash
set -e

TAG=dev-$(date +%s)

SERVICES="api notifier"
for SERVICE in $SERVICES;
do
    cp Dockerfile services/$SERVICE/Dockerfile.tmp

    docker build . -f services/$SERVICE/Dockerfile.tmp --build-arg="SERVICE=$SERVICE" -t $SERVICE:$TAG
    kind load docker-image $SERVICE:$TAG

    rm services/$SERVICE/Dockerfile.tmp
done

helm upgrade canigowatchthis helm/canigowatchthis --install --create-namespace --set global.image.tag=$TAG -f helm/canigowatchthis/values.yaml