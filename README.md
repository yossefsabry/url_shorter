# URL Shortener Service

This is a simple URL shortener service built using Go, the Gin framework,
and Redis. It allows you to shorten long URLs and provides an endpoint 
for redirection. The service accepts `POST` requests to create short 
URLs and redirects users when they visit those shortened links.

## Features

- **Create Short URLs:** Allows users to shorten long URLs via a POST request.
- **Redirection:** The shortened URL will redirect to the original long URL.
- **Supports Multiple Shortened URLs:** You can shorten any number of URLs.

## Prerequisites

- **Go** (1.18 or newer) installed
- **Redis** installed and running
- **curl** or any REST client installed for making API requests
- **jq** for format the json response in termianl

## Setup Instructions

### Step 1: Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yossefsabry/url_shorter
cd url_shorter
```


### Step 2: Build and Run the Server

Make sure you have Go installed and Redis running on your machine.

- if there is no vendor fordler do this !!
```bash
go mod tidy
go mod vendor
```

Build the project:
```bash
go build -o server
./server
```

### Step 3: Create a Short URL
To create a shortened URL, send a POST request to 
http://localhost:4440/create-short-url with the following body:

```bash
{
    "long_url": "https://www.guru3d.com/news-story/spotted-ryzen-" + 
      "threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html",
    "user_id" : "e0dba740-fc4b-4977-872c-d360239e6b10"
}

```


```bash
curl --request POST \
--data '{
    "long_url": "https://www.guru3d.com/news-story/spotted-ryzen-" + 
"threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html",
    "user_id" : "e0dba740-fc4b-4977-872c-d360239e6b10"
}' \
  http://localhost:9808/create-short-url | jq

```

### Step 4: Response
Upon successfully creating the short URL, the server will respond with the 
following JSON:
```bash
{
    "message": "short url created successfully",
    "short_url": "http://localhost:9808/9Zatkhpi"
}
```

> [!NOTE]
> - The service uses Redis for storing the mapping between shortened 
    URLs and original long URLs.
> - You can specify any user_id when creating short URLs, which can help in
    identifying users in a real-world application.




