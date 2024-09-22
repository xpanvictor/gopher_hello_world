package main

import (
	"errors"
	"fmt"
)
import "github.com/gin-gonic/gin"

func handler(c *gin.Context) {

	// slice
	var arrSlice = make([]int32, 3, 8)
	fmt.Println(arrSlice)

	// array
	nums := [...]int32{1, 2, 4}

	// append to slice
	arrSlice = append(arrSlice, nums[2])
	fmt.Println(arrSlice)

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
	router.GET("/:name", handler)
	router.Run(":5000")
}
