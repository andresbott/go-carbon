FROM carbon-builder:latest

WORKDIR /project
COPY . .

RUN make verify