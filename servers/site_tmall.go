
package main

import (
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Tmall class
type Tmall struct {
	baseURL string
	requestClient Request
}

// TmallItem goods
type TmallItem struct {
	itemID int64
	sellerID int64
	
	totalPage int64
	paginator string
	rateCount string
	rateDanceInfo string
}

// TmallItemComment comments
type TmallItemComment struct {
	itemID int64
	sellerID int64

	displayUserNick string
	gmtCreateTime int64
	goldUser bool
	id int64
	pics []string
	rateContent string
	rateDate string
	reply string
	tradeEndTime int64
	useful bool
	userVipLevel int64
}

// NewTmall create Tmall class
func NewTmall() Tmall {
	return Tmall{
		baseURL: "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&order=3&currentPage=%d&callback=jsonp",
		requestClient: NewRequest(RequestOptions{
			useRandomProxy: true, 
			redisClient: redisClient,
		}),
	}
}

func (tmall *Tmall) getComments(itemID string, sellerID string) {
	// chnl := make(chan string)
	url := fmt.Sprintf(tmall.baseURL, itemID, sellerID, 1)

	// 第一次请求获取评论页码数量，并新建商品
	text, err := tmall.requestClient.get(url)
	if err != nil {
		log.Error(fmt.Sprintf("get page error: %s", err))
		return
	}


	reg := regexp.MustCompile(`"page":(\d)}`)
	matches := reg.FindStringSubmatch(text)
	if len(matches) < 1 {
		log.Error(fmt.Sprintf("page 1 error: %s", text))
		return
	}
	totalPage, err := strconv.Atoi(matches[1])
	if err != nil {
	}
	fmt.Println(totalPage)

	// for page := 1; page <= 5; page++ {
	// 	go getTmallCommentRequest(chnl, fmt.Sprintf(baseURL, itemID, sellerID, page), page, totalPage)
	// }

	// for msg := range chnl {
	// 	fmt.Println("Received ", msg)
	// }
}

func (tmall *Tmall) sendRequest(chnl chan string, url string, page int, totalPage int) {
	text, err := tmall.requestClient.get(url)
	if err != nil {
		log.Error(fmt.Sprintf("error in page: %d, %s", page, err))
		if page == totalPage {
			close(chnl)
		}
	}

	if strings.Contains(text, "rateList") {
		log.Info(fmt.Sprintf("ok in page: %d", page))
	}
	chnl <- strconv.Itoa(page)
	if page == totalPage {
		close(chnl)
	}
}

func (tmall *Tmall) parseResult(itemID int64, sellerID int64, result string) (TmallItem, []TmallItemComment) {
	reg := regexp.MustCompile(`jsonp.*?\((.*)`)
	match := reg.FindStringSubmatch(result)
	text := match[1]

	results := gjson.GetMany(
		text, 
		"rateDetail.paginator.lastPage", 
		"rateDetail.paginator", 
		"rateDetail.rateCount", 
		"rateDetail.rateDanceInfo", 
		"rateDetail.rateList",
	)
	
	tmallItem := TmallItem{
		itemID: itemID,
		sellerID: sellerID,
		totalPage: results[0].Int(),
		paginator: results[1].String(),
		rateCount: results[2].String(),
		rateDanceInfo: results[3].String(),
	}

	var tmallItemComments []TmallItemComment

	for _, comment := range results[4].Array() {
		var pics []string
		for _, pic := range comment.Get("pics").Array() {
			pics = append(pics, pic.String())
		}

		tmallItemComments = append(tmallItemComments, TmallItemComment{
			itemID: itemID,
			sellerID: sellerID,
			
			displayUserNick: comment.Get("displayUserNick").String(),
			gmtCreateTime: comment.Get("gmtCreateTime").Int(),
			goldUser: comment.Get("goldUser").Bool(),
			id: comment.Get("id").Int(),
			pics: pics,
			rateContent: comment.Get("rateContent").String(),
			rateDate: comment.Get("rateDate").String(),
			reply: comment.Get("reply").String(),
			tradeEndTime: comment.Get("tradeEndTime").Int(),
			useful: comment.Get("usefule").Bool(),
			userVipLevel: comment.Get("userVipLevel").Int(),
		})
	}

	return tmallItem, tmallItemComments
}
