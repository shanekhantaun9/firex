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

    "github.com/fatih/color"
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

    // Check if any URLs were found
    if len(firebaseURLs) == 0 {
        color.Red("[-] Firebase url is not found")
        return
    }

    // Print the URLs found
    for _, url := range firebaseURLs {
        // Remove </string> from the URL
        url = strings.Replace(url, "</string>", "", -1)
        url = strings.TrimRight(url, "\r")
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
            color.Red("[-] Firebase Permission denied")
        } else if strings.Contains(bodyString, "has been deactivated.") {
            color.Red("[-] Firebase has been deactivated.")
        } else {
            color.Green("[+] Firebase is accessible, go to this url: %s", url)
        }
    }

    // Print status message
    fmt.Println("[*] Removing extracted folder...")

    // Remove the extracted folder
    var cmdArgs []string
    if runtime.GOOS == "windows" {
        cmdArgs = []string{"/C", "rmdir", "/s", "/q", strings.Replace(apkPath, ".apk", "", 1)}
        cmd = exec.Command("cmd", cmdArgs...)
    } else {
        cmdArgs = []string{"-rf", strings.Replace(apkPath, ".apk", "", 1)}
        cmd = exec.Command("rm", cmdArgs...)
    }
    err = cmd.Run()
    if err != nil {
        fmt.Println("[-] Error removing extracted folder:", err)
        return
    }
}
