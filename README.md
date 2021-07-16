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
│   CHAT TOPIC ├────┼───┤            │----│            │
│              │    │   │   Archive  │    │  MongoDB   │
└──────────────┘    │   └────────────┘    └────────────┘
                    │
                    │
                    │   ┌────────────┐    ┌─────────────┐
                    └───┤            │----│             │
                        │   Bot      │    │  ...TODO(NLP)
                        └────────────┘    └──────────────
```

## Tech Stack

- 消息中间件[NSQ](https://nsq.io/)
- NoSQL数据库[MongoDB](https://www.mongodb.com/)
- Web框架[Fiber Web Framework](https://gofiber.io/)

## Usage

1. start
```bash
./launch.sh
```

2. stop
```
./stop.sh
```