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


