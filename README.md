# nsq-chat

A Chat App based on MOM(NSQ) in golang.

## Overview

```
                                            ┌──────────┐
                                            │          │
                                         ┌──┤   Client1│
                                         │  └──────────┘
                                         │
                                         │  ┌──────────┐
                                         │  │          │
                                         ├──┤   Client2│
                        ┌────────────┐   │  └──────────┘
                        │            ├───┤
                    ┌───┤   Server   │   │  ┌──────────┐
                    │   └────────────┘   │  │          │
                    │                    └──┤  Client...
                    │                       └───────────
┌──────────────┐    │
│              │    │   ┌────────────┐    ┌────────────┐
│   CHAT TOPIC ├────┼───┤            │    │            │
│              │    │   │   Archive  │    │  MongoDB   │
└──────────────┘    │   └────────────┘    └────────────┘
                    │
                    │
                    │   ┌────────────┐    ┌─────────────┐
                    └───┤            │    │             │
                        │   Bot      │    │  ...TODO(NLP)
                        └────────────┘    └──────────────
```

## Tech Stack

- [NSQ](https://nsq.io/)
- [MongoDB](https://www.mongodb.com/)
- [Gin Web Framework](https://gin-gonic.com/)

## Usage

1. start MongoDB

```bash
mongod
```
2. start NSQ

```bash
nsqlookupd

nsqd --lookupd-tcp-address=127.0.0.1:4160
```

3. start chat server

```bash
go run main.go
```

4. start archive

```bash
go run archive/archive.go
```

5. start bot

```bash
go run bot/bot.go
```