# Pingdom Operator [![Build Status](https://travis-ci.org/markelog/pingdom-operator.svg?branch=master)](https://travis-ci.org/markelog/pingdom-operator) [![GoDoc](https://godoc.org/github.com/markelog/pingdom-operator?status.svg)](https://godoc.org/github.com/markelog/pingdom-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/markelog/pingdom-operator)](https://goreportcard.com/report/github.com/markelog/pingdom-operator)

> Manage your [Pingdom](https://pingdom.com) checks via kubernetes operator

## Set up

Assuming you already set up kubernetes config on your machine

Clone the repo

```sh
$ git clone git@github.com:Nalum/pingdom-operator.git
$ cd pingdom-operator
```

Apply all suplementary kubernetes configs in order to set up the operator

```sh
$ make setup
```

Set up pingdom [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

```sh
$ kubectl apply -f ./deploy/crds/pingdom.giantswarm.io_checks_crd.yaml
```

```sh
$ kubectl apply -f ./deploy/crds/pingdom.giantswarm.io_checks_crd.yaml
```

Now check out the an [example](https://github.com/markelog/pingdom-operator/blob/2a1b15d34086e48ed19a6bf85cc62b5a5c0baa47/deploy/crds/pingdom.giantswarm.io_v1alpha1_checks_cr.yaml) for a check. Edit it to your pleasure.

It is recommended to also set up secret for the check too, also see an [example](https://github.com/markelog/pingdom-operator/blob/master/deploy/example_secret.yaml). But if you couldn't be bother with it :), you can set up your credentials right in the checks (the one above) definition.

Then apply your configs with `kubectl apply -f ...`. Done and done
