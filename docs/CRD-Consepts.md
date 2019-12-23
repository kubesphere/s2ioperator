# S2I Operator

In S2I, all resource and (CI/CD) steps are defined with Custom Resource Defintion(CRD). Put another way, you can operate all s2i resource by call k8s api directly. In the case of [Kubesphere](https://github.com/kubesphere/kubesphere), this makes it easy to encapsulate configuration into `s2ibuilders` and `s2ibuildertemplates`.

Following CRD will be used in S2I :

1. s2ibuildertemplates: defines information about S2I builder image.
2. s2ibuilders: all configuration information used in building are stored in this CRD.
3. s2iruns: defines an action about build

Here is a Architecture  to figout relationship about all CRD:

![](s2i_arch.png)

For developer who are interested in S2IRun, please read [doc](https://github.com/kubesphere/s2irun#s2irun) about details

