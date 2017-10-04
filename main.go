package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	shouldsave        = false
	smoothingdetected = false
)

func main() {
	local := os.Getenv("localappdata")
	path := local + `\FortniteGame\Saved\Config\WindowsClient\GameUserSettings.ini`
	// Make sure we can read and write to this file right now
	os.Chmod(path, 0600)

	fmt.Println("Backing up existing config")
	backupFile(path)

	// Dump the config to a string, it's relatively short
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %v", err)
		os.Exit(1)
	}

	lines := strings.Split(string(cfg), "\r\n")

	for i, line := range lines {
		accelerationCheck(line, lines, i)
		sensitivityCheck(line, lines, i)
		smoothingCheck(line, lines, i)
	}

	// If smoothing isn't in the config (it isn't default), prompt to add it
	if !smoothingdetected {
		text := linePrompt("Do you want to disable mouse smoothing? 'yes' or 'no'")
		if text == "yes" {
			lines = append(lines, "bEnableMouseSmoothing=False")
			shouldsave = true
		}
	}

	if shouldsave {
		saveCfg(lines, path)
	}
}

// Ask a question and collect some input from the user, leading and trailing newlines for seperation
func linePrompt(question string) string {
	// Leading newline to help seperate questions
	fmt.Printf("\n%s\n", question)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// Save our split config back to the file and set it read only, then exit
func saveCfg(cfg []string, path string) {
	output := strings.Join(cfg, "\r\n")
	err := ioutil.WriteFile(path, []byte(output), 0400)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write file: %v", err)
	}
	os.Chmod(path, 0400)
	fmt.Println("Done!")
	os.Exit(0)
}

// Helper to create a copy of a file
func backupFile(from string) {
	// Read all content of 'from' to data
	data, err := ioutil.ReadFile(from)
	if err != nil {
		log.Fatal(err)
	}
	// Write data to 'to'
	err = ioutil.WriteFile(from+"_BACKUP", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Check for the mouse accel config line, prompt to enable or disable
func accelerationCheck(line string, lines []string, i int) {
	if strings.Contains(line, "bDisableMouseAcceleration") {
	AccelInput:
		for {
			text := linePrompt("Do you want mouse acceleration 'on' or 'off'?")

			switch text {
			case "on":
				lines[i] = "bDisableMouseAcceleration=False"
				shouldsave = true
				break AccelInput
			case "off":
				lines[i] = "bDisableMouseAcceleration=True"
				shouldsave = true
				break AccelInput
			default:
				continue
			}
		}
	}
}

// Check for the mouse sensitivity config line, prompt to change it
func sensitivityCheck(line string, lines []string, i int) {
	if strings.Contains(line, "MouseSensitivity=") {
		split := strings.Split(line, "=")
		sens := strings.TrimSpace(split[1])

		fmt.Printf("\n(Currently: %s)", sens)
		text := linePrompt("Please enter your desired mouse sensitivity (just hit enter, or type 'keep', to keep current):")
		if text != "" && text != "keep" {
			lines[i] = "MouseSensitivity=" + text
			shouldsave = true
			fmt.Printf("Sensitivity set to %s\n", text)
		}
	}
}

// Check for the mouse smoothing config line, prompt to enable or disable
func smoothingCheck(line string, lines []string, i int) {
	if strings.Contains(line, "bEnableMouseSmoothing") {
		smoothingdetected = true
	SmoothingInput:
		for {
			text := linePrompt("Do you want mouse smoothing 'on' or 'off'?")

			switch text {
			case "on":
				lines[i] = "bEnableMouseSmoothing=True"
				shouldsave = true
				break SmoothingInput
			case "off":
				lines[i] = "bEnableMouseSmoothing=False"
				shouldsave = true
				break SmoothingInput
			default:
				continue
			}
		}
	}
}
