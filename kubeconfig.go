package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var kubeConfigDir = flag.String("kubeConfigDir", "./config", "The kube config dir will be merged.")
var suffix = flag.String("suffix", ".yaml", "Kube config file suffix.")
var mergeFile = flag.String("mergeFile", "merged.yaml", "Kube config merged file.")

func main() {
	flag.Parse()
	combinedConfig := api.NewConfig()

	err := filepath.Walk(*kubeConfigDir, func(path string, info os.FileInfo, err error) error {
		singleConfig := api.NewConfig()
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), *suffix) {
			config, err := clientcmd.LoadFromFile(path)
			if err != nil {
				return err
			}

			clusterName := strings.TrimSuffix(info.Name(), *suffix)
			userName := clusterName + "-admin"
			// contextName := clusterName + "-context"

			for clusterKey, cluster := range config.Clusters {
				delete(singleConfig.Clusters, clusterKey)
				singleConfig.Clusters[clusterName] = cluster
			}

			for authKey, authInfo := range config.AuthInfos {
				delete(singleConfig.AuthInfos, authKey)
				singleConfig.AuthInfos[userName] = authInfo
			}

			for contextKey, context := range config.Contexts {
				delete(singleConfig.Contexts, contextKey)
				for clusterName := range singleConfig.Clusters {
					context.Cluster = clusterName
				}
				for authKey := range singleConfig.AuthInfos {
					context.AuthInfo = authKey
				}
				singleConfig.Contexts[fmt.Sprintf("%s@%s", userName, clusterName)] = context
			}
			// merge
			for clusterKey, cluster := range singleConfig.Clusters {
				combinedConfig.Clusters[clusterKey] = cluster
			}

			for authKey, authInfo := range singleConfig.AuthInfos {
				combinedConfig.AuthInfos[authKey] = authInfo
			}

			for contextKey, context := range singleConfig.Contexts {
				log.Println(contextKey)
				combinedConfig.Contexts[contextKey] = context
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	combinedConfigBytes, err := clientcmd.Write(*combinedConfig)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	err = os.WriteFile(*mergeFile, combinedConfigBytes, 0644)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Printf("Combined kubeconfig file created: %s", *mergeFile)
}

