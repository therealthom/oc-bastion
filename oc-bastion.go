package main

import (
	"fmt"
	"log"

	"github.com/burntcarrot/ricecake"
)

func main() {
	// Create a new CLI.
	cli := ricecake.NewCLI("oc-bastion", "OC Bastion Command", "v0.1")

	// Set long description for the CLI.
	cli.LongDescription("This command helps to provision everything " +
		"needed on the bastion node in the Openshift UFI installation process.")

	// -f, --file flag.
	var homedir string
	var version string
	var pullsecret string
	var basedomain string
	var clustername string
	cli.StringFlagP("homedir", "h", "Filename", &homedir)
	cli.StringFlagP("version", "v", "OC Version", &version)
	cli.StringFlagP("pullsecret", "p", "Pull Secret", &pullsecret)
	cli.StringFlagP("basedomain", "b", "Base Domain", &basedomain)
	cli.StringFlagP("clustername", "c", "Cluster Name", &clustername)

	// Set the action for the CLI.
	cli.Action(func() error {
		fmt.Println("I am the root command!")
		fmt.Printf("-h, --homedir flag value: %s\n", homedir)
		fmt.Printf("-v, --version flag value: %s\n", version)
		fmt.Printf("-p, --pullsecret flag value: %s\n", pullsecret)
		fmt.Printf("-b, --basedomain flag value: %s\n", basedomain)
		fmt.Printf("-c, --clustername flag value: %s\n", clustername)
		return nil
	})

	// Run the CLI.
	err := cli.Run()
	if err != nil {
		log.Fatalf("failed to run oc-bastion; err: %v", err)
	}
}
