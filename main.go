package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/ioutil"
	"bufio"
	"regexp"
	"strconv"
	"net/http"
)

func get_files(dir_root string) []string {
	files := []string{}
	filepath.Walk(dir_root, func (path string, info os.FileInfo, _ error) error {
		file_tree := strings.Split(path, "/")
		if !info.IsDir() && !(file_tree[len(file_tree) - 1][0] == '.') {
			files = append(files, path)
		}
		return nil
	})

	return files
}

func write_resources_file(file_body string) {
	file_content := fmt.Sprintf(resources_format, file_body)
	fd, _ := os.Create("resources.go")
	defer fd.Close()

	w := bufio.NewWriter(fd)
	w.Write([]byte(file_content))
	w.Flush()

}

func parse_URLs(url_map *map[string]string, file_content string) {
	matches := html_pattern.FindAllStringSubmatch(file_content, -1)
	for _, match := range matches {
		url := match[2]
		if _, ok := (*url_map)[url]; ok {
			continue
		}

		if url_pattern.MatchString(url) {
			resp, _ := http.Get(url)
			body, _ := ioutil.ReadAll(resp.Body)

			(*url_map)[url] = string(body)
			parse_URLs(url_map, string(body))

			resp.Body.Close()
		}
	}
}

func subsituteURLs(file_map *map[string]string, url_map *map[string]string) map[string]string {
	num := 0
	name_swap := map[string]string{}

	for old_url := range *url_map {
		split := strings.Split(old_url, ".")
		ending := "." + split[len(split) - 1]

		file_name := "web" + strconv.Itoa(num) + ending
		new_url := "http://localhost:8080/" + file_name

		name_swap[old_url] = file_name
		for k := range *url_map {
			(*url_map)[k] = strings.ReplaceAll((*url_map)[k], old_url, new_url)
		}
		for k := range *file_map {
			(*file_map)[k] = strings.ReplaceAll((*file_map)[k], old_url, new_url)
		}
	}

	return name_swap
}

func create_file_body(name_swap map[string]string,
	file_map *map[string]string,
	url_map *map[string]string) string {

	ret := ""

	for name, val := range *file_map {
		relitive_filename := strings.Replace(name, os.Args[1], "", 1)

		if strings.HasSuffix(name, ".html") {
			val = strings.ReplaceAll(val, "`", "`+\"`\"+`")
		}

		bytes_decl := bytes_to_bytes_decl([]byte(val))
		ret += fmt.Sprintf(file_format, relitive_filename, bytes_decl)
	}

	for name, val := range *url_map {
		swaped_name := name_swap[name]

		if strings.HasSuffix(name, ".html") {
			val = strings.ReplaceAll(val, "`", "`+\"`\"+`")
		}

		bytes_decl := bytes_to_bytes_decl([]byte(val))
		ret += fmt.Sprintf(file_format, swaped_name, bytes_decl)
	}

	return ret
}

func bytes_to_bytes_decl(bytes []byte) string {
	printed_bytes := fmt.Sprint(bytes)
	printed_bytes = printed_bytes[1 : len(printed_bytes) - 1]
	printed_bytes = strings.ReplaceAll(printed_bytes, " ", ",")

	ret := fmt.Sprintf("[]byte{%s}", printed_bytes)
	return ret
}

func read_bytes(file *os.File) []byte {
	buffer := []byte{}
	buffer_reader := bufio.NewReader(file)

	for {
		read, err := buffer_reader.ReadByte()
		if err != nil {
			break
		} else {
			buffer = append(buffer, read)
		}
	}

	return buffer
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./quark <web_folder>")
		os.Exit(1)
	}

	html_pattern = regexp.MustCompile(html_regex)
	url_pattern = regexp.MustCompile(url_regex)

	files := get_files(os.Args[1])

	file_map := map[string]string{}
	url_map := map[string]string{}

	for _, filename := range files {
		file, _ := os.Open(filename)
		content_bytes := read_bytes(file)
		content := string(content_bytes)

		file_map[filename] = content
		parse_URLs(&url_map, content)
	}

	spawed_names := subsituteURLs(&file_map, &url_map)
	file_body := create_file_body(spawed_names, &file_map, &url_map)
	write_resources_file(file_body)
}
