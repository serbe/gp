package main

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

func cleanBody(body []byte) []byte {
	for i := range replace {
		re := regexp.MustCompile(replace[i][0])
		if re.Match(body) {
			body = re.ReplaceAll(body, []byte(replace[i][1]))
		}
	}
	return body
}

func getListURL(baseURL string, body []byte) []string {
	var urls []string
	for i := range reURL {
		host, err := getHost(baseURL)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if re.Match(body) {
			allResults := re.FindAllSubmatch(body, -1)
			for _, result := range allResults {
				fullURL := host + "/" + string(result[1])
				if isOld(links.get(fullURL)) {
					links.set(fullURL)

					urls = append(urls, fullURL)
				}
			}
		}
	}
	return urls
}

func getListIP(body []byte) {
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				ip := string(res[1])
				port := string(res[2])
				portInt, _ := strconv.Atoi(port)
				if ip != "0.0.0.0" && portInt < 65535 {
					ipWithPort := ip + ":" + port
					if ips.get(ipWithPort).Addr == "" {
						numIPs++
						ips.set(ipWithPort, newIP(ip, port))
					}
				}
			}
		}
	}
	return
}

func newIP(addr, port string) ipType {
	var ip ipType
	ip.Addr = addr
	ip.Port = port
	ip.CreateAt = time.Now()
	return ip
}

func (ip ipType) encode() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(ip)
	return b.Bytes(), err
}

func (ip *ipType) decode(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(&ip)
}

func (link linkType) encode() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(link)
	return b.Bytes(), err
}

func (link *linkType) decode(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(&link)
}

func isOld(link linkType) bool {
	currentTime := time.Now()
	return currentTime.Sub(link.CheckAt) > time.Duration(15*time.Minute)
}

func compressZlib(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	_, err := w.Write(in)
	if err != nil {
		fmt.Println("error w.Write(in)")
		return nil, err
	}
	err = w.Close()
	return b.Bytes(), err
}

func decompressZlib(b []byte) ([]byte, error) {
	var in bytes.Buffer
	_, err := in.Read(b)
	if err != nil {
		fmt.Println("error in.Read(b))")
		return nil, err
	}
	var out bytes.Buffer
	r, err := zlib.NewReader(&in)
	if err != nil {
		fmt.Println("error zlib.NewReader(&in)")
		return nil, err
	}
	defer r.Close()
	_, err = io.Copy(&out, r)
	if err != nil {
		fmt.Println("error io.Copy(&out, r)")
		return nil, err
	}
	return out.Bytes(), err
}

func readDB() error {
	f, err := os.Open("db.zlib")
	if err != nil && err != io.EOF {
		fmt.Println("error os.Open('db.zlib')")
		return err
	}
	defer f.Close()
	fb := make([]byte, 5)
	_, err = f.Read(fb)
	if err != nil {
		fmt.Println("error f.Read(fb)")
		return err
	}
	dec, err := decompressZlib(fb)
	if err != nil {
		fmt.Println("error decompressZlib(fb)")
		return err
	}
	err = ioutil.WriteFile("ips.db", dec, 0644)
	if err != nil {
		fmt.Println("error ioutil.WriteFile('ips.db')")
		return err
	}
	f.Close()
	return os.Remove("db.zlib")
}

func saveDB() error {
	f, err := os.Open("ips.db")
	if err != nil {
		fmt.Println("error os.Open('ips.db')")
		return err
	}
	defer f.Close()
	fb := make([]byte, 5)
	_, err = f.Read(fb)
	if err != nil {
		fmt.Println("error f.Read(fb)")
		return err
	}
	comp, err := compressZlib(fb)
	if err != nil {
		fmt.Println("error compressZlib(fb)")
		return err
	}
	err = ioutil.WriteFile("db.zlib", comp, 0644)
	if err != nil {
		fmt.Println("error ioutil.WriteFile('db.zlib')")
		return err
	}
	f.Close()
	return os.Remove("ips.db")
}
