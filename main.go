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
		if strings.Contains(line, "bDisableMouseAcceleration") {
			reader := bufio.NewReader(os.Stdin)
		AccelInput:
			for {
				fmt.Println("\nDo you want mouse acceleration 'on' or 'off'?")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)

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

		if strings.Contains(line, "MouseSensitivity=") {
			split := strings.Split(line, "=")
			sens := strings.TrimSpace(split[1])
			fmt.Printf("\nPlease enter your desired mouse sensitivity (Current sens: %s)\n(just hit enter, or type 'keep', to keep current):\n", sens)
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text == "" || text == "keep" {
				continue
			} else {
				lines[i] = "MouseSensitivity=" + text
				shouldsave = true
				fmt.Printf("Sensitivity set to %s\n", text)
			}
		}

		if strings.Contains(line, "bEnableMouseSmoothing") {
			smoothingdetected = true
			reader := bufio.NewReader(os.Stdin)
		SmoothingInput:
			for {
				fmt.Println("\nDo you want mouse smoothing 'on' or 'off'?")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)

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

	if !smoothingdetected {
		fmt.Println("\nDo you want to disable mouse smoothing? 'yes' or 'no'")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "yes" {
			lines = append(lines, "bEnableMouseSmoothing=False")
			shouldsave = true
		}
	}

	if shouldsave {
		saveCfg(lines, path)
	}

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
