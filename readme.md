# Pingdom Operator [![Build Status](https://travis-ci.org/markelog/pingdom-operator.svg?branch=master)](https://travis-ci.org/markelog/pingdom-operator) [![GoDoc](https://godoc.org/github.com/markelog/pingdom-operator?status.svg)](https://godoc.org/github.com/markelog/pingdom-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/markelog/pingdom-operator)](https://goreportcard.com/report/github.com/markelog/pingdom-operator)

> Manage your [Pingdom](https://pingdom.com) checks via kubernetes operator

## Set up

Assuming you already set up kubernetes config on your machine.
Clone the repo:

```sh
$ git clone git@github.com:markelog/pingdom-operator.git
$ cd pingdom-operator
```

Apply all supplementary kubernetes configs in order to set up the operator –

```sh
$ make setup
```

It is recommended to also set up secret for the check too (see an [example](https://github.com/markelog/pingdom-operator/blob/8b64fad921dbaf455b11f13f48f81b2abc7f5fa8/deploy/example_secret.yaml)). But if you couldn't be bothered with it :), you can set up your credentials right in the check definition (see the last paragraph).

Set up the deployment for the operator, so execute either

```sh
$ kubectl apply -f ./deploy/operator.yaml
```

Or without secrets with

```sh
$ kubectl apply -f ./deploy/operator_without_secret.yaml
```

Set up pingdom [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources) –

```sh
$ kubectl apply -f ./deploy/crds/pingdom.giantswarm.io_checks_crd.yaml
```

Now check out the [example](https://github.com/markelog/pingdom-operator/blob/8b64fad921dbaf455b11f13f48f81b2abc7f5fa8/deploy/crds/pingdom.giantswarm.io_v1alpha1_checks_cr.yaml) for an actual check. Edit it to your pleasure.
Then apply your check and secret (you did used the secret, right?) configs with `kubectl apply -f ....`
Done and done.

