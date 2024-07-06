package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)



var ctx = context.Background()


func initRedisClient() *redis.Client{
	rdb := redis.NewClient(&redis.Options{

		Addr: os.Getenv("Redis_Client_URL"),
		Password: os.Getenv("Redis_Pass"),
		DB: 0,
	})
	return rdb
}


const (
	WindowDuration = time.Minute 
	RequestLimit = 5

)


func ReuquestHandler(rdb *redis.Client, userIdentifier string){
	cT := time.Now().Unix()
	addRequestTimeStampToRedis(rdb, userIdentifier, cT )
	a := requestChecker(rdb, userIdentifier, cT)
	if a {
		fmt.Println("your are alowed siuuuuuuu", userIdentifier)
	}else{
		fmt.Println("your request is denied coz our rate limiter is working fine")
	}

}

func addRequestTimeStampToRedis(rdb *redis.Client, userIdentifier string, cT int64){
	redisKey := "requests:" + userIdentifier
	rdb.ZAdd(ctx, redisKey, &redis.Z{
		Score: float64(cT),
		Member: cT,
	})
}

func requestChecker(rdb *redis.Client, userIdentifier string, cT int64) bool{
	redisKey := "requests:" + userIdentifier
	windowStartTime := cT - int64(WindowDuration.Seconds())

	rCount, err := rdb.ZCount(ctx, redisKey, fmt.Sprintf("%d", windowStartTime), fmt.Sprintf("%d", cT)).Result()
	if err != nil{
		fmt.Println("Error counting requests mtsm", err)
		return false
	}

	return rCount < RequestLimit
}

////// background worker////////////////\\

// func cleanOldRequest(rdb *redis.Client, userIdentifier string){
// 	redisKey := "requests:" + userIdentifier

// 	cT := time.Now().Unix()
// 	windowStartTime := cT - int64(WindowDuration.Seconds())

// 	rdb.ZRemRangeByScore(ctx, redisKey, "0",  fmt.Sprintf("%d", windowStartTime))

// }


// func backgroundCleanJob(rdb *redis.Client, cleanUpInterval time.Duration, users []string){
// 	for {
// 		time.Sleep(cleanUpInterval)
// 		for _, userIdentifier := range users {
// 			cleanOldRequest(rdb, userIdentifier)
// 		}
// 	}
// }

// rate limiter midlewawre 
func rateLimiterMidleware(rdb *redis.Client) fiber.Handler{
	return func (c *fiber.Ctx) error{
		uI := c.IP()


		cT := time.Now().Unix()

		log.Println(uI, cT)

		addRequestTimeStampToRedis(rdb, uI, cT)

		if !requestChecker(rdb, uI, cT){
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})

		}

		return c.Next()
		
	}
}



func main() {
	godotenv.Load()

	app := fiber.New()
	rdb := initRedisClient()


	
	// users := []string{"user1", "user2", "user3"}

	// for _, user := range users {
	// 	for i := 0; i < 15; i++ {
	// 		ReuquestHandler(rdb, user)
	// 	}
	// }

	app.Use(rateLimiterMidleware(rdb))

	app.Get("/rate-limit", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"data" :"this is working fine",
		})
	})

	// Start the background cleanup job
	// go backgroundCleanJob(rdb, time.Minute, users)

	app.Listen(":8000")

	// select {}

}
