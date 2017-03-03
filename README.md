# Kubernetes Fault Injection for a single Pod

Kubernetes pod have multiple user's containers attach on a system-default "pause" container with sharing network namespace.That is, each containers in a pod can easily use localhost to communicate.

This work reveal simple HTTP fault injection actions before arriving at the endpoint container.
* Simply change fault injection policy over query string
* 80% for normal traffic and 20% for tested traffic based on L7 Reverse Proxy

* Current Functions
```
simpleResponse
delayResponse
statusResponse
```

