package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	//"sync"
	"time"
)

func executeKubectlGetNodes(file_currentTime string) {
	filename := fmt.Sprintf("%s_nodes_version.log", file_currentTime)
	
	cmd := exec.Command("kubectl", "get", "nodes")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing 'kubectl get nodes': %v\n", err)
	} else {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		logLine := fmt.Sprintf("==============================================================\ncurrent time: %s\n %s\n", currentTime, string(output))
		fmt.Println(logLine)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening file '%s': %v\n", filename, err)
		}
		_, err = file.WriteString(logLine)
		if err != nil {
			fmt.Printf("Error writing to file '%s': %v\n", filename, err)
		}
		file.Close()
	}
}

func processDomain(domain, file_currentTime string) {

	//var mutex sync.Mutex
	nakedDomain := strings.ReplaceAll(domain, "https://", "")
	nakedDomain = strings.ReplaceAll(nakedDomain, "http://", "")
	nakedDomain = strings.ReplaceAll(nakedDomain, ":", "_")
	filename := fmt.Sprintf("%s_%s.log", file_currentTime, nakedDomain)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file '%s': %v\n", filename, err)
		return
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: -1 * time.Second,
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   1 * time.Second,
	}

	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	resp, err := httpClient.Get(domain)
	if err != nil {
		//mutex.Lock() // 加锁
		file.WriteString(fmt.Sprintf("current time: %s || domain: %s || Error: %v\n", currentTime, domain, err))
		fmt.Printf("current time: %s || domain: %s || Error: %v\n", currentTime, domain, err)
		//mutex.Unlock() // 解锁
	} else {
		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		//mutex.Lock() // 加锁
		file.WriteString(fmt.Sprintf("current time: %s || domain: %s || httpcode: %d\n", currentTime, domain, resp.StatusCode))
		fmt.Printf("current time: %s || domain: %s || httpcode: %d\n", currentTime, domain, resp.StatusCode)
		//mutex.Unlock() // 解锁
	}

	file.Close()
}

func main() {
	// 读取 domains.txt 文件
	file, err := os.Open("domains.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var domainsList []string

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		domainsList = append(domainsList, domain)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	ticker := time.Tick(1 * time.Second)

    file_currentTime := time.Now().Format("2006-01-02_15-04-05")

	for range ticker {
		for _, domain := range domainsList {
			processDomain(domain,file_currentTime)
		}
		executeKubectlGetNodes(file_currentTime)
	}
}

