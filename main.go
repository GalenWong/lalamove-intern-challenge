package main

import (
    "context"
    "fmt"
    "os"    // for reading command-line arguments
    "bufio" // for reading file by lines
    "errors"// for error handling
    "sort"  // for sorting
    "github.com/coreos/go-semver/semver"
    "github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
    var versionSlice []*semver.Version

    if minVersion == nil {
        return versionSlice
    }
    if len(releases) == 0 { // empty slice
        return versionSlice
    }
    for _, release := range releases {       // adding all versions >= minVersion into versionSlice
        if !(release.LessThan(*minVersion)) {
            versionSlice = append(versionSlice, release)
        }
    }

    if len(versionSlice) == 0 {
        return versionSlice
    }

    DescendingSort(versionSlice)    // sorting verssionSlice into descending order

    var prevMaxMajor = versionSlice[0].Major
    var prevMaxMinor = versionSlice[0].Minor
    var result []*semver.Version

    result = append(result, versionSlice[0])

    for _, release := range versionSlice {
        if release.Major < prevMaxMajor {
            prevMaxMinor = release.Minor
            prevMaxMajor = release.Major
            result = append(result, release)
            continue
        }
        if release.Minor < prevMaxMinor {
            prevMaxMinor = release.Minor
            result = append(result, release)
        }
    }

    return result
}


// sort required pre-defined functions
type Versions []*semver.Version

func (s Versions) Len() int {
    return len(s)
}

func (s Versions) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s Versions) Less(i, j int) bool {
    return s[i].LessThan(*s[j])
}

// Sort sorts the given slice of Version
func DescendingSort(versions []*semver.Version) {
    sort.Sort(sort.Reverse(Versions(versions)))
}

// separate a string into author, repo, minVer
func ProcessString(str string) (author string, repo string, minVer *semver.Version, err error) {
    var i int
    runes := []rune(str)
    var tokenPos int

    // default error value to be returned
    author  = ""
    repo    = ""
    minVer  = nil
    err     = errors.New(fmt.Sprintf("invalid string: \"%s\"", str))

    if len(str) == 0 {
        return
    }
    for i = 0; i < len(str); i++ {
        if(str[i] == '/'){
            author = string(runes[0:i])
            tokenPos = i
            break
        }
    }
    if i == len(str) {
        return // error value
    }
    for ; i < len(str); i++ {
        if(str[i] == ','){
            repo = string(runes[tokenPos+1:i])
            tokenPos = i
            break
        } 
    }
    if i == len(str) {
        return // error value
    }
    minVer, err = GetVersion(string(runes[tokenPos+1:]))
    if err != nil {
        return
    }

    return author, repo, minVer, nil
}

func GetVersion(versionString string) (ver *semver.Version, err error){
    defer func(){
        if recover() != nil{
            ver = nil
            err = errors.New(fmt.Sprintf("Invalid version string: \"%s\"", versionString))
        }
    }()
    if versionString[0] == 'v' {
        versionString = versionString[1:]
    }
    ver = semver.New(versionString)
    err = nil
    return
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {

    if len(os.Args) < 2 {
        panic(errors.New("No path provided"))
    }

    pathToFile := os.Args[1]

    file, err := os.Open(pathToFile) 

    if err != nil{
        panic(err)
    }

    defer file.Close()

    // Github
    client := github.NewClient(nil)
    ctx := context.Background()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {    // for each line 
        author, repo, minVer, err := ProcessString(scanner.Text());
        if err != nil {
            fmt.Println(err)
            continue
        }
        var minVerIsReached = false
        var allReleases []*semver.Version
        var pageNum = 1
        var isInvalidRepo = false

        for !minVerIsReached {  // this loop attempts to find the minVersion by iterating through pages
            opt := &github.ListOptions{Page: pageNum, PerPage: 10}

            releases, _, err := client.Repositories.ListReleases(ctx, author, repo, opt)
            if err != nil {
                fmt.Println(err)
                isInvalidRepo = true
                switch err.(type){
                case *github.ErrorResponse:
                    fmt.Printf("Invalid repo: %s/%s\n",author, repo)
                }
                break
            }
            if len(releases) == 0 {  // may not have any releases
                break
            }

            for _, release := range releases {
                versionString := *release.TagName
                version, err := GetVersion(versionString)
                if err != nil{
                    continue
                }
                allReleases = append(allReleases, version)
                if version.LessThan(*minVer) || version.Equal(*minVer) {    
                    minVerIsReached = true  // if find version equal or less than minVer, stop querying
                }
            }

            if pageNum > 7 {    // since github impose a rate call limit, for each repo, we read at most 7 pages
                break
            }
            pageNum++
        }

        if isInvalidRepo && pageNum==1 {
            continue
        }
        versionSlice := LatestVersions(allReleases, minVer)
        fmt.Printf("latest versions of %s/%s: %s\n", author, repo, versionSlice)
    }

    file.Close()
}
