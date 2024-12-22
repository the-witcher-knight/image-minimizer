# Image minimizer

> The Image Minimizer Server is a sample project written in Go to demonstrate image resizing and optimization. It provides a simple API to reduce the size of images, making it a great learning resource for understanding image manipulation and working with Go web servers. 

## Requirements
- **Docker** version 4+
- **make** version 3+

## How to run

```bash
# Build docker image to run the server
make build-dev-image

# Run server
make run
```

## How to Test

We’ve included a sample [resize request](resize-image.http) for easy testing of the API using tools like REST Client in VSCode or IntelliJ IDEA
, or using cURL
```curl
curl --location 'localhost:8080/image/resize' \
--form 'file=@"/Users/locdang/IdeaProjects/image-minimize-go/testdata/sample-image.jpg"'
```