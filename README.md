build:
```bash
go mod tidy
go build
```

run:
```bash
./merge_kubeconfig --help
```

usage:
The program will traverse all files ending with `-suffix` in the directory specified by `-kubeConfigDir`, 
read and merge them, and the merged file will be the file specified by `-mergeFile`

You need to put the kubeconfig files of different k8s clusters in the directory specified by `-kubeConfigDir`,
and pay attention to naming these files. The merge will use these kubeconfig file names to distinguish different k8s clusters.

```bash
./merge_kubeconfig -kubeConfigDir ./config -suffix .yaml -mergeFile ./merged.yaml

cat ./merged.yaml
```

Have a good time!

!!!
## Current command functionality has been merged into `https://github.com/snail2sky/bbx`
## usage
```bash
bbx merge kubeconfig -c CONFIG_DIR MERGED_FILENAME
```
