Hi! I am working on a discord bot in golang, and am having some issues with channels.
The code can be found here: https://goplay.space/#f-qvueIa_-v

The basic flow of the bot should be:
- get input from the user
- parse the input, and respond with the parsed input
- start a goroutine to sleep for the amount of time specified in the input
- once done sleeping, get something off the channel from the goroutine
However, I am having a few issues. First, for some reason the function `sendReminder()` which is called by the goroutine (called here: https://goplay.space/#f-qvueIa_-v,116) is running multiple times per call, rather than just once per call like I would expect. 

Second, I am not able to actually get anything out of the channel. I send something through the channel here: https://goplay.space/#f-qvueIa_-v,146 
However, when I go to get the input off the channel here:
https://goplay.space/#f-qvueIa_-v,50
nothing ever happens.

I know I'm messing something stupid up, I just can't for the life of me figure out what. I've made some more basic tests with channels and goroutines and those have all worked. Does anyone have some pointers?