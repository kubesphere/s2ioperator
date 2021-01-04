# Quick start

In this example you will use S2I Operator to create a S2iBuilder and S2iRun, that simply creates a ready-to-run image and push to DockerHub.

1. Create a secret of `docker-registry` type to authenticate with container registry to push image.

   Create a Secret by providing credentials on the command line:
   
   ```shell
   kubectl create secret docker-registry s2i-secret-sample --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
   ```
   
2. Create a new `s2ibuilder`, all configuration information used in building are stored in `s2ibuilder` . 

   ```yaml
   kubectl apply -f - <<EOF
   apiVersion: devops.kubesphere.io/v1alpha1
   kind: S2iBuilder
   metadata:
     name: s2ibuilder-sample
   spec:
     config:
       export: true
       displayName: "For Test"
       sourceUrl: "https://github.com/kubesphere/devops-java-sample"
       builderImage: kubesphere/java-8-centos7:v2.1.0
       imageName: <your-registry-username>/hello-world # please replace your registry username 
       tag: v0.0.1
       builderPullPolicy: if-not-present
       pushAuthentication:
         secretRef:
           name: s2i-secret-sample # This secret is created by step 1.
   EOF
   ```
   
3. You can use `kubectl get s2ib` to check `s2ibuilder` status.

   ```shell
   kubectl get s2ib
   NAME                RUNCOUNT   LASTRUNSTATE   LASTRUNNAME
   s2ibuilder-sample   2          Successful     s2irun-sample1
   ```

4. To start a building by `s2irun`, `s2irun` define an action about build, and set filed `builderName` to select which `s2ibuilder` in use.

   ```yaml
   kubectl apply -f - <<EOF
   apiVersion: devops.kubesphere.io/v1alpha1
   kind: S2iRun
   metadata:
       name: s2irun-sample
   spec:
       builderName: s2ibuilder-sample
   EOF
   ```

5. Observe created `s2ibuilder` and `s2irun`:

   ```
   kubectl get s2ib
   NAME                RUNCOUNT   LASTRUNSTATE   LASTRUNNAME     LASTRUNSTARTTIME
   s2ibuilder-sample   1          Successful     s2irun-sample   11m
   ```

   ```
   kubectl get s2ir
   NAME            STATE        K8SJOBNAME                       STARTTIME   COMPLETIONTIME   IMAGENAME
   s2irun-sample   Successful   s2irun-sample-d46fc027083d-job   6m39s       5m15s            kubesphere/hello-world
   ```

   ```
   kubectl get po
   NAME                                   READY   STATUS      RESTARTS   AGE
   s2irun-sample-d46fc027083d-job-p44rx   0/1     Completed   0          11m
   ```

6. Finally, an image named `hello-world:v0.0.1` will be builded automatically and pushed to your registry.
