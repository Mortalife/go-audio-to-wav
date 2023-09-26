FROM golang:1.21.1-alpine  as build

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o main ./src


FROM jrottenberg/ffmpeg:4.2-alpine311

# Create user and change workdir
RUN adduser --disabled-password --home /home/ffmpgapi ffmpgapi
WORKDIR /home/ffmpgapi

# Copy files from build stage
COPY --from=build /app/main .
RUN chown ffmpgapi:ffmpgapi * && chmod 755 main

EXPOSE 3000

# Change user
USER ffmpgapi

ENTRYPOINT []
CMD [ "./main" ]