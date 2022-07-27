FROM golang:1.18-bullseye as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY ./ ./
RUN rm -rf ussd_sim.go

RUN go build -o /server


FROM gcr.io/distroless/base-debian11 as deploy

COPY --from=build /server /server
COPY data/data.json data/

EXPOSE 8004

USER nonroot:nonroot

ENTRYPOINT [ "/server" ]