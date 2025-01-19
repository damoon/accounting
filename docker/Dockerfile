FROM golang:1.23.5

RUN apt update -yq \
    && apt install poppler-utils -yq \
    && apt purge -y --auto-remove

RUN apt update -yq \
    && apt install bash-completion -yq \
    && apt purge -y --auto-remove \
    && echo "source /etc/profile.d/bash_completion.sh" >> /etc/bash.bashrc
