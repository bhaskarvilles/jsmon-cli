package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Helper to validate and handle API key
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

// Helper to handle domains (replacing repeated logic)
func parseDomains(domainStr string) []string {
	domains := strings.Split(domainStr, ",")
	for i, domain := range domains {
		domains[i] = strings.TrimSpace(domain)
	}
	return domains
}

// Central execution hub for most general flags
func executeCommand(
	scanFileId, uploadFile *string,
	viewurls *bool, viewurlsSize int,
	uploadUrl *string, headers stringSliceFlag,
	rescanDomainFlag *string, totalAnalysisDataFlag *bool,
	searchUrlsByDomainFlag *string, urlswithmultipleResponse *bool,
	viewEmails *string, getResultByJsmonId *string,
	reverseSearchResults *string, getResultByFileId *string,
	s3domains *string, ips *string, gql *string,
	domainUrl *string, apiPath *string, scanUrl *string,
	scanDomainFlag *string, wordsFlag *string, usageFlag *bool,
	getDomainsFlag *bool, getAllResults *string, size int,
	socialMediaUrls *string,
) {
	switch {
	case *scanFileId != "":
		scanFileEndpoint(*scanFileId)
	case *uploadFile != "":
		uploadFileEndpoint(*uploadFile, headers)
	case *viewurls:
		viewUrls(viewurlsSize)
	case *uploadUrl != "":
		uploadUrlEndpoint(*uploadUrl, headers)
	case *rescanDomainFlag != "":
		rescanDomain(*rescanDomainFlag)
	case *totalAnalysisDataFlag:
		totalAnalysisData()
	case *searchUrlsByDomainFlag != "":
		searchUrlsByDomain(*searchUrlsByDomainFlag)
	case *urlswithmultipleResponse:
		urlsmultipleResponse()
	case *viewEmails != "":
		domains := parseDomains(*viewEmails)
		getEmails(domains)
	case *getResultByJsmonId != "":
		getAutomationResultsByJsmonId(strings.TrimSpace(*getResultByJsmonId))
	case *reverseSearchResults != "":
		parts := strings.SplitN(*reverseSearchResults, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format for reverseSearchResults. Use field=value format.")
			return
		}
		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		getAutomationResultsByInput(field, value)
	case *getResultByFileId != "":
		getAutomationResultsByFileId(strings.TrimSpace(*getResultByFileId))
	case *s3domains != "":
		domains := parseDomains(*s3domains)
		getS3Domains(domains)
	case *ips != "":
		domains := parseDomains(*ips)
		getAllIps(domains)
	case *gql != "":
		domains := parseDomains(*gql)
		getGqlOps(domains)
	case *domainUrl != "":
		domains := parseDomains(*domainUrl)
		getDomainUrls(domains)
	case *apiPath != "":
		domains := parseDomains(*apiPath)
		getApiPaths(domains)
	case *scanUrl != "":
		rescanUrlEndpoint(*scanUrl)
	case *scanDomainFlag != "":
		words := []string{}
		if *wordsFlag != "" {
			words = strings.Split(*wordsFlag, ",")
		}
		automateScanDomain(*scanDomainFlag, words)
	case *usageFlag:
		callViewProfile()
	}
}

// Refactored usage output
func printUsage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Printf("  %s [flags]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "INPUT FLAGS:\n")
	options := []struct {
		flagName, description string
	}{
		{"-scanUrl <URL>", "Scan URL or scan ID for rescanning"},
		{"-uploadUrl <URL>", "Upload file or URL for scanning"},
		{"-scanFile <fileId>", "Provide file ID for scanning"},
		{"-uploadFile <string>", "Upload local file to the system"},
		{"-scanDomain <domain>", "Trigger scan automation on domain"},
		{"-getAutomationData <domain>", "Retrieve all automation data results"},
		{"-apikey <string>", "API Key for authentication"},
	}

	for _, opt := range options {
		fmt.Fprintf(os.Stderr, "  %-30s %s\n", opt.flagName, opt.description)
	}

	// Call flag's default usage for other options
	fmt.Fprintf(os.Stderr, "\nAdditional flags:\n")
	flag.PrintDefaults()
}

func main() {
	// Flag declarations
	scanUrl := flag.String("scanUrl", "", "URL or scan ID to rescan")
	uploadUrl := flag.String("uploadUrl", "", "URL to upload for scanning")
	apiKeyFlag := flag.String("apikey", "", "API key for authentication")
	scanFileId := flag.String("scanFile", "", "File ID for scanning")
	uploadFile := flag.String("uploadFile", "", "Path to local file to upload for scanning")
	getAllResults := flag.String("getAutomationData", "", "Get all automation results")
	size := flag.Int("size", 10000, "Number of results to fetch (default 10000)")
	getScannerResultsFlag := flag.Bool("getScannerData", false, "Get scanner results")
	viewurls := flag.Bool("urls", false, "View URLs")
	viewurlsSize := flag.Int("urlSize", 10, "Number of URLs to fetch")
	scanDomainFlag := flag.String("scanDomain", "", "Automate scan for domain")
	wordsFlag := flag.String("words", "", "Comma-separated list of words for scan")
	urlswithmultipleResponse := flag.Bool("changedUrls", false, "Check for URLs with multiple responses")
	getDomainsFlag := flag.Bool("getDomains", false, "Get all user domains")
	// Declaring custom header -H flag
	var headers stringSliceFlag
	flag.Var(&headers, "H", "Custom headers provided as 'Key: Value' (multiple allowed)")

	// Additional flags
	usageFlag := flag.Bool("usage", false, "View user profile")
	viewfiles := flag.Bool("getFiles", false, "View files")
	viewEmails := flag.String("getEmails", "", "View Emails for specified domains")
	s3domains := flag.String("getS3Domains", "", "Get S3 domains for specified domains")
	ips := flag.String("getIps", "", "Get IPs for specified domains")
	gql := flag.String("getGqlOps", "", "Get GraphQL operations")
	domainUrl := flag.String("getDomainUrls", "", "Get domain URLs")
	apiPath := flag.String("getApiPaths", "", "Get API paths for domain")
	totalAnalysisDataFlag := flag.Bool("totalAnalysisData", false, "Retrieve total count of overall analysis data")
	reverseSearchResults := flag.String("reverseSearchResults", "", "Specify search by input field (e.g., emails, domainname)")
	getResultByJsmonId := flag.String("getResultByJsmonId", "", "Get automation results by jsmon ID")
	getResultByFileId := flag.String("getResultByFileId", "", "Get file automation results based on file ID")
	rescanDomainFlag := flag.String("rescanDomain", "", "Rescan all URLs for a specific domain")

	// Custom usage function
	flag.Usage = printUsage
	flag.Parse()

	// Handle API key validation
	handleAPIKey(apiKeyFlag)

	// Check if no flag passed, default to `Usage`
	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *apiKeyFlag != "") {
		fmt.Println("No action specified. Use -h or --help for usage information.")
		flag.Usage()
		os.Exit(1)
	}

	// Execute commands based on parsed flags
	executeCommand(
		scanFileId, uploadFile,
		viewurls, *viewurlsSize,
		uploadUrl, headers,
		rescanDomainFlag, totalAnalysisDataFlag,
		scanDomainFlag, urlswithmultipleResponse,
		viewEmails, getResultByJsmonId, reverseSearchResults,
		getResultByFileId, s3domains, ips,
		gql, domainUrl, apiPath,
		scanDomainFlag, wordsFlag, usageFlag, getDomainsFlag,
		getAllResults, *size, s3domains,
	)
}
