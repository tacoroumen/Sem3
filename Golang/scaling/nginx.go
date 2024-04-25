package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"os/exec"
	"os"
)

func fetchNginxStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func saveRequestCount(count int) error {
	data := []byte(strconv.Itoa(count))
	err := ioutil.WriteFile("./nginx_request_count.txt", data, 0664)
	if err != nil {
		return err
	}
	return nil
}

func readRequestCount() (int, error) {
	data, err := ioutil.ReadFile("./nginx_request_count.txt")
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

func main() {
	nginxStatusURL := "http://10.0.0.11/nginx_status"
	inventoryPath := "../ansible/inventory.json"

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
		inventoryData, err := ioutil.ReadFile(inventoryPath)
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
		if (requestsPerHost >= 50) {
			fmt.Printf("Upscaling\n")
			cmd := exec.Command("ansible-playbook", "../ansible/scaling_upscale.yml", "--vault-password-file", "../ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name_prefix=scaling -e number_of_vms=%d", numHosts+1))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running Ansible upscale playbook:", err)
			}
		} else if (numHosts > 1) {
			fmt.Printf("Downscaling\n")
		        cmd := exec.Command("ansible-playbook", "../ansible/scaling_downscale.yml", "--vault-password-file", "../ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name=scaling%d", numHosts))
                        err := cmd.Run()
                        if err != nil {
                                fmt.Println("Error running Ansible downscale playbook:", err)
                        }
			cmd = exec.Command("ansible-playbook", "../ansible/scaling_upscale.yml", "--vault-password-file", "../ansible/vault_file.txt", "-e", fmt.Sprintf("vm_name_prefix=scaling -e number_of_vms=%d", numHosts-1))
        		err = cmd.Run()
        		if err != nil {
            		fmt.Println("Error running Ansible playbook:", err)
			}
       		}
	} else {
		fmt.Println("No requests data found in the fetched content.")
	}
}
