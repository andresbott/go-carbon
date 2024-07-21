FROM carbon-builder:latest

WORKDIR /project
COPY . .

ENV GITHUB_TOKEN=
RUN ech

#RUN make package-ui
#RUN goreleaser --clean
#RUN chmod -R 0777 dist/