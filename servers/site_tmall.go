package main

import (
	"strconv"
	"strings"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// get tmall comment
func getTmallComment(itemID string, sellerID string) {
	// chnl := make(chan string)
	baseURL := "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&order=3&currentPage=%d&callback=jsonp"
	url := fmt.Sprintf(baseURL, itemID, sellerID, 1)

	// 第一次请求获取评论页码数量
	var requestOptions RequestOptions
	requestOptions.proxy = "http://haofly.net"
	// requestOptions.proxy = "http://183.232.223.142:80"

	response, err := requestGet(url, requestOptions)
	if err != nil {
		log.Error(fmt.Sprintf("get page error: %s", err))
		return
	}
	fmt.Println(response)

	// reg := regexp.MustCompile(`"page":(\d)}`)
	// matches := reg.FindStringSubmatch(response)
	// if len(matches) < 1 {
	// 	log.Error(fmt.Sprintf("page 1 error: %s", response))
	// 	return
	// }
	// totalPage, err := strconv.Atoi(matches[1])
	// if err != nil {
	// }

	// for page := 1; page <= 5; page++ {
	// 	go getTmallCommentRequest(chnl, fmt.Sprintf(baseURL, itemID, sellerID, page), page, totalPage)
	// }

	// for msg := range chnl {
	// 	fmt.Println("Received ", msg)
	// }
}


func getTmallCommentRequest(chnl chan string, url string, page int, totalPage int) {
	var requestOptions RequestOptions
	// requestOptions.proxy = "http://haofly.net"

	response, err := requestGet(url, requestOptions)
	if err != nil {
		log.Error(fmt.Sprintf("error in page: %d, %s", page, err))
		if page == totalPage {
			close(chnl)
		}
	}

	if strings.Contains(response, "rateList") {
		log.Info(fmt.Sprintf("ok in page: %d", page))
	}
	chnl <- strconv.Itoa(page)
	if page == totalPage {
		close(chnl)
	}
}