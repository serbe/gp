package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func getListURL(baseURL string, body []byte) error {
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
				if listIsOld(fullURL) {
					updateList(fullURL)

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
	mutex.Lock()
	defer mutex.Unlock()
	ipBytes, err := get([]byte("ips"), []byte(s))
	if err != nil {
		return false
	}
	_, err = bytesToIP(ipBytes)
	if err != nil {
		return false
	}
	return true
}

func existLink(s string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	linkBytes, err := get([]byte("links"), []byte(s))
	if err != nil {
		return false
	}
	_, err = bytesToLink(linkBytes)
	if err != nil {
		return false
	}
	return true
}

func listIsOld() {

}

func saveLink(s string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var linkS link

	linkB, err := get([]byte("links"), []byte(s))
	if err == nil {
		linkS, _ = bytesToLink(linkB)
	}
	linkS.lastCheck = time.Now()
	linkB, err = linkToBytes(linkS)
	if err != nil {
		return err
	}
	return put([]byte("links"), []byte(s), linkB)
}

func saveIP(addr, port string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var ipS ip

	fullAddress := addr + ":" + port

	ipB, err := get([]byte("ips"), []byte(fullAddress))
	if err == nil {
		ipS, _ = bytesToIP(ipB)
		ipS.addr = addr
		ipS.port = port
	}
	ipS.createAt = time.Now()
	ipB, err = ipToBytes(ipS)
	if err != nil {
		return err
	}
	return put([]byte("ips"), []byte(fullAddress), ipB)
}

func ipToBytes(s ip) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func linkToBytes(s link) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bytesToIP(b []byte) (ip, error) {
	var value ip
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}

func bytesToLink(b []byte) (link, error) {
	var value link
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		return value, err
	}
	return value, nil
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
