FROM golang:1.13-alpine as build
WORKDIR /code
RUN apk add --no-cache \
  git \
  ca-certificates \
  tzdata \
  jq
RUN update-ca-certificates

ARG UID=1000
ARG GID=1000

RUN adduser\
  --disabled-password \    
  --gecos "" \    
  --home "/var/run/klystron" \    
  --shell "/sbin/nologin" \    
  --uid "${UID}" \
  "klystron"

COPY main.go go.mod pdf /code/
RUN go mod download
RUN go mod verify

ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GOARCH=amd64
RUN go build \
  -a \
  -installsuffix cgo \
  -ldflags="-w -s" \
  -o /bin/klystron
WORKDIR /var/run/klystron
RUN chown -R $UID:$GID /var/run/klystron

FROM scratch as prod
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /bin/klystron /bin/klystron
COPY --from=build /var/run/klystron /var/run/klystron
VOLUME /var/run/klystron
EXPOSE 12321
ARG UID=1000
ARG GID=1000
USER $UID:$GID
WORKDIR /var/run/klystron/
CMD ["/bin/klystron", "-server"]

FROM alpine as dev
COPY --from=prod /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=prod /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=prod /etc/passwd /etc/passwd
COPY --from=prod /etc/group /etc/group
COPY --from=prod /bin/klystron /bin/klystron
COPY --from=build /var/run/klystron /var/run/klystron
VOLUME /var/run/klystron
RUN apk add --no-cache jq parallel
EXPOSE 12321
USER 1000:1000
WORKDIR /var/run/klystron/
CMD ["/bin/klystron", "-server"]
