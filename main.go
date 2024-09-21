package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Custom string slice flag type to handle multiple `-H` headers
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Centralized function to load and validate the API key
func handleAPIKey(apiKeyFlag *string) {
	if *apiKeyFlag != "" {
		setAPIKey(*apiKeyFlag)
	} else {
		err := loadAPIKey()
		if err != nil {
			fmt.Println("Error loading API key:", err)
			fmt.Println("Please provide an API key using the -apikey flag.")
			os.Exit(1)
		}
	}
}

// Centralized domain splitting and trimming logic
func parseDomains(domainStr string) []string {
	domains := strings.Split(domainStr, ",")
	for i, domain := range domains {
		domains[i] = strings.TrimSpace(domain)
	}
	return domains
}

// Command Execution Functions
func executeCommands(cmdFlag string, headers stringSliceFlag, size int) {
	switch cmdFlag {
	case "scanFile":
		scanFileEndpoint(*scanFileId)
	case "uploadFile":
		uploadFileEndpoint(*uploadFile, headers)
	case "viewURLs":
		viewUrls(size)
	case "viewFiles":
		viewFiles()
	case "rescanDomain":
		rescanDomain(*rescanDomainFlag)
	case "totalAnalysisData":
		totalAnalysisData()
	case "searchUrlsByDomain":
		searchUrlsByDomain(*searchUrlsByDomainFlag)
	case "changedUrls":
		urlsmultipleResponse()
	}
}

// Setup function to initialize flags
func getFlags() (stringSliceFlag, *string, *string, *string, *string, *bool, *int, *bool, *int, *bool) {
	var headers stringSliceFlag
	scanUrl := flag.String("scanUrl", "", "URL or scan ID to rescan")
	uploadUrl := flag.String("uploadUrl", "", "URL to upload for scanning")
	apiKeyFlag := flag.String("apikey", "", "API key for authentication")
	scanFileId := flag.String("scanFile", "", "File ID to scan")
	uploadFile := flag.String("uploadFile", "", "Path to local file to upload for scanning.")
	getAllResults := flag.String("getAutomationData", "", "Get all automation results")
	sizeFlag := flag.Int("size", 10000, "Number of results to fetch (default 10000)")
	getScannerResultsFlag := flag.Bool("getScannerData", false, "Get scanner results")
	viewURLsFlag := flag.Bool("urls", false, "View all URLs")
	viewURLsSize := flag.Int("urlSize", 10, "Number of URLs to fetch")
	flag.Var(&headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")
	return headers, scanUrl, uploadUrl, scanFileId, uploadFile, getScannerResultsFlag, sizeFlag, viewURLsFlag, viewURLsSize
}

func usageFunc(programName string) {
	fmt.Printf("Usage of %s:\n", programName)
	fmt.Println("  [flags]")

	fmt.Fprintln(os.Stderr, "Flags:")
	fmt.Fprintf(os.Stderr, "  -apikey <XXXXXX-XXXX-XXXX-XXXX-XXXXXX>          API key for authentication\n")
	fmt.Fprintf(os.Stderr, "  -scanFile <fileId>         File ID to scan\n")
	fmt.Fprintf(os.Stderr, "  -uploadFile <filePath>     Path to local file to upload for scanning\n")
	fmt.Fprintf(os.Stderr, "  -urls                    View all URLs\n")
	fmt.Fprintf(os.Stderr, "  -urlSize <int>            Number of URLs to fetch (default 10)\n")
	fmt.Fprintln(os.Stderr, "\nCRON JOB FLAGS:")
	fmt.Fprintln(os.Stderr, "  -notifications <string>    Set cronjob notification channel.")
	fmt.Fprintln(os.Stderr, "  -time <int64>              Set cronjob time.")
	fmt.Fprintln(os.Stderr, "\nMORE OPTIONS:")
	fmt.Fprintln(os.Stderr, "  -H <custom header>         Custom headers for requests (can repeat)")
	flag.PrintDefaults()
}

func main() {
	headers, scanUrl, _, scanFileId, uploadFile, getScannerResultsFlag, sizeFlag, viewURLsFlag, viewURLsSize := getFlags()

	// Custom usage when no arguments
	flag.Usage = func() {
		usageFunc(os.Args[0])
	}

	// Parse command-line arguments
	flag.Parse()

	// Handle API key validation
	apiKeyFlag := flag.String("apikey", "", "API key for authentication")
	handleAPIKey(apiKeyFlag)

	// Check if no arguments provided or only an API key
	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *apiKeyFlag != "") {
		fmt.Println("No action specified. Use -h or --help for usage information.")
		flag.Usage()
		os.Exit(1)
	}

	// Execute the appropriate function based on the provided flag
	if *scanFileId != "" {
		executeCommands("scanFile", headers, *sizeFlag)
	} else if *uploadFile != "" {
		executeCommands("uploadFile", headers, *sizeFlag)
	} else if *viewURLsFlag {
		executeCommands("viewURLs", headers, *viewURLsSize)
	} else if *getScannerResultsFlag {
		executeCommands("getScannerResults", headers, *sizeFlag)
	} else {
		// No valid action specified
		fmt.Println("No valid action specified.")
		flag.Usage()
		os.Exit(1)
	}
}

// type Args struct {
// 	Cron             string
// 	CronNotification string
// 	CronTime         int64
// 	CronType         string
// }

// func parseArgs() Args {
// 	//CRON JOB FLAGS ->
// 	cron := flag.String("cron", "", "Set cronjob.")
// 	cronNotification := flag.String("notifications", "", "Set cronjob notification.")
// 	cronTime := flag.Int64("time", 0, "Set cronjob time.")
// 	cronType := flag.String("type", "", "Set type of cronjob.")

// 	flag.Parse()

// 	return Args{
// 		Cron:             *cron,
// 		CronNotification: *cronNotification,
// 		CronTime:         *cronTime,
// 		CronType:         *cronType,
// 	}
// }
