# Fairness of Go channels
Assume there is a channel where routines listen for new data. As soon as there is new data arriving on that channel some routine will pick it up. However, what happens when mulitple routines are ready to receive new data, who will get it? Will the data be distributed evenly?

**Yes, I think so.**

I conducted a small experiment where a producer sends data on a channel but waits a short amount of time before sending the next one, such that all routines have a chance to get ready (sure, better to use `sync.WaitGroup` but a timeout does its job in this case). It turns out for a large enough timeout (10 miliseconds in my case) each routine receives the same amount of data.

Is the routine selected randomly?

**I doubt that.**

If it is random it makes sense that with large number of data each routine will receive roughly equally many (according to the [law of large numbers](https://en.wikipedia.org/wiki/Law_of_large_numbers)). However, even if there are equally many routines as there is data to be sent each routine will receive data exactly one, which is unlikely (`1 / n!` if I'm not mistaken). Thus I think they are deterministically managed to be distributed as evenly as possibly. 