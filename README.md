# Hobbit

Hobbit is a proxy service to expose the Synology server to public networks. You should install the control service
(named `server`) on a server with a public IP address, and install the agent (named `client`) on Synology to get
everything works as expect.

## Requirements

This project is written totally in Pure Golang. But we need an X86 based Synology server for testing the spk client
package.

1. Golang 1.19 or above.
2. Synology 7.1 or above.
3. [Pre-commit](https://pre-commit.com/)

## Architecture

```text
 ┌────────────────┐         ┌────────────────┐         ┌────────────────┐
 │                │         │                │         │                │
 │    HTTP (S)    │         │  Port Forward  │         │   Dashboard    │
 │    (80/443)    │         │    (Any ports) │         │     (8080)     │
 │                │         │                │         │                │
 └────────┬───────┘         └────────┬───────┘         └────────┬───────┘
          │                          │                          │
 ┌────────▼───────────┐              │                          │
 │                    │              │                          │
 │   Hobbit Server    ◄──────────────┘                          │
 │                    │                                         │
 │  (With public IP)  ◄─────────────────────────────────────────┘
 │                    │
 └──────────▲─────────┘
            │
            │ Encrypted Tunnel
            │ (QUIC protocol, 8081 port)
            │
            │ Hobbit Transport Protocol
            │ (Manage, Proxy, PING, Authentication)
            │
 ┌──────────▼─────────┐                             ┌────────────────────┐
 │                    │                             │                    │
 │   Hobbit Client    │    Manage (proxy, certs)    │                    │
 │                    ├─────────────────────────────►  Synology Server   │
 │      (Agent)       │                             │                    │
 │                    │                             │                    │
 └────────────────────┘                             └────────────────────┘
```

## Hobbit Transport Protocol
