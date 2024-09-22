package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)
import "github.com/gin-gonic/gin"

// channels
func process(c chan int) {
	defer close(c) // close the channel at fn end
	for i := 0; i < 10; i++ {
		c <- i
	}
}

// add something generic to the channel
type gasEngine struct {
	gallons uint8
	mpg     float32
}
type electricEngine struct {
	kwh   float32
	mpkwh float32
}
type car[T gasEngine | electricEngine] struct {
	carMake  string
	carModel string
	engine   T
}

func channelHandler(c *gin.Context) {
	ch := make(chan int)
	// use a routine
	go process(ch)
	for i := range ch {
		fmt.Println(i)
	}
	var gasCar = car[gasEngine]{
		carMake:  "Toyota",
		carModel: "1.0",
		engine: gasEngine{
			gallons: 0,
			mpg:     0,
		},
	}
	c.JSON(http.StatusOK, gasCar)
}

// --- Go routine
var wg = sync.WaitGroup{}

// yeepee mutex
var m = sync.RWMutex{}
var dbData = []string{"hello", "world", "master", "senku", "home"}
var result []string

// with RW locks, we can have multiple readers or one writer

// implement a wait group
func dbCall(i int32) {
	// time call sim
	var delay float32 = rand.Float32() * 2000
	time.Sleep(time.Duration(delay) * time.Millisecond)
	m.Lock()
	result = append(result, dbData[i])
	m.Unlock()
	wg.Done()
}

func routineHandler(c *gin.Context) {
	var t0 = time.Now()
	for i := range dbData {
		wg.Add(1)
		go dbCall(int32(i))
	}
	// wait for any active wg element to call done
	wg.Wait()
	fmt.Printf("Total exec time: %v", time.Since(t0))
	c.JSON(200, gin.H{
		"message": result,
	})
	result = []string{}
}

func handler(c *gin.Context) {

	// slice
	var arrSlice = make([]int32, 3, 8)
	fmt.Println(arrSlice)

	// array
	nums := [...]int32{1, 2, 4}

	// append to slice
	arrSlice = append(arrSlice, nums[2])
	fmt.Println(arrSlice)

	// string manipulation
	var helloBuilder = strings.Builder{}
	helloBuilder.WriteString("Hello")
	helloBuilder.WriteString("\nWorld")
	helloWorld := helloBuilder.String()
	fmt.Println(helloWorld)

	// car engine
	var carEngine = CarEngine{mph: 20, gallons: 5}
	var electricEngine = ElectricEngine{kWh: 5, battery: 20}
	fmt.Printf("Car can make trip: %v\n", canMakeTrip(&carEngine))
	fmt.Printf("Electic can make trip: %v\n", canMakeTrip(&electricEngine))

	num, err := log(nums[1])
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	var name = c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "Hello" + string(num) + "world, " + name,
	})
}

type CarEngine struct {
	mph     uint8
	gallons uint8
}

type ElectricEngine struct {
	kWh     uint8
	battery uint8
}

func (c *CarEngine) availableDuration() uint16 {
	// 1 gallon gives 20 miles
	return uint16(c.gallons) * 20 / uint16(c.mph)
}
func (c *ElectricEngine) availableDuration() uint16 {
	// 1 gallon gives 20 miles
	return uint16(c.battery) * 5 / uint16(c.kWh)
}

// Car Use an interface for a car instead
type Car interface {
	availableDuration() uint16
}

// can make it for journey of 100 hours
func canMakeTrip(c Car) bool {
	return c.availableDuration() > 100
}

func log(num int32) (int32, error) {
	var err error
	fmt.Printf("Hello, number is %d\n", num)
	if num == 2 {
		err = errors.New("number is 2")
		return 0, err
	}
	return num, nil
}

func main() {
	router := gin.Default()
	router.GET("/", routineHandler)
	router.GET("/channel", channelHandler)
	router.GET("/:name", handler)
	router.Run(":5000")
}
