# Fastlinks

## Instructions

Modify /etc/hosts with:
```
# golinks rule
127.0.0.1       go
```

Then modify your localhost iptables to redirect from port 80 to port 12345 (or
whatever port you plan to run fastlinks on). You may n eed to find a way to
persist these iptable rules if you restart often.

```
sudo iptables -t nat -I OUTPUT -p tcp -d 127.0.0.1 --dport=80 -j REDIRECT --to-ports 12345
```


### Firefox
You may need to set `browser.fixup.dns_first_for_single_words` in about:config
to true in order to not have Firefox search when typing `go/<key>`.

Alternatively, you can set `browser.fixup.domainwhitelist.go` to true (you need
to create this new value), since setting the fix for all single first words
might conflict with your other searches. This config may not persist if caches
are wiped however.


## Similar Services
1. [Tailscale](https://github.com/tailscale/golink)
