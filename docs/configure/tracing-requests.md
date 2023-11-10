---
description: "Launch a reverse proxy for tracing requests in to your EC2"
icon: material/shoe-print
status: new
---

# Tracing Requests with a Reverse Proxy

To trace requests sent to your broadcasted EC2, `dns53` comes bundled with an internal reverse proxy. To enable proxying:

```{ .sh .no-select }
dns53 --proxy
```

Once enabled, set the required environment variables to trace both `HTTP` and `HTTPS` requests. It is advised not to proxy any requests to IMDS on your EC2.

```{ .sh .no-select }
export HTTP_PROXY=http://localhost:10080
export HTTPS_PROXY=http://localhost:10080
export NO_PROXY=169.254.169.254
```

```{ .sh .no-select }
curl http://httpbin.org/headers
```

```{ .sh .no-select }
curl https://httpbin.org/ip -k
```

If you do not wish to set any of these environment variables, your preferred CLI tool should support request proxying using a dedicated flag. For `curl`, that is `-x`.

## Changing the proxy port

Feel free to change the default proxy port of `:10080` by using the `proxy-port` flag:

```{ .sh .no-select }
dns53 --proxy --proxy-port 10888
```
