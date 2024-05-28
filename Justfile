set positional-arguments

clean:
    rm -rf svc-*

build component:
    rm -rf svc-{{component}}
    CGO_ENABLED=0 go build -o svc-{{component}} ./cmd/{{component}}
    docker compose build {{component}}
    docker compose up -d {{component}}

build-all:
    #!/bin/bash
    for f in cmd/*; do
        f=$(echo $f | sed 's/cmd\///')
        just build $f
    done

exec component *args="":
    #!/bin/bash
    just build {{component}}
    sudo ./svc-{{component}} ${@:2}

logs svc:
    docker compose logs {{svc}}

ui:
    cd ui && npm run build
    docker compose restart caddy
