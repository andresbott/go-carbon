FROM carbon-builder:latest

WORKDIR /project
COPY . .

RUN make package-ui && make verify