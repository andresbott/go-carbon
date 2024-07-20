FROM carbon-builder:latest

WORKDIR /project
COPY . .

RUN make package-ui
RUN goreleaser --snapshot --clean --rm-dist
RUN chmod -R 0777 dist/