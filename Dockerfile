FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bootstrap ./cmd

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /app/bootstrap /var/runtime/bootstrap

CMD ["bootstrap"]
