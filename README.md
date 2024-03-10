# pkgs-observer-agent

___Report Package Status Changes to pkgs.observer___

## Installing

On Debian or RedHat based systems, get the
[latest version of the package](https://github.com/returnearly/pkgs.observer-agent/releases/latest)
and install it.

Check all your settings are OK in the folder `/etc/pkgs-observer-agent.conf.d`
and then start the service with
`systemctl enable --now pkgs-observer-agent.service`
`systemctl enable --now pkgs-observer-agent.timer`

## License

This code is released under the MIT license.