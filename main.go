package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	settings string
)

type Config struct {
	CloudProvider string   `yaml:"cloud_provider"`
	ClusterToken  string   `yaml:"cluster_token"`
	Image         string   `yaml:"image"`
	Members       []Member `yaml:"members"`
}

type Member struct {
	AdvertiseClientUrls      string `yaml:"advertise_client_urls"`
	CloudProvider            string `yaml:"cloud_provider"`
	ClusterToken             string `yaml:"cluster_token"`
    DataVolumeId             string `yaml:"data_volume_id"`
	WalVolumeId              string `yaml:"wal_volume_id"`
	InitialAdvertisePeerUrls string `yaml:"initial_advertise_peer_urls"`
	Image                    string `yaml:"image"`
	InitialCluster           string `yaml:"initial_cluster"`
	Name                     string `yaml:"name"`
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	tmpl, err := template.New("etcd").Parse(etcdTemplate)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.CloudProvider == "" {
		log.Fatal("cloud provider required. Supported values: aws, gce.")
	}

	initialCluster := make([]string, 0)
	for _, member := range config.Members {
		m := fmt.Sprintf("%s=%s", member.Name, member.InitialAdvertisePeerUrls)
		initialCluster = append(initialCluster, m)
	}
	initialClusterString := strings.Join(initialCluster, ",")

	for i, member := range config.Members {
		member.ClusterToken = config.ClusterToken
		member.CloudProvider = config.CloudProvider
		member.Image = config.Image
		member.InitialCluster = initialClusterString

		f, err := os.Create(fmt.Sprintf("etcd%d-pod.yaml", i))
		if err != nil {
			log.Println(err)
		}
		err = tmpl.Execute(f, member)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("wrote", f.Name())
		f.Close()
	}
}
