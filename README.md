# gifcreator

### Run
Use the command go run github.com/cosmtrek/air to run with hot reload

## Architecture
![Gifhub_simplified_architecture drawio](https://user-images.githubusercontent.com/39282569/168500512-22550801-b681-4c4d-93e5-b58a395327bf.png)


## Propsed solutions for handling the long task of video encoding and user wait
Converting mp4 to gif is a long process that can go over the http time limit. A couple solutions would be 1) the complicated solution involving a message queue. The pro is scalability, con is complexity for an app that probably will not explode in popularity. And 2) Send the user a response, and use Goroutine(s) to continue processing the file(s). Then use websockets to alert the user that their process is done. This has less scalability, I think. But eliminates the complexity of having to add a message queue.
