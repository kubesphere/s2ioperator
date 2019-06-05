# Builder API

> Builder API是以CRD的形式生成的，所以如果想要直接访问原生Https接口，请先参考https://kubernetes.io/docs/concepts/overview/kubernetes-api/ 根据文档中说明构造出需要使用的API

> 目前CRD还没有合适的工具自动生成API文档，本文通过Example的方式说明如何使用Si2Builder
## S2iBuilder

`S2iBuilder`描述的S2iBuilder运行需要的所有信息，下面是一个完整的最普通的例子：
```json
//通过POST的方式发送给APIServer创建一个S2ibuilder
{
    "apiVersion": "devops.kubesphere.io/v1alpha1",
    "kind": "S2iBuilder",
    "metadata": {
        "labels": {
            "controller-tools.k8s.io": "1.0"
        },
        "name": "s2ibuilder-sample", //必须且Namespace下唯一
        "namespace": "default",//必须，但是集群有默认值
    },
    "spec": {
        "config": {//config包含了运行一个S2ibuilder所需的所有信息，之后会转换为一个ConfigMap
            "builderImage": "centos/python-35-centos7",//必须，运行所需镜像
            "builderPullPolicy": "if-not-present", //必须，拉镜像策略
            "displayName": "For Test",//可选，显示名
            "imageName": "kubesphere/hello-python",//必须，最终生成的镜像
            "sourceUrl": "https://github.com/sclorg/django-ex",//必须，仓库所在地址
            "tag": "latest" //可选，最终生成的镜像的标签，默认latest
        }
    },
    "status": {//status 描述了当前S2ibuilder的状态
        "lastRunName": "s2irun-sample1",//当前运行的任务名称
        "lastRunState": "Successful",//当前运行状态
        "runCount": 1//运行次数
    }
}

```
上述例子展示了一个完整的S2ibuilder所有的必须参数（更多的参数请参考）。其中BuilderImage使用的镜像通常不是用户提供的，而是Source2Image事先准备好的，制作一个BuilderImage需要不少人力。所以我们提供了模板功能。

## Template
用户可以在S2iBuilder指定一个模板来构建运行信息，如下：
```json
{
    "apiVersion": "devops.kubesphere.io/v1alpha1",
    "kind": "S2iBuilder",
    "metadata": {
        "name": "s2ibuilder-sample2",
        "namespace": "default",
    },
    "spec": {
        "config": {
            "builderPullPolicy": "if-not-present",
            "displayName": "For Test",
            "imageName": "kubesphere/hello-python",
            "sourceUrl": "https://github.com/sclorg/django-ex",
            "tag": "latest"
        },
        "fromTemplate": {//fromTemplate用于指定模板，并且传递给模板一些参数
            "name": "s2ibuildertemplate-python"
        }
    },
    "status": {
        "lastRunName": "s2irun-sample1",
        "lastRunState": "Successful",
        "runCount": 1
    }
}

```
用户只需要选择模板名称，后台会自动生成完整的ConfigMap。

## S2iBuilderTemplates
S2IOperator预置了多个模板，前端需要根据这些模板的信息渲染，展现给用户。下面展示了一个模板的参数
```json
{
    "apiVersion": "devops.kubesphere.io/v1alpha1",
    "kind": "S2iBuilderTemplate",
    "metadata": {
        "name": "s2ibuildertemplate-python",
    },
    "spec": {
        "baseImage": "centos/python-35-centos7",//模板使用的BaseImage
        "tags":[
            "latest",
            "3.5"
        ],//tags指定了该模板可选的版本
        "codeFramework": "python",//模板对应的语言和框架
        "version": "0.0.1",//模板的版本控制
        "description":"This container image includes Python 3.5 as a S2I base image for your Python 3.5 applications. The resulting image can be run using Docker.
        ",//模板描述性信息
        "iconPath":"https://x.y.z"//与前端约定的Icon地址 
    }
}
```
前端需要事先获取S2ioperator中的自带的模板，生成可供选择的参数。


