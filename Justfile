set positional-arguments

sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/

build component:
    rm -rf {{component}}d
    go build -o {{component}}d ./cmd/{{component}}

exec component:
    just build {{component}}
    sudo ./{{component}}d