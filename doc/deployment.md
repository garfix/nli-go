This is just for me, Patrick, to copy files to the server

The frontend has a separate repository: https://github.com/garfix/blocks-world

## Frontend blocks world

To deploy the frontend:

    yarn quasar build
    rsync -rlt -e "ssh" dist/spa/ patrick@136.144.240.141:/var/www/homepage/blocks-world/

## System blocks

To update the application sources, not the server:

    cd /var/www/homepage/nli-go
    git pull

## Server binary

To rebuild and restart the server

    cd /var/www/homepage/nli-go
    git pull
    sudo systemctl stop nli-go
    go build -o bin/server bin/server.go
    sudo systemctl start nli-go
