set positional-arguments

sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/

build component:
    rm -rf {{component}}d
    go build -o {{component}}d ./cmd/{{component}}

exec component *args="":
    #!/bin/bash
    just build {{component}}
    sudo ./{{component}}d ${@:2}

turbo-on:
    docker compose exec -it turbo curl localhost:80/turbo/on
turbo-off:
    docker compose exec -it turbo curl localhost:80/turbo/off

logs svc:
    docker compose logs {{svc}}
