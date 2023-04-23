package tools

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const defaultCharRange = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type randomGenerator interface {
	GenerateFloat01()
	GenerateString(length int, chars string)
	GenerateInteger(min, max int)
	GenerateUUID() error
}

type randomGeneratorImpl struct{}

var randomGen randomGenerator = randomGeneratorImpl{}

func (r randomGeneratorImpl) GenerateFloat01() {
	fmt.Println(rand.Float64())
}

func (r randomGeneratorImpl) GenerateString(length int, chars string) {
	var buf bytes.Buffer

	for i := 0; i < length; i++ {
		randIndex := rand.Intn(len(chars))
		randChar := chars[randIndex]
		buf.WriteByte(randChar)
	}

	fmt.Println(buf.String())
}

func (r randomGeneratorImpl) GenerateInteger(min, max int) {
	fmt.Println(rand.Intn(max-min) + min)
}

func (r randomGeneratorImpl) GenerateUUID() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	fmt.Println(uuid.String())
	return nil
}

func RandTool() error {
	flag.Usage = printRandomGeneratorHelp

	var helpFlag bool
	var intFlag int
	var strFlag int
	var uuidFlag bool

	// Define command line flags
	flag.BoolVar(&helpFlag, "h", false, "Show help message")
	flag.BoolVar(&helpFlag, "help", false, "Show help message")
	flag.IntVar(&intFlag, "int", -1, "Generate a random integer")
	flag.IntVar(&strFlag, "str", -1, "Generate a random string")
	flag.BoolVar(&uuidFlag, "u", false, "Generate a random uuid")
	flag.BoolVar(&uuidFlag, "uuid", false, "Generate a random uuid")

	// Parse command line flags
	flag.Parse()
	args := flag.Args()

	// Print help message if -h flag is provided
	if helpFlag {
		flag.Usage()
		return nil
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time

	if uuidFlag {
		return randomGen.GenerateUUID()
	}

	if isFlagPassed("int") {
		if len(args) > 1 {
			flag.Usage()
			return nil
		}

		max := intFlag
		min := 0

		if max <= 0 {
			return fmt.Errorf("invalid max value: %v. Must be greater than 0", max)
		}

		if len(args) == 1 {
			minOpt, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			min = minOpt
		}

		if min >= max {
			return fmt.Errorf("invalid min value: %v. Must be less than max value: %v", min, max)
		}

		randomGen.GenerateInteger(min, max)
		return nil
	}

	if isFlagPassed("str") {
		if len(args) > 1 {
			flag.Usage()
			return nil
		}

		length := strFlag
		chars := defaultCharRange

		if length <= 0 {
			return fmt.Errorf("invalid length value: %v. Must be greater than 0", length)
		}

		if len(args) == 1 {
			chars = args[0]
		}

		randomGen.GenerateString(length, chars)
		return nil
	}

	// Get remaining args
	if len(args) != 0 {
		flag.Usage()
		return nil
	}

	randomGen.GenerateFloat01()
	return nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func printRandomGeneratorHelp() {
	fmt.Print(`Usage:
random [options]
Generates random numbers, strings and uuids

-h, --help  shows this help
-u, --uuid  generates a random uuid v4
--int  <max> [min] generates a random non-negative integer between min and max (default min=0)
--str  <len> [<characters>] generates an alphanumeric string of length <len> from the set of <characters> provided (default <characters>=a-zA-Z0-9)
`)
}
