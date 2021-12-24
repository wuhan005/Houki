# 🧹 Houki ![Go](https://github.com/wuhan005/Houki/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/Houki)](https://goreportcard.com/report/github.com/wuhan005/Houki) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/Houki)

Customizable MitM proxy.

[![](./.github/readme/bili2pili-demo.png)]()

## Installation

1. Clone the repository

```bash
git clone git@github.com:wuhan005/Houki.git
 ```

2. Build the binary

```bash
cd Houki

go build .
```

3. Run the binary

```bash
./Houki web
```

## Usage

You can create module to intercept and modify the request and response.

Here is a simple example module configuration which replaces the `bilibili` `哔哩哔哩` to `pilipili` `批哩批哩`。

```json
{
  "title": "Bili2Pili",
  "author": "E99p1ant",
  "description": "This is my first module.",
  "response": {
    "on": "url.contains(\"bilibili.com\")",
    "header": {
      "X-MITM": "Houki"
    },
    "body": {
      "replace": {
        "bilibili": "pilipili",
        "哔哩哔哩": "批哩批哩"
      }
    }
  }
}
```

Then click the `START PROXY` button to start the proxy. You can set your browser's proxy manually or just click
the `OPEN BROWSER` to open a new browser window with the proxy.

Enjoy it!

## What/Who is Houki?

[Houki(ほうき)](https://zh.moegirl.org.cn/index.php?title=%E6%89%AB%E5%B8%9A(%E9%AD%94%E5%A5%B3%E4%B9%8B%E6%97%85)), is the
broom used by Elaina in [*Wandering Witch: The Journey of
Elaina*](https://en.wikipedia.org/wiki/Wandering_Witch:_The_Journey_of_Elaina). Elaina uses magic to turn it into a
human, whose appearance is very similar to hers.

## License

MIT
