# API README

## Pre-requisites

- go
- python3
- virtualenv 
- pip3
- kubectl
- kind
- docker


## Local Testing

### Docker compose 

The docker compose has 2 main features which make the setup very productive:
- build - the images get built instead of you mentioning the image name. here you just pass the directory of the context the docker image build process needs
- watch - the images get rebuild whenver there is any change in the folders. This means that you just need to run `docker compose watch` and whenever you make any changes in the code, the docker compose env is updated automagically!

### Kind Cluster testing


Commands:

```
    kind create cluster --name accio
    kubectl cluster-info --context kind-accio
    kind load docker-image -n accio api-core_api
    kind load docker-image -n accio api-validator_api

```
#postgres-0.postgres.default.svc.cluster.local
notes:

- we load the image to the kind cluster from the local docker repo
- kind cluster has a standard storage provisioner(rancher/local-path) which we will be using for pvc/pv stuff for the postgres statefulset
