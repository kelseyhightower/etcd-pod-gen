# etcd Pod Generator

The `etcd-pod-gen` tool generates a set of pod manifests for deploying a static etcd cluster under the Kubernetes Kubelet manifest directory. The generated templates can be used as a starting point for creating your own pod manifests for a production etcd cluster.

## Usage

```
$ etcd-pod-gen <settings>
```

### Example

Create a settings file:

```
cat <<EOF > aws-settings.yaml
image: "b.gcr.io/kuar/etcd:2.2.0"
cluster_token: "cluster0"
cloud_provider: "aws"
members:
  - name: "etcd0"
    advertise_client_urls: "http://10.0.0.1:2379"
    initial_advertise_peer_urls: "http://10.0.0.1:2380"
    data_volume_id: "vol-1234data0"
    wal_volume_id: "vol-1234wal0"
  - name: "etcd1"
    advertise_client_urls: "http://10.0.0.2:2379"
    initial_advertise_peer_urls: "http://10.0.0.2:2380"
    data_volume_id: "vol-1234data1"
    wal_volume_id: "vol-1234wal1"
  - name: "etcd2"
    advertise_client_urls: "http://10.0.0.2:2379"
    initial_advertise_peer_urls: "http://10.0.0.2:2380"
    data_volume_id: "vol-1234data2"
    wal_volume_id: "vol-1234wal2"
EOF
```

Generate Kubernetes pod configs:

```
$ etcd-pod-gen aws-settings.yaml
```
```
wrote etcd0-pod.yaml
wrote etcd1-pod.yaml
wrote etcd2-pod.yaml
```

Results:

```
$ cat etcd0-pod.yaml
```

```
apiVersion: v1
kind: Pod
metadata: 
  name: etcd
spec: 
  hostNetwork: true
    containers: 
      - name: "etcd"
        image: "b.gcr.io/kuar/etcd:2.2.0"
        args: 
          - "--name=etcd0"
          - "--advertise-client-urls=http://10.0.0.1:2379"
          - "--listen-client-urls=http://0.0.0.0:2379"
          - "--listen-peer-urls=http://0.0.0.0:2380"
          - "--data-dir=/var/lib/etcd/data"
          - "--wal-dir=/var/lib/etcd/wal"
          - "--election-timeout=1000"
          - "--heartbeat-interval=100"
          - "--snapshot-count=10000"
          - "--max-snapshots=5"
          - "--max-wals=5"
          - "--initial-advertise-peer-urls=http://10.0.0.1:2380"
          - "--initial-cluster=etcd0=http://10.0.0.1:2380,etcd1=http://10.0.0.2:2380,etcd2=http://10.0.0.2:2380"
          - "--initial-cluster-state=new"
          - "--initial-cluster-token=cluster0"
        ports:
          - name: client
            containerPort: 2379
            protocol: "TCP"
          - name: peer
            containerPort: 2380
            protocol: "TCP"
        resources:
          limits:
            cpu: "1000m"
            memory: "256Mi"
        volumeMounts:
          - name: "etcd-data"
            mountPath: /var/lib/etcd/data
          - name: "etcd-wal"
            mountPath: /var/lib/etcd/wal
    volumes:
      - name: "etcd-wal"
        awsElasticBlockStore:
          volumeID: vol-1234wal0
          fsType: ext4
      - name: "etcd-data"
        awsElasticBlockStore:
          volumeID: vol-1234data0
          fsType: ext4
```

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
