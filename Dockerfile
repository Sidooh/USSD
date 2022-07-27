FROM golang:1.18 as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY ./ ./
RUN rm -rf ussd_sim.go

RUN CGO_ENABLED=0 go build -o /server


FROM gcr.io/distroless/static-debian11

COPY --from=build /server /server
COPY data/data.json data/

EXPOSE 8004

USER nonroot:nonroot

ENTRYPOINT [ "/server" ]