# pkgs-observer-agent

___Report Package Status Changes to pkgs.observer___

## Installing

On Debian or RedHat based systems, get the
[latest version of the package](https://github.com/returnearly/pkgs.observer-agent/releases/latest)
and install it.


#### Add GPG key for the repository
```bash
wget -qO- https://packages.returnearly.net/public.asc | gpg --dearmor -o /usr/share/keyrings/returnearly.gpg
```

#### Add repository to the sources list
```bash
echo "deb [arch=all signed-by=/usr/share/keyrings/returnearly.gpg] https://packages.returnearly.net/deb stable main" | sudo tee /etc/apt/sources.list.d/returnearly.list > /dev/null
```

#### Update package lists and install the package
```bash
sudo apt update && sudo apt install pkgs-observer-agent -y
```

#### Settings
Check all your settings are OK in the folder `/etc/pkgs-observer-agent.conf.d/service.conf`. It should look like this:
```
INGEST_ENDPOINT=https://ingest.pkgs.observer/ingress/your-key-goes-here/checks
```

#### Enable the Service
and then start the service with

```
systemctl enable --now pkgs-observer-agent.service
systemctl enable --now pkgs-observer-agent.timer
```

## License

This code is released under the MIT license.