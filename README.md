# GoRevolt (Revolt.chat API Library)

GoRevolt is a [revolt.chat](https://revolt.chat) API library allowing you to write fast and large scale bots in Go. Optimized for high performance and stability. Discord.js devs will right at home with this Revolt.chat library.

![Revolt.chat image](https://github.com/revoltchat/.github/raw/master/screenshots/2022-03.png)

## Under Active Development
Please come back after a stable release.

## Road Map
1. Implement all end points for [Delta](https://developers.revolt.chat/stack/delta/permissions).
2. Implement all event handlers for [Bonfire](https://developers.revolt.chat/stack/bonfire/events).
3. Fast customisable caching layer.
4. Fast data rich API for building bots.
5. Reply filter (await user response in channel, no prefix required).

## Robust & Stable
Written with concurrency in mind. Built in unit tests ensure the library is operating as expected.


```sh
go test -v ./...
```
> Run the unit tests your self by using the command above.

## Events
List of currently implemented events. More to come soon.

- [x] OnReady
- [x] OnMessage
- [x] OnMessageUpdate
- [ ] OnMessageAppend
- [ ] OnMessageDelete
- [x] OnChannelCreate
- [x] OnChannelUpdate
- [x] OnChannelDelete
- [ ] OnChannelGroupJoin
- [ ] OnChannelGroupLeave
- [ ] OnChannelStartTyping
- [ ] OnChannelStopTyping
- [ ] OnChannelAck
- [x] OnServerCreate
- [ ] OnServerUpdate
- [ ] OnServerDelete
- [ ] OnServerMemberUpdate
- [ ] OnServerMemberJoin
- [ ] OnServerMemberLeave
- [ ] OnServerRoleUpdate
- [ ] OnServerRoleDelete
- [ ] OnUserUpdate
- [ ] OnUserRelationship
- [ ] OnUserRelationship
- [ ] OnEmojiCreate
- [ ] OnEmojiDelete