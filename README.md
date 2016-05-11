
Tipbot is a simple way to tip people on github. It allows you to send lumens to anyone with a github account. It monitors comments for any mention of @tipbot.

It also allows you to send directly to <github_name>*codetio.io from your Stellar wallet.

# Parts

## tipbot-back
The service that watches for github notifications.

## tipbot-ws
If we allow people to manage their account from the website this is the backend for that.

## tipbot.com
The source of the website that explains how it all works.

## Federation Server
Have to run a modified version of the Federation server to make an account any time we get a federation request.

# Usage
You can send to anyone on github by simply sending to <github_name>*codetip.io in your Stellar client. This is how you can fill up your own tipbot account.

After that you can send tips directly from comments on github.

Simply @mention the tipbot to issue commands to it.

command | result | example
-----|----|------
@tipbot | Will send 100 XLM from the commenter to the owner of the repo. | @tipbot
@tipbot <destination> | Will send 100 XLM from the commenter to the destination. | @tipbot @nullstyle
@tipbot <amount> <destination> | Will send amount XLM from the commenter to the destination. | @tipbot 1000 @irisli
@tipbot empty <stellar address or account ID> | Will send the entire balance of the commenter's account to the give stellar address. | @tipbot empty jedmccaleb*codetip.io

Not implemented yet but future commands:
```
@tipbot <preset> <destination>
@tipbot <amount> <preset> <destination>
@tipbot <currency symbol><amount> <destination>
@tipbot <amount><currency symbol> <destination>
```

You can check your balance in your tipbot account by looking up your stellar accountID and then using any normal account viewer.


## Building

[gb](http://getgb.io) is used for building and testing.

Given you have a running golang installation, you can build the server with:

```
gb build
```

After successful completion, you should find `bin/federation` is present in the project directory.

## Running tests

```
gb test
```