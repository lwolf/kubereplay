[![Build Status](https://travis-ci.org/lwolf/kubereplay.svg?branch=master)](https://travis-ci.org/lwolf/kubereplay)
[![Docker Repository on Quay](https://quay.io/repository/lwolf/kubereplay-controller-amd64/status "Docker Repository on Quay")](https://quay.io/repository/lwolf/kubereplay-controller-amd64)
[![Go Report Card](https://goreportcard.com/badge/github.com/lwolf/kubereplay)](https://goreportcard.com/report/github.com/lwolf/kubereplay)

# kubereplay

Kubereplay aims to make integration of [Goreplay](https://github.com/buger/goreplay) and [Kubernetes](https://github.com/kubernetes/kubernetes) as easy and automated as possible.

# Current status

This is an early alpha version. It is *not* meant to run in production.

# Quickstart

```
# start minikube with Admission capabilities
$ minikube start --extra-config=apiserver.Admission.PluginNames="Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota"

# start kubereplay controller manager in one console
$ kubebuilder run local

# start initializer controller in the second
$ go run cmd/initializer-controller/main.go --kubeconfig=/Users/lwolf/.kube/config

# add initializer config
$ kubect create -f samle/initializer-configuration.yaml

# create harvester, refinery and test deployment
$ kubectl create -f sample/harvester.yaml
$ kubectl create -f sample/refinery.yaml
$ kubectl create -f sample/echoserver.yaml
```

## Pre-requisites

## Deploying kubereplay


# Troubleshooting


# Contributing


# Changelog
The list of [releases](https://github.com/lwolf/kubereplay/releases) is the best place to look for information on changes between releases.

# Support

