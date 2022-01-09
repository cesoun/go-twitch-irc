# go-twitch-irc

See forked for README:
- [gempir/go-twitch-irc](https://github.com/gempir/go-twitch-irc)

Added support for proxy based net-dial. (socks5)

- [WebShare](https://proxy.webshare.io/dashboard/)
- [WebShare API](https://proxy.webshare.io/docs/#introduction)

## env

variables

    WEBSHARE_API_KEY=<api_key>

set env

    export $(grep -v '^#' .env | xargs)

unset env

    unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)