package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

// Tmall class
type Tmall struct {
	siteName      string
	baseURL       string
	requestClient Request
}

// TmallItem goods
type TmallItem struct {
	ItemID   int64
	SellerID int64

	TotalPage     int64
	Paginator     string
	RateCount     string
	RateDanceInfo string
	CreatedTime   int64
	UpdatedTime   int64
}

// TmallComment comments
type TmallComment struct {
	ItemID   int64
	SellerID int64

	DisplayUserNick string
	GmtCreateTime   int64
	GoldUser        bool
	ID              int64
	Pics            []string
	RateContent     string
	RateDate        string
	Reply           string
	TradeEndTime    int64
	Useful          bool
	UserVipLevel    int64
	CreatedTime     int64
	UpdatedTime     int64
}

// NewTmall create Tmall class
func NewTmall() Tmall {
	return Tmall{
		siteName: "tmall",
		baseURL:  "https://rate.tmall.com/list_detail_rate.htm?itemId=%d&sellerId=%d&order=3&currentPage=%d&callback=jsonp",
		requestClient: NewRequest(RequestOptions{
			useRandomProxy: true,
			redisClient:    redisClient,
			retry:          100,
		}),
	}
}

func (tmall *Tmall) getComments(itemID string, sellerID string) {
	_itemID, err := strconv.ParseInt(itemID, 10, 64)
	_sellerID, err := strconv.ParseInt(sellerID, 10, 64)
	url := fmt.Sprintf(tmall.baseURL, _itemID, _sellerID, 1)

	// 第一次请求获取评论页码数量，并新建商品
	text, err := tmall.requestClient.get(url)
	if err != nil {
		log.Error(fmt.Sprintf("get page error: %s", err))
		return
	}
	tmallItem, tmallComments := parseTmallResult(_itemID, _sellerID, text)
	saveTmallItem(tmallItem)
	saveTmallComments(tmallComments)

	chnl := make(chan string)
	// for page := int64(2); page <= tmallItem.TotalPage; page++ {
	totalPage := int64(10)
	for page := int64(2); page <= totalPage; page++ {
		go tmall.goGetTamllCommentRequest(chnl, _itemID, _sellerID, page)
	}

	i := int64(0)
	for msg := range chnl {
		fmt.Println("Received ", msg)
		i = i + 1
		if i == totalPage-1 {
			break
		}
	}
	close(chnl)
}

func (tmall *Tmall) goGetTamllCommentRequest(chnl chan string, itemID int64, sellerID int64, page int64) {
	log.Info(fmt.Sprintf("getting page: %d", page))
	text, err := tmall.requestClient.get(fmt.Sprintf(tmall.baseURL, itemID, sellerID, page))
	if err != nil {
		log.Error(fmt.Sprintf("get page %d error: %s", page, err))
		return
	}

	log.Info(fmt.Sprintf("page %d : ", page) + text)
	tmallItem, tmallComments := parseTmallResult(itemID, sellerID, text)
	saveTmallItem(tmallItem)
	saveTmallComments(tmallComments)
	log.Info(fmt.Sprintf("page %d finished", page))
	chnl <- fmt.Sprintf("page %d finished", page)
}

func parseTmallResult(itemID int64, sellerID int64, result string) (TmallItem, []TmallComment) {
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
		ItemID:        itemID,
		SellerID:      sellerID,
		TotalPage:     results[0].Int(),
		Paginator:     results[1].String(),
		RateCount:     results[2].String(),
		RateDanceInfo: results[3].String(),
		CreatedTime:   time.Now().Unix(),
		UpdatedTime:   time.Now().Unix(),
	}

	var tmallItemComments []TmallComment

	for _, comment := range results[4].Array() {
		var pics []string
		for _, pic := range comment.Get("pics").Array() {
			pics = append(pics, pic.String())
		}

		tmallItemComments = append(tmallItemComments, TmallComment{
			ItemID:   itemID,
			SellerID: sellerID,

			DisplayUserNick: comment.Get("displayUserNick").String(),
			GmtCreateTime:   comment.Get("gmtCreateTime").Int(),
			GoldUser:        comment.Get("goldUser").Bool(),
			ID:              comment.Get("id").Int(),
			Pics:            pics,
			RateContent:     comment.Get("rateContent").String(),
			RateDate:        comment.Get("rateDate").String(),
			Reply:           comment.Get("reply").String(),
			TradeEndTime:    comment.Get("tradeEndTime").Int(),
			Useful:          comment.Get("usefule").Bool(),
			UserVipLevel:    comment.Get("userVipLevel").Int(),
			CreatedTime:     time.Now().Unix(),
			UpdatedTime:     time.Now().Unix(),
		})
	}

	return tmallItem, tmallItemComments
}

func saveTmallItem(tmallItem TmallItem) {
	b, err := json.Marshal(tmallItem)
	if err != nil {
		log.Error(err)
	}

	mongoClient.insertOne("tmall_items", string(b))
}

func saveTmallComments(tmallItemComments []TmallComment) {
	var documents []string
	for _, tmallComment := range tmallItemComments {
		b, err := json.Marshal(tmallComment)
		if err != nil {
			log.Error(err)
		}

		documents = append(documents, string(b))
	}

	mongoClient.insertMany("tmall_comments", documents)
}
