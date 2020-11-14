package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/webview/webview"
)

// Version contains the version of the agent being used
const Version = "v1.0.0-dev"

var config AgentConfig

var enrollFlg = flag.Bool("enroll", false, "todo")

func main() {
	flag.Parse()
	// TODO: Version command + Help command

	var filePath = "./cmd/agent/interface/index.html"
	if *enrollFlg == true {
		filePath = "./cmd/agent/interface/enroll.html"
	} else {
		var err error
		config, err = LoadConfig(AgentConfigDefaultPath)
		if err != nil {
			log.Printf("Error loading configuration: %s\n", err)
			os.Exit(1)
		}
	}

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("MDM")
	w.SetSize(400, 550, webview.HintFixed)

	w.Bind("login", func(upn string, password string) {
		w.Terminate()
		log.Println("Enrolling user " + upn)
		if err := enroll(upn, password); err != nil {
			log.Fatalln(err) // TODO: Show in GUI
			return
		}
		log.Println("User enrolled!")
	})

	w.Bind("info", func() {
		fmt.Println("info")
		w.Eval(`document.getElementById("server-url").textContent = "` + config.Server + `";`)
	})

	w.Bind("sync", func() {
		fmt.Println("Syncing")

		// TODO: Update UI
	})

	enrollFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error loading enrollment interface: %s\n", err)
		os.Exit(1)
	}
	w.Navigate("data:text/html," + string(enrollFile))
	w.Run()
}
