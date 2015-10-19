package main

var etcdTemplate = `apiVersion: v1
kind: Pod
metadata: 
  name: etcd
spec: 
  hostNetwork: true
  containers: 
    - name: "etcd"
      image: "{{.Image}}"
      args: 
        - "--name={{.Name}}"
        - "--advertise-client-urls={{.AdvertiseClientUrls}}"
        - "--listen-client-urls=http://0.0.0.0:2379"
        - "--listen-peer-urls=http://0.0.0.0:2380"
        - "--data-dir=/var/lib/etcd/data"
        - "--wal-dir=/var/lib/etcd/wal"
        - "--election-timeout=1000"
        - "--heartbeat-interval=100"
        - "--snapshot-count=10000"
        - "--max-snapshots=5"
        - "--max-wals=5"
        - "--initial-advertise-peer-urls={{.InitialAdvertisePeerUrls}}"
        - "--initial-cluster={{.InitialCluster}}"
        - "--initial-cluster-state=new"
        - "--initial-cluster-token={{.ClusterToken}}"
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
  volumes:{{if eq .CloudProvider "gce"}}
    - name: "etcd-wal"
      gcePersistentDisk:
        pdName: {{.WalVolumeId}}
        fsType: ext4
    - name: "etcd-data"
      gcePersistentDisk:
        pdName: {{.DataVolumeId}}
        fsType: ext4{{else}}
    - name: "etcd-wal"
      awsElasticBlockStore:
        volumeID: {{.WalVolumeId}}
        fsType: ext4
    - name: "etcd-data"
      awsElasticBlockStore:
        volumeID: {{.DataVolumeId}}
        fsType: ext4{{end}}
`
