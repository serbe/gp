package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func getListURL(baseURL string, body []byte) error {
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
				if linkIsOld(fullURL) {
					saveLink(fullURL)

					resultChan <- fullURL
				}
			}
		}
	}
	return nil
}

func getListIP(body []byte) {
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			fmt.Printf("find %d ip\n", len(results))
			for _, res := range results {
				ip := string(res[1])
				port := string(res[2])
				portInt, _ := strconv.Atoi(port)
				if ip != "0.0.0.0" && portInt < 65535 {
					ipWithPort := ip + ":" + port
					if !existIP(ipWithPort) {
						numIPs++
						saveIP(ip, port)
					}
				}
			}
		}
	}
	return
}

func existIP(s string) bool {
	ip, err := getIP(s)
	if err != nil || ip.Addr == "" {
		return false
	}
	return true
}

func existLink(s string) bool {
	link, err := getLink(s)
	if err != nil || link.Host == "" {
		return false
	}
	return true
}

func linkIsOld(fullURL string) bool {
	link, err := getLink(fullURL)
	if err != nil {
		return true
	}
	return isOld(link.CheckAt)
}

func getLink(fullURL string) (linkType, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var link linkType
	byteArray, err := get([]byte("links"), []byte(fullURL))
	if err != nil {
		return link, err
	}
	link.decode(byteArray)
	return link, err
}

func saveLink(fullURL string) error {
	var link linkType
	link, err := getLink(fullURL)
	if err != nil {
		return err
	}
	link.CheckAt = time.Now()
	link.Host, link.Ssl = getLinkAttribs(fullURL)
	byteArray, err := link.encode()
	if err != nil {
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()
	return put([]byte("links"), []byte(fullURL), byteArray)
}

func getLinkAttribs(s string) (string, bool) {
	if strings.Contains(s, "http://") {
		return s[7:], false
	}
	if strings.Contains(s, "https://") {
		return s[8:], true
	}
	return s, false
}

func getIP(fullAddress string) (ipType, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var ip ipType
	byteArray, err := get([]byte("ips"), []byte(fullAddress))
	if err != nil {
		return ip, err
	}
	ip.decode(byteArray)
	return ip, err
}

func saveIP(addr, port string) error {
	var ip ipType
	fullAddress := addr + ":" + port
	ip, err := getIP(fullAddress)
	if err != nil {
		return err
	}
	ip.Addr = addr
	ip.Port = port
	ip.CreateAt = time.Now()
	byteArray, err := ip.encode()
	if err != nil {
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()
	return put([]byte("ips"), []byte(fullAddress), byteArray)
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

func toBytes(data interface{}) ([]byte, error) {
	var (
		v   []byte
		err error
	)

	switch val := data.(type) {
	case string:
		v = []byte(val)
	case []byte:
		v = val
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v = []byte(fmt.Sprintf("%d", val))
	case float64, float32:
		v = []byte(fmt.Sprintf("%f", val))
	case fmt.Stringer:
		v = []byte(val.String())
	default:
		err = fmt.Errorf("non supported types")
	}
	return v, err
}

func isOld(t time.Time) bool {
	tn := time.Now()
	return tn.Sub(t) > time.Duration(15)*time.Second
}
