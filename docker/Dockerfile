FROM golang:1.22.0

RUN apt update -yq \
    && apt install poppler-utils -yq \
    && apt purge -y --auto-remove

RUN apt update -yq \
    && apt install bash-completion -yq \
    && apt purge -y --auto-remove
