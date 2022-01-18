FROM golang:1.17-alpine AS build
WORKDIR /go/src/foxylock
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/foxylock

FROM scratch
COPY --from=build /go/bin /bin/foxylock/
WORKDIR /bin/foxylock
ENTRYPOINT ["/bin/foxylock/foxylock"]