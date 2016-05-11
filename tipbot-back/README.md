tipbot-back
gets all the github notifications of when @tipbot is mentioned.

User @mentions the tipbot.


@tipbot
@tipbot <amount>
@tipbot <destination>
@tipbot <amount> <destination>
@tipbot <preset> <destination>
@tipbot <currency symbol><amount>
@tipbot <amount><currency symbol>
@tipbot <amount> <preset> <destination>
@tipbot <currency symbol><amount> <destination>
@tipbot <amount><currency symbol> <destination>
@tipbot empty <payment address or account ID>

*200

verbs:
	+1
	tip
	emoji
	beer
	

@tipbot 20 USD  @bob

// TODO: 
- fully support other assets
- allow people to trust gateways in the management page
- add date added to User table

*Listen for incoming notifications
*Parse comments
*Create accounts for people 
*Send payments
*Post replies saying the tip happened
*Automatically make an account for people when they federate.
*Allow people to empty their accounts
*injection attacks


bob*codetip.io



===%
User has sent @blah a beer! Claim at tipbot.com
Keep hacking!

====
@user your intentions are noble but your wallet is poor. Send to user*tipbot.com to refill.

People can define their own units, beer, shot, walrus, etc

User gives Tipbot it's secret key. Tipbot has a DB with github user name associated with account.

Page auth with github.

Create account page
Manage account page
Show pending tips

two processes:
    1) backend for the site
    2) polling for mentions


Endpoints:
    change address
    set tip type
    fetch pending tips?



