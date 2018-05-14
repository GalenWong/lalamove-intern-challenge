package main

import (
	"context"
	"fmt"
	"os"	// for reading command-line arguments
	"bufio"	// for reading file by lines
	"errors"// for error handling
	"sort"	// for sorting
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	if minVersion == nil {
		return versionSlice
	}
	if len(releases) <= 0 {	// empty slice
		return versionSlice
	}
	for _, release := range releases{	// adding all versions >= minVersion into versionSlice
		if !(release.LessThan(*minVersion)) {
			versionSlice = append(versionSlice, release)
		}
	}

	DescendingSort(versionSlice)	// sorting verssionSlice into descending order

	var prevMaxMinor = versionSlice[0].Minor
	var result []*semver.Version

	result = append(result, versionSlice[0])

	for _, release := range versionSlice{
		if (release.Minor < prevMaxMinor) {
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
func ProcessString(str string) (author string, repo string, minVer *semver.Version, err error){
	var i int
	runes := []rune(str)
	var tokenPos int

	// default error value to be returned
	author 	= ""
	repo 	= ""
	minVer 	= nil
	err 	= errors.New(fmt.Sprintf("invalid string: \"%s\"", str))

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

	defer func(){
		if recover() != nil {
			// just return error value
		}
	}()
	
	minVer = semver.New(string(runes[tokenPos+1:]))

	return author, repo, minVer, nil
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	pathToFile := os.Args[1]

	file, err := os.Open(pathToFile) 

	if err != nil{
		panic(err)
	}

	defer file.Close()

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {	// for each line 
		author, repo, minVer, err := ProcessString(scanner.Text());
		if err != nil {
			fmt.Println(err)
			continue
		}
		releases, _, err := client.Repositories.ListReleases(ctx, author, repo, opt)
		if err != nil {
			fmt.Println(err)
			continue
		}

		allReleases := make([]*semver.Version, len(releases))

		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}

		fmt.Println("here with string :", scanner.Text())

		versionSlice := LatestVersions(allReleases, minVer)
		fmt.Printf("latest versions of %s/%s: %s\n", author, repo, versionSlice)
	}
}
