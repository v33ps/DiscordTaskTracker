# DiscordTaskTrackerBot
This bot will get a task from a user with a reminder time, and remind the user in a private message when the time has been reached

## Example
An example command:

```
r: 8 hours t: fill out your time card
```


## Adding bot to your Discord server
Go to the [Developers Applications](https://discordapp.com/developers/applications). Create a new application.
Go down to "Bot" on the left hand side panel menu.

Then browse to: 

```
https://discordapp.com/oauth2/authorize?&client_id=YOUR_CLIENT_ID_HERE&scope=bot&permissions=0
```

replace `YOUR_CLIENT_ID_HERE` with the Client ID of your discoard application. This is found under the "General information" menu panel.
