// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

'use strict';

console.log('aaaaaaaaaaaaa');
chrome.extension.getBackgroundPage().console.log('foo');


try {
    var xhr = new XMLHttpRequest();
    xhr.open(
        'GET', 
        'https://rate.taobao.com/feedRateList.htm?auctionNumId=39852367764&currentPageNum=2&pageSize=20&rateType=&orderType=sort_weight&attribute=&sku=&hasSku=false&folded=0&callback=jsonp_tbcrate_reviews_list',
         false
        );
    xhr.send();

    var result = xhr.responseText;
    var results = JSON.parse(string.JSONresult.match(/jsonp_tbcrate_reviews_list\((.*)\)/)[1]);

// TODO: 不用登录也能拿到评论信息，那么还需要写在这里吗
// 1. 直接js，那么会发送很多个请求，数据量肯定是全量。如果上万条评论，后果还是有点儿不堪设想
// 2. 后端服务器，客户端与服务端交互只有一次，也永远不会返回全量数据，服务端可以自己并发处理，并且可以缓存。貌似没什么缺点。。
// 而且后端可以用Go语言了，异步并发。用户进入网页的时候就可以去并发了，用户点击的时候再直接获取就行了撒，666.唉，复杂是复杂了点的啦
// 唉，用go的话就只能用自己在国内的服务器了。但是自己维护服务器的话，唉，复杂度就有点大了
// 或者说我分开做？做两套，如果查不到或者说如果量太大就走服务端？但是这样就有了两套。。。除非后端用js这样增加代码的重用性。
// 但是js的性能和go的性能能比吗

    results.comments;
    results.currentPageNum;
    results.maxPage;
    results.total;
    results.watershed=100;

    chrome.extension.getBackgroundPage().console.log(result);
} catch (error) {
    chrome.extension.getBackgroundPage().console.log(error);
}



// let changeColor = document.getElementById('changeColor');

// chrome.storage.sync.get('color', function(data) {
//   changeColor.style.backgroundColor = data.color;
//   changeColor.setAttribute('value', data.color);
// });

// changeColor.onclick = function(element) {
//   let color = element.target.value;
//     chrome.tabs.executeScript(
//         tabs[0].id,
//         {code: 'document.body.style.backgroundColor = "' + color + '";'});
// };
