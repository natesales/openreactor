set positional-arguments

sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/

build component:
    rm -rf {{component}}d
    go build -o {{component}}d ./cmd/{{component}}

exec component:
    just build {{component}}
    sudo ./{{component}}d

turbo-on:
    docker compose exec -it turbo curl localhost:80/turbo/on
turbo-off:
    docker compose exec -it turbo curl localhost:80/turbo/off

logs svc:
    docker compose logs {{svc}}
