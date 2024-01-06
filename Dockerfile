FROM golang
COPY ./bin/output-generation ./output-generation
ENTRYPOINT ./output-generation server
