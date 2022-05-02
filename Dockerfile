FROM golang:1.18-bullseye as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build -o /server


FROM gcr.io/distroless/base-debian11

COPY --from=build /server /server
COPY data/data.json data/
#COPY logger ./

EXPOSE 8004

ENTRYPOINT [ "/server" ]