package services

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

func WriteFile(base64_image_content string) (path string, err error) {

	b, err := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)
	if !b {
		return "", err
	}

	//data:image/jpeg;base64,/9j/4R/+RXhpZgAATU0AKgAAA

	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取

	base64Str := re.ReplaceAllString(base64_image_content, "")

	curFileStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(99999)

	dir, err := os.Getwd()
	fmt.Println(dir)

	//  /var/www/html/aogo/seedHabits/images/158618123006229642477069.jpeg
	var file string = dir + "/images/" + curFileStr + strconv.Itoa(n) + "." + fileType
	fmt.Println("file", file)

	dataImgPath := curFileStr + strconv.Itoa(n) + "." + fileType
	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err = ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return dataImgPath, nil
}
