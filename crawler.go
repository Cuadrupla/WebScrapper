/**
@package CS112/Crawler

@author Vrinceanu Radu-Tudor <radu-tudor.vrinceanu@s.unibuc.ro>
@author Ispas Jany-Gabriel <jany-gabriel.ispas@s.unibuc.ro>
@author Gheorghe Liviu-Ionut <liviu-ionut.gheorghe@s.unibuc.ro>

THE IDEA OF THE REGEX CRAWLER IS TO CONSTRUCT A PARSE TREE WITH LINKS CRAWLED FROM A SPECIFIC DOMAIN GIVEN
AS AN INPUT. FOR THAT WE WILL ALSO NEED TO EXTRACT INFORMATION ABOUT THE LINK TO CHECK WHETHER IS AN ACCESSIBLE
ROUTE TO GO IN THE RECURSION'S DEPTH OR STOP RIGHT THERE. THE ALGORITHM REPRESENTS A DEPTH-FIRST SEARCH PROCEDURE.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

var countUrls = 0
var visitedUrls = make(map[string]bool)
var uniqueUrls = make(map[string]bool)

func crawilingRoutine(beginUrl *string, recursionDepth *int, currentDepth int) {
	if recursionDepth == nil || beginUrl == nil || *recursionDepth == currentDepth || visitedUrls[*beginUrl] {
		return
	}
	visitedUrls[*beginUrl] = true

	req, err := http.NewRequest("GET", *beginUrl, nil)
	if err != nil {
		fmt.Println(err, *recursionDepth, currentDepth)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err, *recursionDepth, currentDepth)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%s (recursion depth: %d, current depth: %d)\n", *beginUrl, *recursionDepth, currentDepth)
		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err, *recursionDepth, currentDepth)
		return
	}

	r := regexp.MustCompile(`(([\w]+:)//)(([\d\w]|%[a-fA-f\d]{2,2})+(:([\d\w]|%[a-fA-f\d]{2,2})+)?@)?([\d\w][-\d\w]{0,253}[\d\w]\.)+[\w]{2,63}(:[\d]+)?(/([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)*(\?(&?([-+_~.\d\w]|%[a-fA-f\d]{2,2})=?)*)?(#([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)?`)
	matches := r.FindAllStringSubmatch(string(body), -1)

	countUrls += len(matches)
	for _, v := range matches {
		uniqueUrls[v[0]] = true
		crawilingRoutine(&v[0], recursionDepth, currentDepth+1)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ######################################## THE CLI Arguments ########################################
	beginUrl := flag.String("beginUrl", "", "This is used to set the beginning url for the crawler.")
	recursionDepth := flag.Int("depth", 0, "This is the maximum recursion depth for the crawler.")

	requiredArgs := []string{"beginUrl", "depth"}
	flag.Parse()

	checker := make(map[string]bool)
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != "" && f.Value.String() != "0" {
			checker[f.Name] = true
		}
	})
	for _, req := range requiredArgs {
		if !checker[req] {
			_, err := fmt.Fprintf(os.Stderr, "[CRAWLER]: missing required --%s argument/flag\n", req)
			if err != nil {
				return
			}
			os.Exit(2)
		}
	}

	crawilingRoutine(beginUrl, recursionDepth, 0)
	fmt.Printf("[CRAWLER RESULTS]:\nUsed link: %s (recursion depth: %d)\nFound (total: %d | unique: %d | visited: %d) URLS with crawler\n",
		*beginUrl, *recursionDepth, countUrls, len(uniqueUrls), len(visitedUrls))
}
