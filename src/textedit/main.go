package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func main() {
	// Get the file descriptor for standard input
	fd := int(os.Stdin.Fd())

	// Enable raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)

	// Set up a signal handler to catch interrupt signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nCaught signal, exiting...")
		term.Restore(fd, oldState)
		os.Exit(0)
	}()

	// Read and process characters one at a time in raw mode
	fmt.Println("Raw mode enabled. Press 'q' to quit.")
	for {
		var buf [1]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			break
		}

		// Exit if 'q' is pressed
		if buf[0] == 'q' {
			break
		}

		// Print the hexadecimal value of the key pressed
		fmt.Printf("%c", buf[0])
	}

	fmt.Println("\nExiting...")
}

