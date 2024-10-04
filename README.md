# GoAlbumApi

GoAlbumApi is a Golang REST API project using [Gin](https://github.com/gin-gonic/gin) to enable a fast image distribution.

**GoAlbumApi Current Features are:**

- creating albums
- configuring resized images for generation
- uploading images
- creating and updating a password for administrative APIs

## Installation and run

Must have Go installed as and pull the respository. change into the directory.

run the command

```
go get -u github.com/gin-gonic/gin
```

and then run

```
go run main.go
```

Visit http://localhost:3000/ to configure the album. After saving the password you should be directed to the temporary administrative form. You must have your current password entered at the top to create or update albums, update the configuration file, or to upload images. The password is delivered in the Authorization HTTP header


## Endpoints

[Review the REST Endpoints available for GoAlbumApi here](APIREFERENCE.md)


## Future Features

* create react app to complement GoAlbumApi and replace temporary administration forms
* add endpoint to delete images
* add endpoint to delete albums
* pull image meta data - latitude, longitude
* allow setting image title
* define default image resizes
* add parameter to hide albums from the public albums endpoint
* require Authorization password to be encoded
* add backend validation to endpoints
* handle logging and exceptions better, return errors in response
* add whitelisting domains to config.json and setup CORS middleware to respect that
