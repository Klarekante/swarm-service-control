package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// filterServices returns an array of services that are currently deployed in the Swarm mode cluster
// that match the services provided in the input array
func filterServices(services []string) ([]string, error) {
	// Get all deployed services in the Swarm mode cluster
	allDeployedServices, err := getDeployedServices()
	if err != nil {
		return nil, err
	}
	var matchedServices []string
	// Iterate through the input services and check if they match any of the deployed services
	for _, service := range services {
		for _, runningService := range allDeployedServices {
			//if service == runningService {
			if strings.Contains(strings.ToLower(runningService), strings.ToLower(service)) {
				matchedServices = append(matchedServices, runningService)
				break
			}
		}
	}
	return matchedServices, nil
}

// getDeployedServices returns an array of all the services currently deployed in the Swarm mode cluster
func getDeployedServices() ([]string, error) {
	// Run the "docker service ls --format {{.Name}}" command and retrieve the output
	cmd := exec.Command("docker", "service", "ls", "--format", "{{.Name}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting running services: %s", err)
	}
	services := strings.Split(strings.TrimSpace(string(output)), "\n")
	return services, nil
}

// getScale returns the scale of the service provided
func getScale(service string) (string, error) {
	// Run the "docker service inspect --format {{.Spec.Mode.Replicated.Replicas}} service" command and retrieve the output
	cmd := exec.Command("docker", "service", "inspect", "--format", "{{.Spec.Mode.Replicated.Replicas}}", service)
	scale, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting service scale: %s", err)
	}
	return strings.TrimSpace(string(scale)), nil
}

// stopServices stops the services provided in the input array
func stopServices(services []string) error {
	for _, service := range services {
		// Run the "docker service scale service=0" command for each service
		cmd := exec.Command("docker", "service", "scale", service+"=0")
		_, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error stopping service: %s", err)
		}
		fmt.Printf("Service %s stopped successfully\n", service)
	}
	return nil
}

// startServices starts the services provided in the input map, where the key is the service name and the value is the number of replicas
func startServices(services map[string]string) error {
	for name, replicas := range services {
		replicaArray := strings.SplitN(replicas, "/", 2)
		replicas = replicaArray[0]
		// Run the "docker service scale service=replicas" command
		cmd := exec.Command("docker", "service", "scale", name+"="+replicas)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error starting service: %s", err)
		}
		fmt.Printf("Service %s started successfully\n", name)
	}
	return nil
}

// restartServices restarts the services provided in the input array
func restartServices(services []string) error {
	for _, service := range services {
		// Get the current scale of the service
		scale, err := getScale(service)
		if err != nil {
			return fmt.Errorf("error restarting service %s: %s", service, err)
		}
		// Stop the service
		if err := stopServices([]string{service}); err != nil {
			return fmt.Errorf("error restarting service %s: %s", service, err)
		}
		// Start the service with the previous scale
		if err := startServices(map[string]string{service: scale}); err != nil {
			return fmt.Errorf("error restarting service %s: %s", service, err)
		}
		fmt.Printf("Service %s restarted successfully\n", service)
	}
	return nil
}

// backupServices backs up the services to a local file
func backupServices() {
	cmd := exec.Command("docker", "service", "ls", "--format", "{{.Name}} {{.Replicas}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	services := out.String()

	file, err := os.Create("swarm-service-backup.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Split the services string by newline character
	for _, service := range strings.Split(services, "\n") {
		if service != "" {
			// Write the service to the file, followed by a newline
			fmt.Fprintln(file, service)
		}
	}
	fmt.Println("Services backed up successfully")
}

// restoreServices restore the status of all services from a local file
func restoreServices() {
	file, err := os.Open("swarm-service-backup.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a map to store the service name and replica count
	services := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split the line by spaces and store the service name and replica count in the map
		service := strings.Fields(scanner.Text())
		services[service[0]] = service[1]
	}
	// Call the startServices function with the map of services
	if err := startServices(services); err != nil {
		fmt.Errorf("error restoring services: %s", err)
		return
	}
	fmt.Println("Services restored successfully")
}

func main() {
	var services string
	var restart, stop, backup, restore, all bool
	flag.StringVar(&services, "services", "", "Comma separated names of services")
	flag.BoolVar(&restart, "restart", false, "Restart the specified services")
	flag.BoolVar(&stop, "stop", false, "Stop the specified services")
	flag.BoolVar(&backup, "backup", false, "Backup running services")
	flag.BoolVar(&restore, "restore", false, "Restore services from backup")
	flag.BoolVar(&all, "all", false, "Operate on all running services")
	flag.Parse()

	if !restart && !stop && !backup && !restore {
		fmt.Println("Please specify an action to perform (restart, stop, backup, restore)")
		os.Exit(1)
	}
	if (restart && stop) || (restart && backup) || (restart && restore) || (stop && backup) || (stop && restore) || (backup && restore) {
		fmt.Println("Please specify only one action to perform (restart, stop, backup, restore)")
		os.Exit(1)
	}

	if all {
		runningServices, err := getDeployedServices()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		services = strings.Join(runningServices, ",")
	}
	filteredServices, err := filterServices(strings.Split(services, ","))
	fmt.Println("filteredServices: ", filteredServices)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if restart {
		if err := restartServices(filteredServices); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if stop {
		if err := stopServices(filteredServices); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if backup {
		backupServices()
	}
	if restore {
		restoreServices()
	}
}
