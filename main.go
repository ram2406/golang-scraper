package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"encoding/json"

	"modules/transform_html"

	"github.com/akamensky/argparse"
)


func wlk_(
	url string, rule []* transform_html.ParserTransfromRule,
) {
	
	ret_data := map[string]any{}

	req, err := http.NewRequest(`GET`, url, nil)
	req.Header.Set(`user-agent`, `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rvsca:109.0) Gecko/20100101 Firefox/113.0`)
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode >= 400 {
		fmt.Printf(`%d`, resp.StatusCode)
		return
	}
	fmt.Printf("\n%s", err)
	bta, err := io.ReadAll(resp.Body)
	fmt.Printf("\n%s", err)
	
	err = transform_html.TransformHtmlText(&ret_data, string( bta ), rule)
	//lst := ret_data[`menu_items`].(*list.List)
	
	fmt.Printf("\nerr %s", err)
	bta2, err := json.Marshal(ret_data)
	fmt.Printf(`\n%s`, bta2)

}

func wlk_2(walker *Walker, url string, rule []* transform_html.ParserTransfromRule,) {
	
	ret_data, err := walker.extract_data(url, rule)
	fmt.Printf("\nerr %s", err)
	bta2, err := json.Marshal(ret_data)
	fmt.Printf(`\n%s`, bta2)
}

func main() {
	start_time := time.Now()
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	etl_config_path := parser.String("p", "etl_config_path", &argparse.Options{Required: true, Help: "Etl config file location"})
	source_name := parser.String("s", "source_name", &argparse.Options{Required: true, Help: "Source name from config file"})
	filter_url := parser.String("f", "filter_url", &argparse.Options{Required: true, Help: "Filter url pattern"})
	output_file_path := parser.String("o", "output_file_path", &argparse.Options{Required: true, Help: "Output file with JSON format"})
	begin_page := parser.Int("b", "begin_page", &argparse.Options{Required: true, Help: "Scraping will start from begin page"})
	end_page := parser.Int("e", "end_page", &argparse.Options{Required: true, Help: "Scraping will stop on end page"})
	max_os_threads := parser.Int("m", "max_os_threads", &argparse.Options{Required: false, Help: "Set max OS threads"})
	goroutine_disable := parser.Flag("d", "goroutine_disable", &argparse.Options{Required: false, Help: "Work without green thread and go-routines", Default: false})
	
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}
	if max_os_threads != nil {
		runtime.GOMAXPROCS(*max_os_threads)
	}
	// Finally print the collected string
	walker, err := (&Walker{}).Init(*source_name, *etl_config_path)
	if err != nil {
		panic(err)
	}
	output := map[uint]any{}
	walk_fun := walker.Walk
	if *goroutine_disable {
		walk_fun = walker.WalkSync
	}
	err = walk_fun(*filter_url, uint(*begin_page), uint(*end_page), func(tm []*transform_html.TransformMap, num uint) {
		output[num] = tm
		fmt.Printf("\n Handled menu page %v count %v", num, len(tm))
	})
	if err != nil {
		panic(err)
	}

	data := OutputJson{*begin_page, *end_page, *etl_config_path, *source_name, output}
	json_str, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	
	err = ioutil.WriteFile(*output_file_path, json_str, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n Done in %v seconds, wrote to %v", time.Now().Sub(start_time), *output_file_path)
}
