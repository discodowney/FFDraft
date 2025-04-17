package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// Check if Docker is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		fmt.Println("Docker is not running. Please start Docker and try again.")
		os.Exit(1)
	}

	// Check if the container is running
	cmd = exec.Command("docker", "ps")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error checking Docker containers: %v\n", err)
		os.Exit(1)
	}

	if !strings.Contains(string(output), "go_app_db") {
		fmt.Println("PostgreSQL container is not running. Starting it...")
		cmd = exec.Command("docker-compose", "up", "-d")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			os.Exit(1)
		}

		// Wait for the container to be ready
		fmt.Println("Waiting for container to be ready...")
		time.Sleep(5 * time.Second)
	}

	// Run the tests
	fmt.Println("Running tests...")
	cmd = exec.Command("go", "test", "-v", "./services/player/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}
}
