##### explain how the projects works and state flow

> staring with the main file `main.go`
```go
func main() {
    r := gin.Default() // using gin server framework for building
    r.Use(gin.Recovery()) // Add Recovery manually
    r.GET("/", func(c *gin.Context) { // starting the welcome end point
        c.JSON(200, gin.H{
            "message": "Hey Go URL Shortener !",
        })
    })

    // it's the most importent route
    // for creating a shourturl incoded using sha256 and base64 
    // store in the redis cache and incoded using the url and user_id
    // from website information
    r.POST("/:create-short-url", func(c *gin.Context) { // 
        handler.CreateShortUrl(c)
    })

    // this url get the real url from the redis cache based in the give
    // url in the reqeust encoded it and found real url
    r.GET("/:shortUrl", func(c *gin.Context) {
        handler.HandleShourtUrlRedirect(c)
    })


    // create a redis client here initalization for redis cache
    store.InitializationRedisClient()

    err := r.Run(":4440") // starting running server
    if err != nil {
        panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
    }
}
```

## for more information about gen server i suggest visited this 
[gin github](https://github.com/gin-gonic/gin) 

># now content for handler file `handlers.go`
```go
package handler

import (
    "url_shorter/shortener" // another package with some tools for help
    "url_shorter/store" // ....

    "net/http" // for http reqeust
    "github.com/gin-gonic/gin" // gin server
)

// request definition data type for reqeuest
type UrlCreationRequest struct { 
    LongUrl string `json:"long_url" binding:"required"`
    UserId string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
    var creationRequest UrlCreationRequest // create object from type
    // of struct  for checking if data is similer to specific data type
    // ShouldBindJSON => for checking if the same structure and type
    if err := c.ShouldBindJSON(&creationRequest); err != nil { 
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // staring creating a shourt url using LongUrl, UserID)
    shortUrl := shortener.GenertateShortLink(creationRequest.LongUrl,
        creationRequest.UserId)

    // saving the url in redis cache(shourtUrl, realUrl, UserId)
    store.SaveUrlMapping(shortUrl, creationRequest.LongUrl,
        creationRequest.UserId)

    // setup the response for the endpoint
    host := "http://localhost:4440/"
    c.JSON(200, gin.H{
        "message": "shourted url created successfuly",
        "short_url": host+shortUrl,
    })
}

// this funciton for rediect endpoint for the shourturl
func HandleShourtUrlRedirect(c *gin.Context) {
    shourtUrl := c.Param("shortUrl") // getting date from reqeust
    // getting date from redis cache using shourtUrl and return longUrl
    initialUrl := store.RetrieveInitialUrl(shourtUrl) 
    c.Redirect(302, initialUrl) // starting redirect or long Url
}
```

## staring for the file store and retive date from redis `store.go`
```go
package store
// Define the struct wrapper around raw Redis client
type StoreService struct  {
    redisClient *redis.Client
}

// Top level declarations for the storeService and Redis context
var (
    storageService = &StoreService{}
    ctx = context.Background()
)
const CasheDuration = 6 * time.Hour

/* initliazation connection for redis 
function that initalization redis connection and return redis
warpper object
*/
func InitializationRedisClient() *StoreService {
    port := 6379 // port for redis
    client := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("localhost:%v", port),
        Password: "",
        DB: 0,
    })

    // in older verision of redis ping don't take a ctx in ping 
    // so if you want to adding this you must update the redis to v8 or v9
    pong, err := client.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("error happend in redis connection: %v", err)
    }

    fmt.Printf("redis Client running in port: %v, pong message: %v", port, pong)
    // get the object from client and store 
    storageService.redisClient = client
    return storageService // return redis wrapper object
}


/* We want to be able to save the mapping between the originalUrl 
and the generated shortUrl url
*/
func SaveUrlMapping(shortUrl string, originalUrl string, userId string){ 
    // saving the shourtUrl and origianlUrl and 
    // Set() function stores a key-value pair in Redis. The third argument 
    // is the expiration time (0 means no expiration).
    // and don't count ctx
    err := storageService.redisClient.Set(ctx, shortUrl, originalUrl, 
        CasheDuration).Err()
    if err != nil {
        log.Fatalf("error happend in saving data redis: %v", err)
    }
}

/* We should be able to retrieve the initial long URL once the short 
is provided. This is when users will be calling the shortlink in the 
url, so what we need to do here is to retrieve the long url and
think about redirect.
*/
func RetrieveInitialUrl(shortUrl string) string {
    // getting long url based on shourtUrl
    result, err := storageService.redisClient.Get(ctx, shortUrl).Result()
    if err != nil {
        log.Fatalf("error happend in retreiving data redis: %v", err)
    }
    return result
}
```
and this how the storage work in project


> [!IMPORTANT]
> what is ctx context.content => in go every thing using singals
    so if you want to cancel an operation you can do it using signals 
    and can settting timeout for each signals so it's using for handle prcess
    and manage and store things

# last thing how ot shourt and ecodeUrl `shoutener.go`
```go
package shortener
import (
    "crypto/sha256" // for hashing
    "fmt"
    "log"
    "math/big" // for using Sum in math

    "github.com/itchyny/base58-go" // for hashing
)

func sha2560f(input string) []byte{
    algorthim := sha256.New() // creating object from sha256
    // write convert date into bytes
    // and injected data into the hash object
    algorthim.Write([]byte(input)) 
    return algorthim.Sum(nil)
}

/*
The Encode method from the base58.BitcoinEncoding object is
    used to encode the provided bytes array into a Base58-encoded string.
- bytes is a []byte array that represents the data you want to encode 
    (for example, a hash or some binary data).
- The Encode method returns two values:
- encoded: This is the Base58-encoded result of your input data, which 
    will be a string of characters from the Base58 alphabet.
- err: This will be nil if the encoding was successful, or it will contain an 
    error if something went wrong during the encoding process (e.g., invalid input).
*/
func base58Encoding(bytes []byte) string{
    encoding := base58.BitcoinEncoding
    encoded, err := encoding.Encode(bytes)
    if err != nil {
        log.Fatalf("error happend in encoding: %v", err)
    }
    return string(encoded)
}


/* The final algorithm will be super straightforward now as we now have our 2
main building blocks already setup, it will go as follow :

- Hashing  initialUrl + userId url with sha256.  Here userId is added to 
prevent providing similar shortened urls to separate users in case they 
want to shorten exact same link, it's a design decision, so some implementations 
do this differently.
- Derive a big integer number from the hash bytes generated during the hasing.
- Finally apply base58  on the derived big integer value and pick the first 8 characters
*/
func GenertateShortLink(initalLink string, userId string) string {
    /*
        Input: You pass in an initalLink (the original URL) and a userId.
        - Generate a SHA-256 hash: The concatenation of initalLink and 
                userId is hashed using SHA-256, resulting in a 32-byte hash.
        - Convert to a number: The hash is converted to a 64-bit 
                integer (genreteNumber).
        - Encode to Base58: The number is then converted to a Base58 
                string (finalString).
        - Return shortened URL: Only the first 8 characters of the Base58 
                string are returned as the final "short link".
    */
    urlHashBytes := sha2560f(initalLink+userId)
    genreteNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
    finalString := base58Encoding([]byte(fmt.Sprintf("%d", genreteNumber)))
    return finalString[:8]
}
```

> [!IMPORTANT]
> for sha256 more info
    [hash256](https://www.youtube.com/watch?v=y3dqhixzGVo) 
> for base85Encoding 
    [baseEncoding85](https://b64encode.com/blog/base85-encoding/) 
