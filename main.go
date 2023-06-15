package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "runtime"
    "strings"
)

// ANSI color codes for console output
const (
    green = "\033[32m"
    red   = "\033[31m"
    reset = "\033[0m"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("[?] Usage: firex.go example.apk")
        fmt.Println("[+] @shanekhantaun9")
        return
    }

    apkPath := os.Args[1]

    // Print status message
    fmt.Println("[*] Decompiling APK file...")

    // Decompile the APK file using apktool
    cmd := exec.Command("apktool", "d", apkPath)
    err := cmd.Run()
    if err != nil {
        fmt.Println("[-] Error decompiling APK file:", err)
        return
    }

    // Print status message
    fmt.Println("[*] Extracting firebase urls...")

    // Read the strings.xml file
    stringsXMLPath := strings.Replace(apkPath, ".apk", "", 1) + "/res/values/strings.xml"
    stringsXMLBytes, err := ioutil.ReadFile(stringsXMLPath)
    if err != nil {
        fmt.Println("[-] Error reading strings.xml file:", err)
        return
    }

    // Search for Firebase URLs using regular expressions
    firebaseURLRegex := regexp.MustCompile(`https://.*firebase.*`)
    firebaseURLs := firebaseURLRegex.FindAllString(string(stringsXMLBytes), -1)

    // Print the URLs found
    for _, url := range firebaseURLs {
        // Remove </string> from the URL
        url = strings.Replace(url, "</string>", "", -1)
        fmt.Printf("[+] Firebase URL found: %s\n", url)

        // Append /.json to the URL
        url = url + "/.json"

        // Send a GET request to the URL
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("[-] Error sending GET request:", err)
            return
        }
        defer resp.Body.Close()

        // Check if the response shows permission denied
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("[-] Error reading response body:", err)
            return
        }
        bodyString := string(bodyBytes)
        if strings.Contains(bodyString, "Permission denied") {
            fmt.Println(red + "[-] Firebase Permission denied" + reset)
        } else {
            fmt.Println(green + "[+] Firebase is accessible, go to this url: " + url + reset)
        }
    }

    // Print status message
    fmt.Println("[*] Removing extracted folder...")

    // Remove the extracted folder
    var cmdArgs []string
    if runtime.GOOS == "windows" {
        cmdArgs = []string{"/C", "rmdir", "/s", "/q", strings.Replace(apkPath, ".apk", "", 1)}
    } else {
        cmdArgs = []string{"-rf", strings.Replace(apkPath, ".apk", "", 1)}
    }
    cmd = exec.Command("rm", cmdArgs...)
    err = cmd.Run()
    if err != nil {
        fmt.Println("[-] Error removing extracted folder:", err)
        return
    }

}
