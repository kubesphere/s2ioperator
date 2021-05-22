FROM gitpod/workspace-full

# More information: https://www.gitpod.io/docs/config-docker/
RUN sudo rm -rf /usr/bin/hd && \
    curl -L https://github.com/LinuxSuRen/http-downloader/releases/download/v0.0.29/hd-linux-amd64.tar.gz | tar xzv && \
    sudo mv hd /usr/local/bin && \
    hd install cli/cli
