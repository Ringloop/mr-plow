FROM mcr.microsoft.com/vscode/devcontainers/base:0-debian-10

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV DEBIAN_FRONTEND=noninteractive

# Configure apt and install packages
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get -y install sudo build-essential bash-completion docker-compose \
    && apt-get autoremove -y

# Create the user
ADD puthnfo.sh /usr/bin/
COPY foo hnfo / 
RUN if [ "$USERNAME" != "vscode" ]; then \
    puthnfo.sh \
    ; fi

# Configure apt and install go
RUN wget https://golang.org/dl/go1.17.7.linux-amd64.tar.gz && tar xvf go1.17.7.linux-amd64.tar.gz
RUN chown -R root:root ./go && mv go /usr/local
RUN echo "export GOPATH=/home/vscode/work" >> /home/vscode/.profile
RUN echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> /home/vscode/.profile

USER $USERNAME
ENV DEBIAN_FRONTEND=dialog
