sync:
    rsync -raz --exclude venv --exclude .idea --progress . reactor:~/openreactor/
