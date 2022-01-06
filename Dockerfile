FROM golang:1.16-alpine
COPY ./ /mrplow
RUN cd /mrplow && go build
CMD ["/mrplow/mr-plow", "-config", "/mrplow/config.yml"]
