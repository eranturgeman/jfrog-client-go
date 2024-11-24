package utils

import (
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	XraySuffix                        = "/xray/"
	xscSuffix                         = "/xsc/"
	apiV1Suffix                       = "api/v1"
	XscInXraySuffix                   = apiV1Suffix + xscSuffix
	MinXrayVersionXscTransitionToXray = "3.108.0"
)

// From Xray version 3.108.0, XSC is transitioning to Xray as inner service. This function will return compatible URL.
func XrayUrlToXscUrl(xrayUrl, xrayVersion string) string {
	if !IsXscXrayInnerService(xrayVersion) {
		log.Debug(fmt.Sprintf("Xray version is lower than %s, XSC is not an inner service in Xray.", MinXrayVersionXscTransitionToXray))
		return strings.Replace(xrayUrl, XraySuffix, xscSuffix, 1) + apiV1Suffix + "/"
	}
	// Newer versions of Xray will have XSC as an inner service.
	return xrayUrl + XscInXraySuffix
}

func IsXscXrayInnerService(xrayVersion string) bool {
	if err := utils.ValidateMinimumVersion(utils.Xray, xrayVersion, MinXrayVersionXscTransitionToXray); err != nil {
		return false
	}
	return true
}

func GetGitRepoUrlKey(gitRepoUrl string) string {
	if len(gitRepoUrl) == 0 {
		// No git context was provided
		return ""
	}
	if !strings.HasSuffix(gitRepoUrl, ".git") {
		// Append .git to the URL if not included
		gitRepoUrl += ".git"
	}
	// Remove the Http/s protocol from the URL
	if strings.HasPrefix(gitRepoUrl, "http") {
		return strings.TrimPrefix(strings.TrimPrefix(gitRepoUrl, "https://"), "http://")
	}
	// Remove the SSH protocol from the URL
	if strings.Contains(gitRepoUrl, "git@") {
		return strings.Replace(strings.TrimPrefix(strings.TrimPrefix(gitRepoUrl, "ssh://"), "git@"), ":", "/", 1)
	}
	return gitRepoUrl
}
