# traefik-waf

A very opinionated Web Application Firewall (WAF) for Traefik and **primarily
for myself**.

## Demo

You can try out the WAF at https://wafdemo.leodev.xyz. It uses the OWASP Core
Rule Set (CRS) and has the Javascript challenge enabled.

## Rules

This WAF uses [Coraza](https://coraza.io/) under the hood and doesn't ship with
any rules by default. We recommend making a `rules` directory, filling it with
the rules from https://coraza.io/docs/tutorials/coreruleset/ and then mounting
it to `/rules` in the container.

You must pass the environment variable `WAF_RULESET_ENABLED=true` to enable the
ruleset.

## Javascript Check

This WAF includes a Javascript challenge to mitigate bots. It's disabled by
default. To enable it, set the environment variable `WAF_REQUIREJS` to `true` or
set the `x-waf-require-js` header to `1` in your requests (e.g. in a traefik
middleware).

The challenge page will be served to clients that don't have a valid cookie. The
challenge is ~3.5kb in size (uncompressed) using the default settings.

The challenge will perform a proof-of-work check in the browser (using
[jspowobfdata](https://github.com/le0developer/jspowobfdata)).

## Customization

### Partial customization

You can customize parts of the default challenge/blocked pages.

You can set the following environment variables to customize the pages:

| Environment Variable | Description                   | Default Value            |
| -------------------- | ----------------------------- | ------------------------ |
| WAF_NAME             | Name of the WAF               | Web Application Firewall |
| WAF_FOOTER_NAME      | Name shown in the footer      | Web Application Firewall |
| WAF_FOOTER_URL       | URL linked in the footer name | This Github repository   |

You can also insert your own HTML into the `<head>` section of both pages by
mounting a file to `/assets/head.html` in the container. This can be used for
adding custom styles, meta tags or analytics scripts.

### Full customization

You can customize the blocked page and the challenge page by mounting your own
HTML files to `/assets/blocked.html` and `/assets/challenge.html` in the
container.

The `challenge.html` **MUST** include a `<!--CHALLENGE-->` marker where the
challenge script will be injected.

### REF

The `<!--REF-->` marker in both HTML files will be replaced with a reference to
the current request ID for debugging purposes.

We will try to automatically detect this based on incoming request headers. We
look for:

- `X-Request-ID`
- `CF-Ray` (Cloudflare)
- `CDN-Requestid` (Bunny.net)

If none or multiple headers are present, a WAF reference will be generated. This
is to prevent the user from being able to manipulate the reference value by
sending the header themselves.

You can also manually set the header to use by setting the `WAF_REF_HEADER`
environment variable to the name of the header you want to use. You can also set
the header to an impossible value like `-` to always make the WAF generate a
reference.

## Logging

By default only WAF blocks are logged to stdout. You can increase the verbosity
by setting the `WAF_VERBOSITY` environment variable to:

| Level       | Blocks | New Challenges | Other Logs | Coraza Logs |
| ----------- | ------ | -------------- | ---------- | ----------- |
| 0           |        |                |            |             |
| 1 (default) | ✓      |                |            |             |
| 2           | ✓      | ✓              |            |             |
| 3           | ✓      | ✓              | ✓          |             |
| 4           | ✓      | ✓              | ✓          | ✓           |
| 5           | ✓      | ✓              | ✓          | ✓ and TRACE |

## Usage

```yml
services:
	traefik-waf:
		image: ghcr.io/le0developer/traefik-waf:latest
		environment:
			- WAF_RULESET_ENABLED=true
			- WAF_REQUIREJS=true
			- WAF_REF_HEADER=CF-Ray
		volumes:
			- ./rules:/rules:ro
			- ./assets:/assets:ro
		labels:
			traefik.enable: true
			traefik.http.routers.traefik-waf.rule: PathPrefix(`/.waf`)
			traefik.http.routers.traefik-waf.middlewares: waf-replace-assets@docker
			traefik.http.middlewares.waf-replace-assets.replacepathregex.regex: ^/\.waf/(.*)
			traefik.http.middlewares.waf-replace-assets.replacepathregex.replacement: /assets/$1
			traefik.http.services.traefik-waf.loadbalancer.server.port: 8080
			traefik.http.middlewares.waf.forwardauth.address: http://traefik-waf:8080
			traefik.http.middlewares.waf.forwardauth.trustForwardHeader: true
			traefik.http.middlewares.waf-requirejs.headers.customRequestHeaders.x-waf-require-js: "1"
```

And then use the `waf[@docker]` or `waf-requirejs[@docker]` middleware in your
routers (if y ou use `waf-requirejs`, it must be listed before the `waf`
middleware).
