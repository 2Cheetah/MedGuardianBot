FROM golang:1-bookworm

RUN apt-get update
RUN apt-get install -y zsh tree vim
RUN wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh
RUN sh install.sh --unattended
RUN chsh -s "/usr/bin/zsh" "root"

ENTRYPOINT [ "zsh" ]