This is just for me, Patrick, to copy files to the server

## Local

start the server

    cd bin
    ./server 3333 /data/go/src/nli-go/resources /data/go/nli-go/var

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
    go build -o bin/server bin/server.go
    sudo systemctl restart nli-go

## Install the server

To register with systemd allows you to auto-start at startup and to start and stop it.

sudo nano /usr/lib/systemd/system/nli-go.service

~~~
[Unit]
Description=NLI-GO semantic parser and execution engine
After=network.target

[Service]
Type=simple

User=patrick
Group=www-data

WorkingDirectory=/var/www/homepage/nli-go/bin

ExecStart=/var/www/homepage/nli-go/bin/server 3333 /var/www/homepage/nli-go/resources /var/www/homepage/nli-go/var
ExecReload=/bin/kill -s HUP $MAINPID

[Install]
WantedBy=multi-user.target
~~~

Install the service with

    sudo systemctl daemon-reload
    sudo systemctl enable --now nli-go


