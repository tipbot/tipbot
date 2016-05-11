package main

import (
	"fmt"
)

func (self *App) SOURCE_ACCOUNT_NOT_MADE(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + ", the sentiment is appreciated but you must create a tipping wallet first. Send lumens to " +
		sourceName + "*" + self.config.FEDERATION_DOMAIN + " to get started. Checkout [Stellar](https://stellar.org) or [" +
		self.config.FEDERATION_DOMAIN + "](http://" + self.config.FEDERATION_DOMAIN + ") if this makes no sense.")
}

func (self *App) NOT_ENOUGH_BALANCE(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + ", your heart is larger than your wallet! Send lumens to " + sourceName + "*" + self.config.FEDERATION_DOMAIN + " to refill")

}

func (self *App) TIP_NOT_HIGH_ENOUGH(sourceName string, destName string) string {
	return fmt.Sprintf("@%s, %s doesn't have a Stellar account yet. You must send at least %.0f lumens to create their account.",sourceName,destName,self.config.MIN_XML_BALANCE)

}

func (self *App) ONLY_XLM() string {
	return fmt.Sprintln("Right now you can only send lumens.")
}

func (self *App) REPORT_TIP(sourceName string, destName string) string {

	return fmt.Sprintln("@" + destName + ", you have recieved a tip from @" + sourceName + ". Claim at " + self.config.WEBSITE_URL)

}

func (self *App) MALFORMED_COMMAND(sourceName string) string {
	return fmt.Sprintln("@" + sourceName +" command is malformed. Check [" + self.config.FEDERATION_DOMAIN + "](http://" + self.config.FEDERATION_DOMAIN + ") for instructions on how to talk to the tipbot.")

}

func (self *App) ALREADY_EMPTY(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + ", your wallet is empty except for the min balance. Checkout [Stellar](https://stellar.org) or [" +
		self.config.FEDERATION_DOMAIN + "](http://" + self.config.FEDERATION_DOMAIN + ") if this makes no sense.")
}

func (self *App) EMPTY_SUCCESS(sourceName string) string {

	return fmt.Sprintln("@" + sourceName + " you have emptied your wallet.")
}

func (self *App) INVALID_ACCOUNT(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + " that isn't a valid Stellar account.")
}

func (self *App) ADDRESS_NOT_FOUND(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + " that address couldn't be found.")
}

func (self *App) SOMETHING_WRONG(sourceName string) string {
	return fmt.Sprintln("@" + sourceName + " there's a problem but don't panic. Maybe you entered in something weird or some server is down. Hard for a simple bot to know. I've notified the Humans and they are looking into it.")
}


