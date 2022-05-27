package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/internals/types"
	"net/url"
	"os"
	"os/exec"
	"runtime"
)

var LogsText []string
var LinksToCheck = GetLinksToCheck()
var ProgressBarValue = 0.0

func GetLinksToCheck() []types.LinkToCheck {
	links := make([]types.LinkToCheck, 0)
	sources := GetEnv("URLS_TO_CHECK", "[{\"url\": \"https://aejuice.com\", \"name\": \"Website and API\"},{\"url\": \"https://nyc3.digitaloceanspaces.com/aejuice/update/updater/version.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=SSK2B4AURYVVYMUF75K3%2F20220524%2Fnyc3%2Fs3%2Faws4_request&X-Amz-Date=20220524T214003Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=8a9d408ac988642ffeb3021eb5447f3a423e96ea335aa6ad46ff63a66cb0a83d\", \"name\": \"Possibility to download files\"}]")
	_ = json.Unmarshal([]byte(sources), &links)

	return links
}

func GetEnv(envName string, fallback string) string {
	value := os.Getenv(envName)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func IncrementProgressValue() {
	floatLength := 1.00 / float64(len(LinksToCheck))
	ProgressBarValue = ProgressBarValue + floatLength
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func GetTracerouteFunction() string {
	if IsWindows() {
		return "tracert"
	}

	return "traceroute"
}

func GetMaxHopsArg() string {
	if IsWindows() {
		return "-h 30"
	}

	return "-m30"
}

func Traceroute(href string) error {
	tracerouteFunction := GetTracerouteFunction()
	currentUrl, _ := url.Parse(href)
	hopsArg := GetMaxHopsArg()
	cm := exec.Command(tracerouteFunction, hopsArg, currentUrl.Host)
	stdout, _ := cm.StdoutPipe()
	stderr, _ := cm.StderrPipe()
	err := cm.Start()
	if err != nil {
		LogsText = append(LogsText, fmt.Sprintf("%s", err)+"\n")
		return err
	}
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			LogsText = append(LogsText, stderrScanner.Text()+"\n")
		}
	}()
	go func() {
		for stdoutScanner.Scan() {
			LogsText = append(LogsText, stdoutScanner.Text()+"\n")
		}
	}()
	err = cm.Wait()
	if err != nil {
		LogsText = append(LogsText, fmt.Sprintf("%s", err)+"\n")
		return err
	}

	return nil
}

func SaveLogs() {
	var bytesArr []byte
	for i := 0; i < len(LogsText); i++ {
		b := []byte(LogsText[i])
		for j := 0; j < len(b); j++ {
			bytesArr = append(bytesArr, b[j])
		}
	}

	err := ioutil.WriteFile("network-diagnostic-tool.log", bytesArr, 0644)
	if err != nil {
		panic(err)
	}
}
