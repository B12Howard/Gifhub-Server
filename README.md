# kithli-api

### Requirements
[PostgreSQL](https://www.postgresql.org/) for storing user usage and user created file links
[Firebase](https://firebase.google.com/) Authentication for login 
[GCP Cloud Storage](https://cloud.google.com/) to temporarily store user converted files 
[FFMPEG](https://ffmpeg.org/) installed on local machine or host machine 


### Getting Started
Install FFMPEG, PostgreSQL. Install Go dependencies.

#### Database
Put PostgreSQL credentials in .env.

#### Enable External Service Credentials
Enable the GCP Cloud Storage service. Put the keys into config.json under GCPCLOUDSTORAGE

Create a Firebase account and enable Firebase Authentication. Put the keys into config.json under FIREBASE

#### Run
go install `https://github.com/cosmtrek/air`
Initialize air with `air init`
Use the command `air -c .air.toml` to run with hot reload

#### Client
[https://github.com/B12Howard/gifiviewer](https://github.com/B12Howard/gifiviewer)

#### TODO Docker

#### Postman Collection
Env and Collection here
[PostmanEnvCollection20220906.zip](https://github.com/B12Howard/kithli-api/files/9500089/PostmanEnvCollection20220906.zip)


## Architecture
![Gifhub_simplified_architecture drawio (3)](https://user-images.githubusercontent.com/39282569/196551643-9d64515f-128e-4c8c-af39-071ce5d43226.png)


## Propsed solutions for handling the long task of video encoding and user wait
Converting mp4 to gif is a long process that can go over the http time limit. A couple solutions would be 1) the complicated solution involving a message queue. The pro is scalability, con is complexity for an app that probably will not explode in popularity. And 2) Send the user a response, and use Goroutine(s) to continue processing the file(s). Then use websockets to alert the user that their process is done. This has less scalability, I think. But eliminates the complexity of having to add a message queue.


### Future Improvements
[Trello](https://trello.com/b/34GbTIKL/gifhub)

- Spin off notification into its own thing: Redis
- Clip mp4 to mp4 clips with sound. Compression will be important
- Spin off conversion into it's own service: AWS Lambda with FFMPEG support
- Add AWS S3 as a storage option
- Is it possible for users to enter their AWS S3 or GCP Storage so users can just use the conversion service and store in their own buckets? OAuth2 then the user enters in their bucket name, etc??? Reason - more user privacy, less data storage on our part
- More editing options like cropping, quality (1080p resolution, more exotic stuff like bit rate or what not)
- Add youtube-dl to get youtube videos
- Add ability to parse m3u8 to access streaming videos
- Have users connect with each other to share lists and content
