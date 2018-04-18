[![Build Status](https://travis-ci.org/lwolf/kubereplay.svg?branch=master)](https://travis-ci.org/lwolf/kubereplay)
[![Docker Repository on Quay](https://quay.io/repository/lwolf/kubereplay-controller-amd64/status "Docker Repository on Quay")](https://quay.io/repository/lwolf/kubereplay-controller-amd64)
[![Go Report Card](https://goreportcard.com/badge/github.com/lwolf/kubereplay)](https://goreportcard.com/report/github.com/lwolf/kubereplay)

# kubereplay

Kubereplay aims to make integration of [Goreplay](https://github.com/buger/goreplay) and [Kubernetes](https://github.com/kubernetes/kubernetes) as easy and automated as possible.

# Current status

This is an early alpha version. It is *not* meant to run in production.

# About

Kubereplay is Kubernetes add-on to automate capturing and redirection of traffic using [Goreplay](https://github.com/buger/goreplay).
It consist of 2 parts that need to run in the cluster - controller-manager and initializer-controller.

## How it works:

Kubereplay creates and manage 2 CRDs: Harvesters and Refineries.

Refinery - is responsible for managing dedicated GoReplay deployment used for receiving data from workloads (harvesters).
 It listens to traffic on tcp socket and then sends it to configured output (stdout, elasticsearch, kafka, http).

Harvester - is used to configure which deployments should controlled by Kubereplay.
Based on selector in Harvester spec Kubereplay will add GoReplay-sidecar to matching deployments.
More about initialization process is in the [docs](docs/initialization.md)


# Quickstart

```
# start minikube with Admission capabilities
$ minikube start --extra-config=apiserver.Admission.PluginNames="Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota"

# start kubereplay controller manager in one console
$ go run cmd/controller-manager/main.go --kubeconfig=~/.kube/config

# start initializer controller in the second
$ go run cmd/initializer-controller/main.go --kubeconfig=~/.kube/config

# add initializer config
$ kubect create -f sample/initializer-configuration.yaml

# create harvester, refinery and test deployment
$ kubectl create -f sample/harvester.yaml
$ kubectl create -f sample/refinery.yaml
$ kubectl create -f sample/echoserver.yaml
```

## Pre-requisites

* Kubernetes v1.9+ with admission capabilities enabled.

## Deploying kubereplay


# Troubleshooting

If you encounter any issues while using Kubereplay, and your issue is not documented, please file an [issue](https://github.com/lwolf/kubereplay/issues).

# Contributing

All kinds of contributions are very much welcome!

* Fork it
* Create your feature branch (git checkout -b my-new-feature)
* Commit your changes (git commit -am 'Added some feature')
* Push to the branch (git push origin my-new-feature)
* Create new Pull Request

# Changelog
The list of [releases](https://github.com/lwolf/kubereplay/releases) is the best place to look for information on changes between releases.

# Support

If you're using kubereplay and want to support the development, buy me a beer at Beerpay!

[![Beerpay](https://beerpay.io/lwolf/kubereplay/badge.svg?style=beer-square)](https://beerpay.io/lwolf/kubereplay)