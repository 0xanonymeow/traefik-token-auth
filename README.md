# traefik-token-auth

### Configuration

The following declaration (given here in YAML) defines a plugin:

```
# Static configuration

experimental:
  plugins:
    traefik-token-auth:
      moduleName: github.com/0xanonymeow/traefik-token-auth
      version: v0.1.1

```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the http.middlewares section:

```
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - tokenAuth@docker

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:80

  middlewares:
    tokenAuth:
      plugin:
        traefik-token-auth:
          headerField: X-Api-Token
          removeHeader: true
          algorithm: sha256
          hashedToken: d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1
```