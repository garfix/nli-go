# Deployment of an nli-go app

To deploy an application on a web server

## Build and run the server 

Inside the directory 'bin', create the server on the platform of choice (i.e. Linux)

    go build server.go

This will create a file 'server' in the directory. Then run the server at port 3333

    ./server 3333

## Web apps

There are two web apps: `index.html` (the DBpedia app) and `block.html` (the blocks world app)
