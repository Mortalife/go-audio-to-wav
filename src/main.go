package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Start a new fiber app
	app := fiber.New()

	// Send a string back for GET calls to the endpoint "/"
	app.Post("/", func(c *fiber.Ctx) error {
		// Read the request body into a temporary file
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println("No file passes:", err)
			return err
		}

		// Create a temporary file
		openFile, err := file.Open()

		if err != nil {
			fmt.Println("Failed to open file:", err)
			return err
		}
		defer openFile.Close()

		tempFile, err := os.CreateTemp("", "input_*")
		if err != nil {
			fmt.Println("Failed to create temporary file:", err)
		}

		defer os.Remove(tempFile.Name())

		// Save the incoming file to the temporary file
		if _, err := io.Copy(tempFile, openFile); err != nil {
			fmt.Println("Failed to save incoming file:", err)
		}

		// Close the temporary file
		if err := tempFile.Close(); err != nil {
			fmt.Println("Failed to close temporary file:", err)
		}

		// Generate a random filename
		randomFilename := getRandomFilename()

		// Start FFmpeg process to read from the named pipe and output to a converted file
		outputFilePath := filepath.Join("/tmp", randomFilename)
		ffmpegCmd := exec.Command("ffmpeg", "-i", tempFile.Name(), outputFilePath)
		ffmpegCmd.Stderr = os.Stderr

		// Wait for FFmpeg process to finish
		err = ffmpegCmd.Run()
		if err != nil {
			fmt.Println("Unable to run ffmpeg:", err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				// There was an error, capture and log the error output
				errorOutput := string(exitErr.Stderr)
				fmt.Println("FFmpeg error output:", errorOutput)
			}

			return err
		}

		// Read the converted file into memory
		convertedFile, err := os.ReadFile(outputFilePath)
		if err != nil {
			fmt.Println("Unable to read converted file:", err)
			return err
		}

		// Delete the temporary files
		os.Remove(outputFilePath)

		// Set the response headers for attachment
		c.Set("Content-Disposition", "attachment; filename="+randomFilename)
		c.Set("Content-Type", "audio/wav")

		// Return the converted file as the response
		return c.Send(convertedFile)
	})

	// Listen on PORT 300
	app.Listen(":3000")
}

func getRandomFilename() string {
	// Generate a random filename based on timestamp and random number
	timestamp := time.Now().UnixNano()
	randomNumber := strconv.Itoa(rand.Intn(100000))
	return "converted_" + strconv.FormatInt(timestamp, 10) + "_" + randomNumber + ".wav"
}
