FROM golang:1.14.2-alpine3.11 as builder

WORKDIR /src

ADD ./ /src/

RUN ./build.sh

FROM scratch

ENV MINIO_ENDPOINT=s3.mio.io
ENV MINIO_ACCESS_KEY_ID=AccessKeyID
ENV MINIO_SECRET_ACCESS_KEY=secretAccessKey
ENV MINIO_USE_SSL=false
ENV MINIO_BUCKET_NAME=default
ENV MINIO_OBJECT_NAME=object
ENV MINIO_FILE_NAME=tmp
ENV MINIO_FILE_PATH=/data

COPY --from=builder /src/release/linux/amd64/mincli /

ENTRYPOINT ["/mincli"]