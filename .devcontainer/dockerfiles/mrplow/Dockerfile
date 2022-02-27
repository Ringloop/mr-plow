# To fully customize the contents of this image, use the following Dockerfile as a base and add the RUN statement from this file:
# https://github.com/microsoft/vscode-dev-containers/blob/v0.112.0/containers/debian-10-git/.devcontainer/Dockerfile
FROM mcr.microsoft.com/vscode/devcontainers/base:0-debian-10

# This Dockerfile's base image has a non-root user with sudo access. Use the "remoteUser"
# property in devcontainer.json to use it. On Linux, the container user's GID/UIDs
# will be updated to match your local UID/GID (when using the dockerFile property).
# See https://aka.ms/vscode-remote/containers/non-root-user for details.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Configure apt and install packages
RUN apt-get update \
    && apt-get -y install sudo build-essential bash-completion docker-compose\
    # Clean up
    && apt-get autoremove -y

# Create the user
ADD puthnfo.sh /usr/bin/
COPY foo hnfo / 
RUN if [ "$USERNAME" != "vscode" ]; then \
    puthnfo.sh \
    ; fi

USER $USERNAME

# Configure apt and install go
# RUN wget https://golang.org/dl/go1.17.7.linux-amd64.tar.gz && tar xvf go1.17.7.linux-amd64.tar.gz
# RUN chown -R root:root ./go && mv go /usr/local
# RUN echo "export GOPATH=/home/vscode/work" >> /home/vscode/.profile
# RUN echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> /home/vscode/.profile

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog
