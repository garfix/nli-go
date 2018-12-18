# Deployment of an nli-go app

To deploy an application on a web server

## Build the executable

Inside the directory 'app', create the executable on the platform of choice (i.e. Linux)

 go build nli.go

This will create a file 'nli' in the directory of nli.go

## File structure

Create the following directory structure on the server:

 public_html
    css
    img
    js
 resources

And copy all files of the web app, and the necessary resources (for example dbpedia), to the server.
