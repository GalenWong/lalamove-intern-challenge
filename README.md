# Technical Challenge - Lalamove

## Afterthoughts
It was a really fun challenge. I learnt a lot from writing codes and reading documentations.

I included ```data.txt``` which contains some test strings that I used to test the main routine. 

They include: 
- correctly formatted repo
- non-existent repo
- incorrect formatting
- repo that has a large number of releases (node.js) with a low minVer
- repo that does not have releases
- repo with tagName that is invalid to semver.New
- empty line


I also implemented extra test cases for the main routine, with extra condition checks. Namely, the cases where the test fails are different length of result and expected, and expected of length 0.

### Difficulties
- Mastering Golang syntax
- Error handling
- Get as much releases as possible without breaking the git API rate limit

The way I handled the api rate limit issue is to set a upper limit on the number of calls. Theoretically, if the there is no such upper limit, the routine will work perfectly fine even if it has to check until version 0.0.1, since the routine itself will limit calls by checking whether the current call contains version lower or equal to the minVer. If so, stop querying. 


## Tech/Infrastructure Intern

This exercise is designed to assess some basic knowledge of writing applications. These include, but are not limited to:
- Writing code in Golang
- Use the standard library for solving various algorithmic tasks
- Implement well-documented third-party libraries

Bonus points for:
- Thorough test cases
- Separation of logic
- optimising the amount of API calls

# Preamble
We use a lot of Open Source software at Lalamove, and we want to be able to track when new versions of all of these applications are released. The open source community has mostly settled on using Github and its releases feature to publish releases and are also mostly using Semantic Versioning as their versioning structure.

## The Challenge
We want you to write a simple application that gives us the highest patch version of every release between a minimum version and the highest released version.
It should be written in Go, reads the Github Releases list, uses SemVer for comparison and takes a path to a file as its first argument when executed. It reads this file, which is in the format of:
```
repository,min_version
kubernetes/kubernetes,1.8.0
prometheus/prometheus,2.2.0
```
and it should produce output to stdout in the format of:
```
latest versions of kubernetes/kubernetes: [1.10.1 1.9.6 1.8.11]
latest versions of prometheus/prometheus: [2.2.1]
```

In this repository you will find two go files, one main.go which contains a skeleton for you to start developing your application, and one main_test.go, which contains some test cases that should pass when the application is ready. You run the tests by writing `go test`, and there are some cases that aren't tested for - can you figure them out?

We will run your application through a number of test cases not mentioned here and compare the output with an application that produces the correct output.
