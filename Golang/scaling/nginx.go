package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func fetchNginxStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func saveRequestCount(count int) error {
	data := []byte(strconv.Itoa(count))
	err := os.WriteFile("/home/troumen/scaling/nginx_request_count.txt", data, 0664)
	if err != nil {
		return err
	}
	return nil
}

func readRequestCount() (int, error) {
	data, err := os.ReadFile("/home/troumen/scaling/nginx_request_count.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
		// If file doesn't exist, initialize it with 0
		err := saveRequestCount(0)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	count, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func saveDownscaleStatus(count int) error {
	data := []byte(strconv.Itoa(count))
	err := os.WriteFile("/home/troumen/scaling/downscale.txt", data, 0664)
	if err != nil {
		return err
	}
	return nil
}

func readDownscaleStatus() (int, error) {
	data, err := os.ReadFile("/home/troumen/scaling/downscale.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
		// If file doesn't exist, initialize it with 0
		err := saveDownscaleStatus(0)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	count, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func main() {
	nginxStatusURL := "http://10.0.0.11/nginx_status"
	inventoryPath := "/home/troumen/ansible/inventory.json"

	// Fetching data from the URL
	data, err := fetchNginxStatus(nginxStatusURL)
	if err != nil {
		fmt.Printf("Failed to fetch data from the URL: %v\n", err)
		return
	}

	// Regular expression pattern to extract the number of requests
	requestsPattern := regexp.MustCompile(`Active connections: (\d+)\s+server accepts handled requests\s+\d+\s+\d+\s+(\d+)`)
	match := requestsPattern.FindStringSubmatch(data)

	if len(match) > 0 {
		// Extracting the number of requests
		newRequests, err := strconv.Atoi(match[2])
		if err != nil {
			fmt.Printf("Error parsing number of requests: %v\n", err)
			return
		}

		// Read old request count
		oldRequests, err := readRequestCount()
		if err != nil {
			fmt.Printf("Error reading old request count: %v\n", err)
			return
		}

		// Compare old and new request counts
		diff := newRequests - oldRequests
		fmt.Printf("Difference in requests: %d\n", diff)

		// Save new request count
		err = saveRequestCount(newRequests)
		if err != nil {
			fmt.Printf("Error saving new request count: %v\n", err)
			return
		}

		// Load inventory file
		inventoryData, err := os.ReadFile(inventoryPath)
		if err != nil {
			fmt.Printf("Failed to read inventory file: %v\n", err)
			return
		}

		// Parse inventory JSON
		var inventory map[string]interface{}
		err = json.Unmarshal(inventoryData, &inventory)
		if err != nil {
			fmt.Printf("Failed to parse inventory JSON: %v\n", err)
			return
		}

		// Calculate number of hosts
		hosts := inventory["hosts"].(map[string]interface{})
		numHosts := len(hosts)

		// Calculate requests per host
		requestsPerHost := float64(diff) / float64(numHosts)
		fmt.Printf("New requests per host: %.2f\n", requestsPerHost)
		adding_hosts := requestsPerHost / 50
		fmt.Printf("Adding hosts: %.2f\n", adding_hosts)
		if requestsPerHost >= 50 {
			fmt.Printf("Upscaling\n")
			cmd := exec.Command("ansible-playbook", "/home/troumen/ansible/scaling_upscale.yml", "--vault-password-file", "/home/troumen/ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name_prefix=scaling -e number_of_vms=%d", numHosts+int(math.Ceil(adding_hosts))))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running Ansible upscale playbook:", err)
			}
		} else if numHosts > 1 {
			// Load downscale status
			downscaleStatus, err := readDownscaleStatus()
			if err != nil {
				fmt.Printf("Error reading downscale status: %v\n", err)
				return
			}

			if diff < 50 {
				downscaleStatus++
			} else {
				downscaleStatus = 0
			}

			if downscaleStatus >= 3 {
				// Downscale
				fmt.Printf("Downscaling\n")
				cmd := exec.Command("ansible-playbook", "/home/troumen/ansible/scaling_downscale.yml", "--vault-password-file", "/home/troumen/ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name=scaling%d", numHosts))
				err := cmd.Run()
				if err != nil {
					fmt.Println("Error running Ansible downscale playbook:", err)
				}
				cmd = exec.Command("ansible-playbook", "/home/troumen/ansible/scaling_upscale.yml", "--vault-password-file", "/home/troumen/ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name_prefix=scaling -e number_of_vms=%d", numHosts-1))
				err = cmd.Run()
				if err != nil {
					fmt.Println("Error running Ansible playbook:", err)
				}
				// Reset downscale status
				err = saveDownscaleStatus(0)
				if err != nil {
					fmt.Printf("Error resetting downscale status: %v\n", err)
					return
				}
			} else {
				// Save downscale status
				err = saveDownscaleStatus(downscaleStatus)
				if err != nil {
					fmt.Printf("Error saving downscale status: %v\n", err)
					return
				}
			}

			//fmt.Printf("Downscaling\n")
			//cmd := exec.Command("ansible-playbook", "/home/troumen/ansible/scaling_downscale.yml", "--vault-password-file", "/home/troumen/ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name=scaling%d", numHosts))
			//err := cmd.Run()
			//if err != nil {
			//	fmt.Println("Error running Ansible downscale playbook:", err)
			//}
			//cmd = exec.Command("ansible-playbook", "/home/troumen/ansible/scaling_upscale.yml", "--vault-password-file", "/home/troumen/ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name_prefix=scaling -e number_of_vms=%d", numHosts-1))
			//err = cmd.Run()
			//if err != nil {
			//	fmt.Println("Error running Ansible playbook:", err)
			//}
		}
	} else {
		fmt.Println("No requests data found in the fetched content.")
	}
}