## 完整的S2iBuilder Config参数（比较重要的加了中文描述，其他的参数都属于高级用法，暂不需要理会）
```go
type S2iConfig struct {
	// DisplayName is a result image display-name label. This defaults to the
	// output image name.
	DisplayName string `json:"displayName,omitempty"`

	// Description is a result image description label. The default is no
	// description.
	Description string `json:"description,omitempty"`

	// BuilderImage describes which image is used for building the result images.
	BuilderImage string `json:"builderImage,omitempty"`

	// BuilderImageVersion provides optional version information about the builder image.
	BuilderImageVersion string `json:"builderImageVersion,omitempty"`

	// BuilderBaseImageVersion provides optional version information about the builder base image.
	BuilderBaseImageVersion string `json:"builderBaseImageVersion,omitempty"`

	// RuntimeImage specifies the image that will be a base for resulting image
	// and will be used for running an application. By default, BuilderImage is
	// used for building and running, but the latter may be overridden.
	RuntimeImage string `json:"runtimeImage,omitempty"`

	// RuntimeImagePullPolicy specifies when to pull a runtime image.
	RuntimeImagePullPolicy PullPolicy `json:"runtimeImagePullPolicy,omitempty"`

	// RuntimeAuthentication holds the authentication information for pulling the
	// runtime Docker images from private repositories.
	RuntimeAuthentication AuthConfig `json:"runtimeAuthentication,omitempty"`

	// RuntimeArtifacts specifies a list of source/destination pairs that will
	// be copied from builder to a runtime image. Source can be a file or
	// directory. Destination must be a directory. Regardless whether it
	// is an absolute or relative path, it will be placed into image's WORKDIR.
	// Destination also can be empty or equals to ".", in this case it just
	// refers to a root of WORKDIR.
	// In case it's empty, S2I will try to get this list from
	// io.openshift.s2i.assemble-input-files label on a RuntimeImage.
	RuntimeArtifacts []VolumeSpec `json:"runtimeArtifacts,omitempty"`

	// DockerConfig describes how to access host docker daemon.
	DockerConfig *DockerConfig `json:"dockerConfig,omitempty"`

	// PullAuthentication holds the authentication information for pulling the
	// Docker images from private repositories
	PullAuthentication AuthConfig `json:"pullAuthentication,omitempty"` //拉取代码的凭证，需要提供用户名和密码

	// PullAuthentication holds the authentication information for pulling the
	// Docker images from private repositories
	PushAuthentication AuthConfig `json:"pushAuthentication,omitempty"`//推送代码的凭证

	// IncrementalAuthentication holds the authentication information for pulling the
	// previous image from private repositories
	IncrementalAuthentication AuthConfig `json:"incrementalAuthentication,omitempty"`

	// DockerNetworkMode is used to set the docker network setting to --net=container:<id>
	// when the builder is invoked from a container.
	DockerNetworkMode DockerNetworkMode `json:"dockerNetworkMode,omitempty"`

	// PreserveWorkingDir describes if working directory should be left after processing.
	PreserveWorkingDir bool `json:"preserveWorkingDir,omitempty"`

	//ImageName Contains the registry address and reponame, tag should set by field tag alone
	ImageName string `json:"imageName"`
	// Tag is a result image tag name.
	Tag string `json:"tag,omitempty"`

	// BuilderPullPolicy specifies when to pull the builder image
	BuilderPullPolicy PullPolicy `json:"builderPullPolicy,omitempty"`

	// PreviousImagePullPolicy specifies when to pull the previously build image
	// when doing incremental build
	PreviousImagePullPolicy PullPolicy `json:"previousImagePullPolicy,omitempty"`

	// Incremental describes whether to try to perform incremental build.
	Incremental bool `json:"incremental,omitempty"`

	// IncrementalFromTag sets an alternative image tag to look for existing
	// artifacts. Tag is used by default if this is not set.
	IncrementalFromTag string `json:"incrementalFromTag,omitempty"`

	// RemovePreviousImage describes if previous image should be removed after successful build.
	// This applies only to incremental builds.
	RemovePreviousImage bool `json:"removePreviousImage,omitempty"`

	// Environment is a map of environment variables to be passed to the image.
	Environment []EnvironmentSpec `json:"environment,omitempty"`//环境变量，环境变量是定制基础镜像的唯一途径

	// LabelNamespace provides the namespace under which the labels will be generated.
	LabelNamespace string `json:"labelNamespace,omitempty"`

	// CallbackURL is a URL which is called upon successful build to inform about that fact.
	CallbackURL string `json:"callbackUrl,omitempty"`

	// ScriptsURL is a URL describing where to fetch the S2I scripts from during build process.
	// This url can be a reference within the builder image if the scheme is specified as image://
	ScriptsURL string `json:"scriptsUrl,omitempty"`

	// Destination specifies a location where the untar operation will place its artifacts.
	Destination string `json:"destination,omitempty"`

	// WorkingDir describes temporary directory used for downloading sources, scripts and tar operations.
	WorkingDir string `json:"workingDir,omitempty"`

	// WorkingSourceDir describes the subdirectory off of WorkingDir set up during the repo download
	// that is later used as the root for ignore processing
	WorkingSourceDir string `json:"workingSourceDir,omitempty"`

	// LayeredBuild describes if this is build which layered scripts and sources on top of BuilderImage.
	LayeredBuild bool `json:"layeredBuild,omitempty"`

	// Specify a relative directory inside the application repository that should
	// be used as a root directory for the application.
	ContextDir string `json:"contextDir,omitempty"`

	// AssembleUser specifies the user to run the assemble script in container
	AssembleUser string `json:"assembleUser,omitempty"`

	// RunImage will trigger a "docker run ..." invocation of the produced image so the user
	// can see if it operates as he would expect
	RunImage bool `json:"runImage,omitempty"`

	// Usage allows for properly shortcircuiting s2i logic when `s2i usage` is invoked
	Usage bool `json:"usage,omitempty"`

	// Injections specifies a list source/destination folders that are injected to
	// the container that runs assemble.
	// All files we inject will be truncated after the assemble script finishes.
	Injections []VolumeSpec `json:"injections,omitempty"`

	// CGroupLimits describes the cgroups limits that will be applied to any containers
	// run by s2i.
	CGroupLimits *CGroupLimits `json:"cgroupLimits,omitempty"`

	// DropCapabilities contains a list of capabilities to drop when executing containers
	DropCapabilities []string `json:"dropCapabilities,omitempty"`

	// ScriptDownloadProxyConfig optionally specifies the http and https proxy
	// to use when downloading scripts
	ScriptDownloadProxyConfig *ProxyConfig `json:"scriptDownloadProxyConfig,omitempty"`

	// ExcludeRegExp contains a string representation of the regular expression desired for
	// deciding which files to exclude from the tar stream
	ExcludeRegExp string `json:"excludeRegExp,omitempty"`

	// BlockOnBuild prevents s2i from performing a docker build operation
	// if one is necessary to execute ONBUILD commands, or to layer source code into
	// the container for images that don't have a tar binary available, if the
	// image contains ONBUILD commands that would be executed.
	BlockOnBuild bool `json:"blockOnBuild,omitempty"`

	// HasOnBuild will be set to true if the builder image contains ONBUILD instructions
	HasOnBuild bool `json:"hasOnBuild,omitempty"`

	// BuildVolumes specifies a list of volumes to mount to container running the
	// build.
	BuildVolumes []string `json:"buildVolumes,omitempty"`

	// Labels specify labels and their values to be applied to the resulting image. Label keys
	// must have non-zero length. The labels defined here override generated labels in case
	// they have the same name.
	Labels map[string]string `json:"labels,omitempty"`

	// SecurityOpt are passed as options to the docker containers launched by s2i.
	SecurityOpt []string `json:"securityOpt,omitempty"`

	// KeepSymlinks indicates to copy symlinks as symlinks. Default behavior is to follow
	// symlinks and copy files by content.
	KeepSymlinks bool `json:"keepSymlinks,omitempty"`

	// AsDockerfile indicates the path where the Dockerfile should be written instead of building
	// a new image.
	AsDockerfile string `json:"asDockerfile,omitempty"`

	// ImageWorkDir is the default working directory for the builder image.
	ImageWorkDir string `json:"imageWorkDir,omitempty"`

	// ImageScriptsURL is the default location to find the assemble/run scripts for a builder image.
	// This url can be a reference within the builder image if the scheme is specified as image://
	ImageScriptsURL string `json:"imageScriptsUrl,omitempty"`

	// AddHost Add a line to /etc/hosts for test purpose or private use in LAN. Its format is host:IP,muliple hosts can be added  by using multiple --add-host
	AddHost []string `json:"addHost,omitempty"`//运行时提供的一些别名

	//Export Push the result image to specify image registry in tag
	Export bool `json:"export,omitempty"`//是否将镜像推送到外部，通常需要PushAuth配合

	//SourceURL is  url of the codes such as https://github.com/a/b.git
	SourceURL string `json:"sourceUrl"`
}
```
