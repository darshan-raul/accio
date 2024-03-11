===== phase 1
- [ ] create a simple api based on net http which serves calls
- [ ] create a templ based api which serves a simple webpage 
- [ ] use htmx to make it interactive with buttons

- once you are comfortable with the concepts of htmx,templ :
- [ ] change to a more modern api framework -> mux/fiber
- [ ] write some test cases and make them pass
- [ ] Create makefile for build,test,run
- [ ] Create a nix environment for same environments everywhere
- [ ] Create dockerfiles and run the service as a container 
- [ ] push the image to github container registry or Gcr, add make commands as you keep adding more actions
- [ ] Run this simple api on a GKE cluster behind a load balancer service
- [ ] Create terraform for the whole setup [gke env + k8s env]
- [ ] Create a gihub actions pipeline to build, test and push the image to GCR
    - [ ] Create a optional step to push to GKE as well


====== phase 2

- [ ] create a mongodb service
- [ ] add test cases for these mongo tasks
- [ ] add secrets/configmaps for the mongo service 
- [ ] manage the secrets using sops
- [ ] update the k8s manifests to use helm
- [ ] use kustomize to then have some patches based on envs
- [ ] Add argocd and have a simple gitops flow for the application
- [ ] Create one more microservice which will be a dashboard for the mongo service
- [ ] Add a keycloak service  [basic setup for now]
- [ ] Front these services using Ingress
- [ ] Have a Linkerd service mesh for the interconnectivity of these services 
- [ ] add relevant firewall rules in the service mesh

====== phase 3

crossplane
actually adding logic to the code
vault for authentication
prometheus
grafana
opentelemetry
local testing scafholling