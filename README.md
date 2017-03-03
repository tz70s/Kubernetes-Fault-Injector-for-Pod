# Kubernetes Fault Injection for a single Pod

### Description
Kubernetes pod have multiple user's containers attach on a system-default "pause" container with sharing network namespace.That is, each containers in a pod can easily use localhost to communicate.

```BASH
#
############################################
#####            ############           ####
#####  FAULT     ############  TESTED   ####
##### INJECtION  ############   HTTP    ####
#####  L7PROXY   ############  SERVER   ####
########  -----  ############  ------   ####
########       container/pause          ####
########                                ####
############################################
```

### Functionality
This work reveal simple HTTP fault injection actions before arriving at the endpoint container.
* Simply change fault injection policy over simple HTTP GET with query string and relative url route 
* 80% for normal traffic and 20% for tested traffic based on L7 Reverse Proxy and route handler

### Current Functions
* simpleResponse
```BASH
curl -qa http://inject.default.svc.cluster.local:8282/inject?policy=simpleResponse
```
* delayResponse
Set 3secs response as default
```BASH
curl -qa http://inject.default.svc.cluster.local:8282/inject?policy=delayResponse
```
* statusResponse
Set status code you like
```bash
curl -qa http://inject.default.svc.cluster.local:8282/inject?status=404
```

