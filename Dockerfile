FROM ubuntu:latest
LABEL authors="dmytry"

ENTRYPOINT ["top", "-b"]