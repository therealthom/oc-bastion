package main

import (
	"fmt"
	"github.com/burntcarrot/ricecake"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	fileName    string
	fullURLFile string
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
		// Check the flag values
		homedir = checkHomeDir(homedir)
		checkVersion(version)
		checkPullSecret(pullsecret)
		checkBaseDomain(basedomain)
		checkClusterName(clustername)

		// Creating dirs

		fmt.Printf("Creating dirs in %s\n", homedir)
		if err := os.MkdirAll(homedir+"/ocpinstall/resources", os.ModePerm); err != nil {
			log.Fatal(err)
		}
		if err := os.MkdirAll(homedir+"/ocpinstall/install", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		// Download files
		customPath, ocpInstallDir := downloadFile(homedir, "client", version)
		fmt.Println(ocpInstallDir)
		unTarFile(customPath)
		/*moveFile("oc", ocpInstallDir)
		customPath, ocpInstallDir = downloadFile(homedir, "install", version)
		unTarFile(customPath)
		moveFile("openshift-install", ocpInstallDir)
		removeFiles()*/
		/*fmt.Println("I am the root command!")
		fmt.Printf("-h, --homedir flag value: %s\n", homedir)
		fmt.Printf("-v, --version flag value: %s\n", version)
		fmt.Printf("-p, --pullsecret flag value: %s\n", pullsecret)
		fmt.Printf("-b, --basedomain flag value: %s\n", basedomain)
		fmt.Printf("-c, --clustername flag value: %s\n", clustername)*/
		return nil
	})

	// Run the CLI.
	err := cli.Run()
	if err != nil {
		log.Fatalf("failed to run oc-bastion; err: %v", err)
	}
}

func checkClusterName(clustername string) {
	if clustername == "" {
		log.Fatal("-c, --clustername is required")
	}
}

func checkBaseDomain(basedomain string) {
	if basedomain == "" {
		log.Fatal("-b, --basedomain is required")
	}
}

func checkHomeDir(homedir string) string {
	if homedir == "" {
		var userhomedir string
		userhomedir, err := os.UserHomeDir()
		fmt.Printf("userhomedir: %s\n", userhomedir)
		if err != nil {
			log.Fatal(err)
		}
		return userhomedir
	}
	return homedir
}

func checkVersion(version string) {
	ocVersions := []string{"4.9", "4.10", "4.11"}
	if version == "" || !slices.Contains(ocVersions, version) {
		log.Fatal("-v, --version is required, it must be one of these: ", ocVersions)
	}
}

func checkPullSecret(pullsecret string) {
	if pullsecret == "" {
		log.Fatal("-p, --pullsecret is required")
	}
}

func downloadFile(homeDir string, fileType string, version string) (string, string) {
	fullURLFile = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable-" +
		version + "/openshift-" + fileType + "-linux.tar.gz"
	//version + "/openshift-" + fileType + "-mac-arm64.tar.gz"

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	// Create blank file
	var ocpInstallDir = homeDir + "/ocpinstall"
	var customPath = ocpInstallDir + "/resources/" + fileName

	file, err := os.Create(customPath)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	defer file.Close()
	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
	fmt.Printf("customPath: %s\n", customPath)
	fmt.Printf("ocpInstallDir: %s\n", ocpInstallDir)
	return customPath, ocpInstallDir
}

func unTarFile(customPath string) {
	cmd := exec.Command("tar", "zxvf", customPath)
	erro := cmd.Run()
	if erro != nil {
		log.Fatal(erro)
	}
}
func removeFiles() {
	cmd := exec.Command("rm", "kubectl")
	erro := cmd.Run()
	if erro != nil {
		log.Fatal(erro)
	}
	cmd = exec.Command("rm", "README.md")
	erro = cmd.Run()
	if erro != nil {
		log.Fatal(erro)
	}
}

func moveFile(fileName string, ocpInstallDir string) {
	cmd := exec.Command("mv", fileName, ocpInstallDir)
	erro := cmd.Run()
	if erro != nil {
		log.Fatal(erro)
	}
}
