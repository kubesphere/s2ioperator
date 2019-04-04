# S2I Operator

## 介绍
`S2I`( source to image )是一款由Openshift开发、自动将代码容器化的工具（<https://github.com/openshift/source-to-image>），通过预置的模板来支持多种语言和框架，诸如Java，Nodejs, python等等。S2I Operator将S2I引入到kubernetes中，相比原生命令行的使用方式，有下面几个优点：

   1. 对外提供API，用户可以直接调用API生成自己需要的镜像，或者进行二次开发
   2. 每一次运行的配置文件都存储在k8s中，能够复用
   3. 提供webhook，实现CI/CD
   4. 提供kubectl apply的方式，对k8s使用者友好 


## 如何安装
```bash
kubectl create ns kubesphere-devops-system
kubectl apply -f  https://github.com/kubesphere/s2ioperator/releases/download/v0.0.2/s2ioperator.yaml
```
最新版的S2iOperator加入了验证功能，由于目前[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)的局限性（会在下个版本增强），需要手动添加当前集群的CA。执行完上述命令之后，执行下面的命令添加SSL证书。
```bash
wget https://raw.githubusercontent.com/kubesphere/s2ioperator/master/hack/certs.sh
chmod +x certs.sh
./certs.sh --service webhook-server-service --namespace kubesphere-devops-system --secret webhook-server-secret
```
执行最下面一句命令，需要系统安装`openssl`，`jq`，并且拥有正确配置了k8s集群config的`kubectl`。需要有k8s的管理员权限。
## 快速开始

1. 新建一个s2ibuilder，s2ibuilder存储了所有需要的配置信息。每一次生成Docker镜像可以复用这些信息，也可以覆盖一些信息。

    ```bash
    kubectl apply -f - <<EOF
    apiVersion: devops.kubesphere.io/v1alpha1
    kind: S2iBuilder
    metadata:
        name: s2ibuilder-sample
    spec:
        config:
            displayName: "For Test"
            sourceUrl: "https://github.com/sclorg/django-ex"
            builderImage: centos/python-35-centos7
            imageName: kubesphere/hello-python
            tag: v0.0.1
            builderPullPolicy: if-not-present
    EOF
    ```
    可以通过`kubectl get s2ib` 查看当前所有的S2ibuilder状态
    ```bash
    kubectl get s2ib
    NAME                RUNCOUNT   LASTRUNSTATE   LASTRUNNAME
    s2ibuilder-sample   2          Successful     s2irun-sample1
    ```

2. 开始一次运行s2irun，在builderName中指定使用的s2ibuilder
    ```bash
    kubectl apply -f - <<EOF
    apiVersion: devops.kubesphere.io/v1alpha1
    kind: S2iRun
    metadata:
        name: s2irun-sample
    spec:
        builderName: s2ibuilder-sample
    EOF
    ```
    通过`kubectl get s2ir s2irun-sample`查看当前运行的状态，如果出现了错误，可以查看当前namespace下以"s2irun-sample"开头的pod日志。
    ```bash
    kubectl get s2ir
    NAME             STATE        COMPLETIONTIME
    s2irun-sample    Successful   1m

    $ kubectl logs -f s2irun-sample-xxxxx ##查看具体POD的日志
    ```
3. 查看Job运行的Node节点，利用命令docker image ls 查看编译好的镜像。（*查看S2iBuilder配置指南，学习如何自动将镜像导出到镜像仓库*）
## 如何配置S2ibuilder和Sirun
  1. [S2IBuilder 配置指南](docs/builder_config.md)
  2. [S2IRun 配置指南](docs/run_config.md)
   
## 开源许可
