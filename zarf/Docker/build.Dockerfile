FROM carbon-builder:latest

WORKDIR /project
COPY . .

RUN make package-ui
RUN goreleaser build --auto-snapshot --clean -f zarf/goreleaser-all.yaml
RUN chmod -R 0777 dist/