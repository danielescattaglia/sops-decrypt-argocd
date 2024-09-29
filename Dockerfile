FROM golang:1.23-alpine AS  builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY /src/go.mod /src/go.sum ./
RUN go mod tidy
RUN go mod download
COPY /src/. .
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/sops-decrypt-argocd ./

FROM alpine
COPY --from=builder /usr/local/bin/sops-decrypt-argocd /usr/local/bin/sops-decrypt-argocd
RUN chmod +x /usr/local/bin/sops-decrypt-argocd
ENTRYPOINT ["sops-decrypt-argocd"]

