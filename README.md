# etcd Pod Generator

The `etcd-pod-gen` tool generates a set of pod manifests for deploying a static etcd cluster under the Kubernetes Kubelet manifest directory. The generated templates can be used as a starting point for creating your own pod manifests for a production etcd cluster.

## Usage

### Create Presistant Disks

#### AWS

Create an etcd data volume:

```
$ aws ec2 create-volume --availability-zone eu-west-1a --size 10 --volume-type gp2 
```

Create an etcd wal volume:

```
$ aws ec2 create-volume --availability-zone eu-west-1a --size 10 --volume-type gp2
```

#### GCE

Create an etcd data volume:

```
$ gcloud compute disks create --size=10GB --zone=us-central1-a etcd0-cluster0-data
```

Create an etcd wal volume:

```
$ gcloud compute disks create --size=10GB --zone=us-central1-a etcd0-cluster0-wal
```
