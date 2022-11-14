# Altcoin exchange rates

[![Build Status](https://travis-ci.org/gyengus/altcoin_exchange_rates.svg?branch=master)](https://travis-ci.org/gyengus/altcoin_exchange_rates) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gyengus/altcoin_exchange_rates)](https://goreportcard.com/report/github.com/gyengus/altcoin_exchange_rates) 
[![GoDoc](https://godoc.org/github.com/gyengus/altcoin_exchange_rates?status.svg)](https://godoc.org/github.com/gyengus/altcoin_exchange_rates)

This program get cryptocoins exchange rates, then publish them to a specified MQTT topic.

### Configuration

See `config.example.json`. Copy it to `config.json` and fill `mqtt.server` field. You can add or remove coins in the coins array.

For Home Assistant setup and more information [click here](https://gyengus.hu/2018/01/arfolyamok-megjelenitese?utm_source=github_repo)!

### Tips
- Bitcoin: bc1qx4q5epl7nsyu9mum8edrvp2my8tut0enrz7kcn
- EVM compatible (Ethereum, Fantom, Polygon, etc.): 0x9F0a70A7306DF3fc072446cAF540F6766a4CC4E8
- Litecoin: ltc1qk2gf43u3lw6vzhvah03wns0nkgetg2c7ea0w5r
- Solana: 14SHwk3jTNYdMkEvpbq1j7Eu9iUJ3GySnaBF4kqBR8Ah
- Flux: t1T3x4HExm4nWD7gN68px9zCF3ZFQyneFSK
