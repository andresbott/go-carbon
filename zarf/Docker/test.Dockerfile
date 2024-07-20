FROM carbon-builder:latest

WORKDIR /project
COPY . .

RUN make package-ui
RUN make test
RUN make lint
RUN make benchmark
RUN make license-check