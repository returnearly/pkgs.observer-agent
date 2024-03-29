package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var version = "unknown"

type Package struct {
	RawLine string `json:"raw"`
}

type Payload struct {
	Version        string    `json:"version"`
	Hostname       string    `json:"hostname"`
	MachineID      string    `json:"machine_id"`
	Distribution   string    `json:"distribution"`
	PackageManager string    `json:"package_manager"`
	Packages       []Package `json:"packages"`
}

func getPackageManagerCommand(packageManager string) *exec.Cmd {
	var cmd *exec.Cmd

	switch packageManager {
	case "apt":
		cmd = exec.Command("apt", "list", "--upgradable")
	case "apk":
		cmd = exec.Command("apk", "list", "--upgradable")
	case "pacman":
		cmd = exec.Command("pacman", "-Qu")
	case "yum":
		cmd = exec.Command("yum", "list", "updates")
	case "homebrew":
		cmd = exec.Command("brew", "outdated", "--verbose")
	default:
		log.Fatal("Unsupported package manager")
	}

	return cmd
}

func getUpgradablePackages(packageManager string) ([]Package, error) {
	var packages []Package

	// Get the list of upgradable packages based on the distribution
	var cmd *exec.Cmd = getPackageManagerCommand(packageManager)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := bytes.Split(output, []byte("\n"))
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}

		packages = append(packages, Package{RawLine: string(line)})
	}

	return packages, nil
}

func getDistro() (string, error) {
	var distro string

	switch runtime.GOOS {
	case "linux":
		switch {
		case isUbuntu():
			distro = "ubuntu"
		case isDebian():
			distro = "debian"
		case isArch():
			distro = "arch"
		case isCentOSOrRedHat():
			distro = "centos"
		default:
			return "", fmt.Errorf("unsupported distribution")
		}
	case "darwin": // macOS
		distro = "macos"
	default:
		return "", fmt.Errorf("unsupported OS")
	}

	return distro, nil
}

func getPackageManager(distro string) (string, error) {
	var packageManager string

	switch distro {
	case "ubuntu", "debian":
		packageManager = "apt"
	case "arch":
		packageManager = "pacman"
	case "centos":
		packageManager = "yum"
	case "macos":
		packageManager = "homebrew"
	default:
		return "", fmt.Errorf("unsupported distribution")
	}

	return packageManager, nil
}

func getMachineID() (string, error) {
	var machineID string

	machineID, err := machineid.ProtectedID("pkgs.observer")
	if err != nil {
		log.Fatal(err)
	}

	return machineID, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isUbuntu() bool {
	return fileExists("/etc/lsb-release")
}

func isDebian() bool {
	return fileExists("/etc/debian_version")
}

func isArch() bool {
	return fileExists("/etc/arch-release")
}

func isCentOSOrRedHat() bool {
	return fileExists("/etc/centos-release")
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return
	}

	machineID, err := getMachineID()
	if err != nil {
		fmt.Println("Error getting machine id:", err)
		return
	}

	distro, err := getDistro()
	if err != nil {
		fmt.Println("Error getting distribution:", err)
		return
	}

	packageManager, err := getPackageManager(distro)
	if err != nil {
		fmt.Println("Error getting package manager:", err)
		return
	}

	packages, err := getUpgradablePackages(packageManager)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	payload := Payload{
		Version:        version,
		Hostname:       hostname,
		MachineID:      machineID,
		Distribution:   distro,
		Packages:       packages,
		PackageManager: packageManager,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get the ingest endpoint from the environment variable
	ingestEndpoint := os.Getenv("INGEST_ENDPOINT")
	if ingestEndpoint == "" {
		fmt.Println("Error: INGEST_ENDPOINT environment variable not set")
		return
	}

	req, err := http.NewRequest("POST", ingestEndpoint, bytes.NewBuffer(jsonData))

	req.Header.Set("User-Agent", "pkgs.observer-agent/"+version)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Packages sent successfully to", ingestEndpoint)
}
