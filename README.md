
# webFimos

webFimos is an HTTP server written in Go that handles firewalld rules, it responds to requests that require authentication with a token. The token is generated when the application is run with the `--regen` flag.


On this release it is only able to create ALLOW rules to requests coming from the IP of the requester thowards the specified destination port in the request's body.
The rules are created using the `org.fedoraproject.FirewallD1.zone.addRichRule` dbus interface


To regen the token or change configs run:
`webFimos /etc/web-fimos/web-fimos.json --regen --host <HOST> --port <PORT>`

Remember to replace \<HOST> with the host you want the server to expose as and \<PORT> as the port it has to listen to.

Es.:

`webFimos /etc/web-fimos/web-fimos.json --regen --host 127.0.0.1 --port 9091`




## Installation

To install webFimos, follow these steps:

1. Clone the repository:

    ```shell
    git clone https://github.com/thedanielenavarra/webFimos.git
    ```

2. Build the application: (OPTIONAL)
   
   Please note that the Go compiler is necessary to execute this step

    ```shell
    cd SOURCES
    go build
    sudo rpmbuild -bb webFimos.spec
    sudo rpm -i ../RPMS/x86_64/webFimos-1.0-1.el8.x86_64.rpm

    ```


3. Start the HTTP server:

    ```shell
    systemctl start web-fimos
    ```

## Usage

To use webFimos, make HTTP requests to the server and include the generated token in the request headers as follows:

```shell
    curl <HOST>:<PORT> -X POST --data '{"destinationPort": <DEST_PORT>, "protocol": "tcp", "duration": <DURATION>, "token": "<TOKEN>"}'
```

Replace:
- HOST with the host where the service is running on
- PORT with the port set
- DEST_PORT with the destination port you want firewalld to "open" to your IP
- DURATION with the duration (in seconds) you want the rule to be active - after that the rule will be removed
- TOKEN with the token generated during the last `--regen`, you can get it by running:

```shell
grep token /var/log/web-fimos/web-fimos.log
```