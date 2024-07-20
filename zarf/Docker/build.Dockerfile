FROM carbon-builder:latest

WORKDIR /project
COPY . .

ARG ACTION
RUN echo $ACTION
RUN if [[ -z "$ACTION" ]] ; then echo Argument not provided && exit 1; elif [[ "$ACTION" == "snapshot" ]]; then  goreleaser --snapshot --clean ; fi




RUN chmod -R 0777 dist/