# gifcreator

### Requirements
PostgreSQL for storing user usage and user created file links
Firebase Authentication for login 
GCP Cloud Storage to temporarily store user converted files 
FFMPEG installed on local machine or host machine

### Run
Use the command go run github.com/cosmtrek/air to run with hot reload

## Architecture
![Gifhub_simplified_architecture drawio](https://user-images.githubusercontent.com/39282569/168500512-22550801-b681-4c4d-93e5-b58a395327bf.png)


## Propsed solutions for handling the long task of video encoding and user wait
Converting mp4 to gif is a long process that can go over the http time limit. A couple solutions would be 1) the complicated solution involving a message queue. The pro is scalability, con is complexity for an app that probably will not explode in popularity. And 2) Send the user a response, and use Goroutine(s) to continue processing the file(s). Then use websockets to alert the user that their process is done. This has less scalability, I think. But eliminates the complexity of having to add a message queue.


### Future Improvements
- Spin off notification into its own thing: Redis
- Clip mp4 to mp4 clips with sound. Compression will be important
- Spin off conversion into it's own service: AWS Lambda with FFMPEG support
- Add AWS S3 as a storage option
- Is it possible for users to enter their AWS S3 or GCP Storage so users can just use the conversion service and store in their own buckets? OAuth2 then the user enters in their bucket name, etc??? Reason - more user privacy, less data storage on our part
- More editing options like cropping, quality (1080p resolution, more exotic stuff like bit rate or what not)
- Add youtube-dl to get youtube videos
- Add ability to parse m3u8 to access streaming videos
- Have users connect with each other to share lists and content
