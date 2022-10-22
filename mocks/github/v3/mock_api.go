package v3

import (
	_ "embed"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Since the GitHub client doesn't have an interface to generate mocks against,
// and since the functions that we'll test use the GitHub client as the first of
// several actions, we can't simply make mock data and wave hands that the mock
// data is to be used for the rest of the actions.  Thus, we need some way to
// halfway mock the GitHub client, so that it doesn't depend on external HTTP
// traffic.  And in particular, we can accomplish that by making a server with
// httptest and have the GitHub client use it for its requests.

/* The API endpoint for repository information */
const pathRepository = "/repos"

/* The API endpoint regexp for a lookup of the list of releases in a repository */
const regexpReleaseList = "^" + pathRepository + "(/.+){2}/releases$"

/* The API endpoint regexp for a lookup of a single repository */
const regexpRepoGet = "^" + pathRepository + "(/.+){2}$"

/* The regexp to verify a client can accept a JSON response */
const regexpResponseJson = "application/([^\\s]+\\+)json"

const timeFormatISO8601 = "2006-01-02T15:04:05Z07:00"

/* Variable details of a particular release of a repository */
type releaseDetails struct {
	name         string
	tag          string
	isPreRelease bool
	isDraft      bool
	assets       []releaseAssetDetails
}

/* Variable details of a particular asset in a release of a repository */
type releaseAssetDetails struct {
	name      string
	updatedAt time.Time
}

/* Maps an owner, to a map of repository name, to an array of release details */
type mockMMApiStateType map[string]map[string][]releaseDetails

//go:embed asset_format.json
var jsonAssetformat string

//go:embed release_format.json
var jsonReleaseformat string

//go:embed repository_format.json
var jsonRepositoryformat string

//go:embed user_format.json
var jsonUserformat string

// Provides a test-integrated HTTP handler for single repository lookups
func makeHandleRepositoryGet(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if match, _ := regexp.MatchString(regexpRepoGet, pathRepository+"/user/foo"); !match {
			if t != nil {
				t.Errorf("Regex Match Fail")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else if match, _ := regexp.MatchString(regexpRepoGet, requestPath); !match {
			if t != nil {
				t.Errorf("Request path was not \"" + pathRepository + "/{owner}/{repo}\"")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else if match, _ := regexp.MatchString(regexpResponseJson, r.Header.Get("Accept")); !match {
			if t != nil {
				t.Errorf("Request does not accept JSON response")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			requestPathParts := strings.Split(r.URL.Path, "/")
			requestOwnerName := requestPathParts[len(requestPathParts)-2]
			requestRepoName := requestPathParts[len(requestPathParts)-1]

			repoFullName := requestOwnerName + "/" + requestRepoName
			repoHTMLURL := r.URL.Scheme + "://" + r.URL.Host + "/" + requestOwnerName + "/" + requestRepoName

			responseOwnerJson := formatUserJson(1, requestOwnerName)
			responseOrganizationJson := formatUserJson(1, requestOwnerName)
			response := formatRepositoryJson(
				1,
				requestRepoName,
				repoFullName,
				responseOwnerJson,
				responseOrganizationJson,
				repoHTMLURL)
			print(response + "\n")
			w.Write([]byte(response))
		}
	}
}

// Provides a test-integrated HTTP handler for looking up all releases of a repository
func makeHandleReleaseList(t *testing.T, state mockMMApiStateType) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if match, _ := regexp.MatchString(regexpReleaseList, "/repos/user/foo/releases"); !match {
			if t != nil {
				t.Errorf("Regex Match Fail")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		if match, _ := regexp.MatchString(regexpReleaseList, requestPath); !match {
			if t != nil {
				t.Errorf("Request path was not \"/repos/{owner}/{repo}/releases\"")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		if match, _ := regexp.MatchString(regexpResponseJson, r.Header.Get("Accept")); !match {
			if t != nil {
				t.Errorf("Request does not accept JSON response")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			requestPathParts := strings.Split(r.URL.Path, "/")
			requestOwnerName := requestPathParts[len(requestPathParts)-3]
			requestRepoName := requestPathParts[len(requestPathParts)-2]
			requestBaseURL := r.URL.Scheme + "://" + r.URL.Host
			var response string
			releases := state[requestOwnerName][requestRepoName]

			if len(releases) == 0 {
				response = "[]"
			} else {
				release := releases[0]
				releaseHTMLURLBase := strings.Join([]string{
					requestBaseURL,
					requestOwnerName,
					requestRepoName,
					"releases"}, "/")

				releaseAuthorJson := formatUserJson(1, requestOwnerName)

				asset_id := 1
				asset := release.assets[0]
				assetDownloadURL := strings.Join([]string{
					requestBaseURL,
					requestOwnerName,
					requestRepoName,
					"releases",
					"download",
					release.tag,
					asset.name}, "/")
				releaseAssetsJSON := formatAssetJson(
					int64(asset_id),
					asset.name,
					assetDownloadURL,
					asset.updatedAt,
					asset.updatedAt)
				for _, asset := range release.assets[1:] {
					asset_id += 1
					assetDownloadURL := strings.Join([]string{
						requestBaseURL,
						requestOwnerName,
						requestRepoName,
						"releases",
						"download",
						release.tag,
						asset.name}, "/")
					releaseAssetsJSON += ",\n" + formatAssetJson(
						int64(asset_id),
						asset.name,
						assetDownloadURL,
						asset.updatedAt,
						asset.updatedAt)
				}
				releaseAssetsJSON = "[\n" + indent(releaseAssetsJSON, "    ", false) + "\n]"

				response = formatReleaseJson(
					1,
					release.name,
					release.tag,
					releaseHTMLURLBase+"/"+release.tag,
					release.isDraft,
					release.isPreRelease,
					releaseAuthorJson,
					releaseAssetsJSON)

				for idxr, release := range releases[1:] {
					asset_id += 1
					asset := release.assets[0]
					assetDownloadURL := strings.Join([]string{
						requestBaseURL,
						requestOwnerName,
						requestRepoName,
						"releases",
						"download",
						release.tag,
						asset.name}, "/")
					releaseAssetsJSON := formatAssetJson(
						int64(asset_id),
						asset.name,
						assetDownloadURL,
						asset.updatedAt,
						asset.updatedAt)
					for _, asset := range release.assets[1:] {
						asset_id += 1
						assetDownloadURL := strings.Join([]string{
							requestBaseURL,
							requestOwnerName,
							requestRepoName,
							"releases",
							"download",
							release.tag,
							asset.name}, "/")
						releaseAssetsJSON += ",\n" + formatAssetJson(
							int64(asset_id),
							asset.name,
							assetDownloadURL,
							asset.updatedAt,
							asset.updatedAt)
					}
					releaseAssetsJSON = "[\n" + indent(releaseAssetsJSON, "    ", false) + "\n]"

					response += ",\n" + formatReleaseJson(
						int64(idxr+1),
						release.name,
						release.tag,
						releaseHTMLURLBase+"/"+release.tag,
						release.isDraft,
						release.isPreRelease,
						releaseAuthorJson,
						releaseAssetsJSON)
				}
			}
			response = "[\n" + indent(response, "    ", false) + "\n]"
			print(response + "\n")
			w.Write([]byte(response))
		}
	}
}

// Checks the request URL and invokes the appropriate handler for the path
func makeHandleAPIEndpoints(t *testing.T, state mockMMApiStateType) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Host == "" {
			r.URL.Host = r.Host
		}
		if r.URL.Scheme == "" {
			if r.TLS != nil {
				r.URL.Scheme = "https"
			} else {
				r.URL.Scheme = "http"
			}
		}
		requestPath := r.URL.Path
		if match, _ := regexp.MatchString(regexpReleaseList, requestPath); match {
			makeHandleReleaseList(t, state)(w, r)
		} else if match, _ := regexp.MatchString(regexpRepoGet, requestPath); match {
			makeHandleRepositoryGet(t)(w, r)
		} else {
			if t != nil {
				t.Errorf("Client requested an API path not mocked")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

// Indents each line of input with the provided indent
func indent(input string, indent string, skipFirstLine bool) string {
	result := strings.ReplaceAll(input, "\n", "\n"+indent)
	if !skipFirstLine {
		result = indent + result
	}
	return result
}

// Provides the string for true or false for use in JSON output
func jsonBool(input bool) string {
	if input {
		return "true"
	} else {
		return "false"
	}
}

// Populates an Asset JSON string with per-asset details
func formatAssetJson(id int64, name string, browserDownloadURL string, created time.Time, updated time.Time) string {
	return fmt.Sprintf(
		strings.Replace(jsonAssetformat, "\"MATCH_RELEASE_ID\"", fmt.Sprint(id), 1),
		browserDownloadURL,
		name,
		created.Format(timeFormatISO8601),
		updated.Format(timeFormatISO8601))
}

// Populates a Release JSON string with per-release details
func formatReleaseJson(id int64, name string, tagName string, htmlURL string, isDraft bool, isPrerelease bool, authorJSON string, assetsJSON string) string {
	result := fmt.Sprintf(
		jsonReleaseformat,
		htmlURL,
		tagName,
		name)
	result = strings.Replace(result, "\"MATCH_RELEASE_ID\"", fmt.Sprint(id), 1)
	result = strings.Replace(result, "\"MATCH_IS_DRAFT\"", jsonBool(isDraft), 1)
	result = strings.Replace(result, "\"MATCH_IS_PRERELEASE\"", jsonBool(isPrerelease), 1)
	result = strings.Replace(result, "\"MATCH_AUTHOR_OBJ\"", indent(authorJSON, "    ", true), 1)
	result = strings.Replace(result, "\"MATCH_ASSETS_ARR\"", indent(assetsJSON, "    ", true), 1)
	return result
}

// Populates a Repository JSON string with per-repository details
func formatRepositoryJson(id int64, name string, fullName string, ownerJSON string, organizationJSON string, htmlURL string) string {
	result := fmt.Sprintf(
		jsonRepositoryformat,
		name,
		fullName,
		htmlURL)
	result = strings.Replace(result, "\"MATCH_REPO_ID\"", fmt.Sprint(id), 1)
	result = strings.Replace(result, "\"MATCH_OWNER_OBJ\"", indent(ownerJSON, "    ", true), 1)
	result = strings.Replace(result, "\"MATCH_ORG_OBJ\"", indent(organizationJSON, "    ", true), 1)
	return result
}

// Populates an User JSON string with per-user details
func formatUserJson(id int64, login string) string {
	return fmt.Sprintf(
		strings.Replace(jsonUserformat, "\"MATCH_USER_ID\"", fmt.Sprint(id), 1),
		login)
}
