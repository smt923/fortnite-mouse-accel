package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	local := os.Getenv("localappdata")
	path := local + `\FortniteGame\Saved\Config\WindowsClient\GameUserSettings.ini`
	// Make sure we can read and write to this file right now
	os.Chmod(path, 0600)

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
			for {
				fmt.Println("Do you want mouse acceleration 'on' or 'off'?")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)

				switch text {
				case "on":
					lines[i] = "bDisableMouseAcceleration=False"
					saveCfg(lines, path)
				case "off":
					lines[i] = "bDisableMouseAcceleration=True"
					saveCfg(lines, path)
				default:
					continue
				}
			}
		}
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
