FROM golang:1.16-buster

RUN apt-get update \
    && apt-get install -y make

ARG USERNAME=user
ARG GROUPNAME=user
ARG UID=1001
ARG GID=1002
ARG PASSWORD=user
RUN apt-get install -y sudo \
    && groupadd -g 998 docker \
    && groupadd -g $GID $GROUPNAME \
    && useradd -m -s /bin/bash -u $UID -g $GID -G sudo -G docker $USERNAME \
    && gpasswd -a $USERNAME docker \
    && echo $USERNAME:$PASSWORD | chpasswd \
    && echo "$USERNAME   ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
