set positional-arguments

clean:
    rm -rf svc-*

sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/

build component:
    rm -rf svc-{{component}}
    go build -o svc-{{component}} ./cmd/{{component}}

exec component *args="":
    #!/bin/bash
    just build {{component}}
    sudo ./svc-{{component}} ${@:2}

logs svc:
    docker compose logs {{svc}}

ui:
    cd ui && npm run build
    docker compose restart caddy

piper-setup:
    cd tts && just
