## Status

Usable but not yet documented nor commented so read the code to figure out how to use them.

## Packages

### `named`

Simple package used to managed named objects.

### `msgbus`

Lightweight pub/sub message bus. I typically use them to limit component dependency using message-based interface.

Message bus topics are used to isolate different types of messages to their own channel.

### `msgbus.named`

Manages named message buses for loosely wiring process-wide components without discovery and life-cycle issues.

I typically use named message buses for broadcasting app events and bridging external messages.
