sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/

turbod:
    rm -rf turbod
    go build -o turbod cmd/turbo/main.go
    sudo ./turbod -v

gauged:
    rm -f gauged
    go build -o gauged cmd/gauge/main.go
    sudo ./gauged -v
