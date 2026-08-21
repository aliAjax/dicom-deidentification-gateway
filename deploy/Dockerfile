FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /out/dicom-server ./cmd/server
FROM alpine:3.20
RUN adduser -D -u 10001 dicom && mkdir -p /var/lib/dicom/spool /var/lib/dicom/exports && chown -R dicom:dicom /var/lib/dicom
COPY --from=build /out/dicom-server /usr/local/bin/dicom-server
USER dicom
EXPOSE 8093 11112
ENTRYPOINT ["/usr/local/bin/dicom-server"]
