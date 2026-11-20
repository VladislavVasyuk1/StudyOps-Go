FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY app.go ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/studyops .

FROM scratch
COPY --from=build /out/studyops /studyops
USER 65532:65532
EXPOSE 8888
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=6 CMD ["/studyops", "--healthcheck"]
ENTRYPOINT ["/studyops"]
