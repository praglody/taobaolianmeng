(function () {
    let app = new Vue({
        el: "#app",
        data: function () {
            return {
                searchList: false,
                loading: false,
                finished: false,
                items: [],
                search_input: "",
                last_search_input: "",
                page: 1,
                recommend: true,
                recommendPage: 1,
                recommendItems: [],
                recommendLoading: false,
                recommendFinished: false,
                shareKey: "",
            }
        },
        methods: {
            searchCommodity: function () {
                this.last_search_input = "";
                this.listCommodity();
            },
            listCommodity: function () {
                let th = this;
                if (this.search_input == "") {
                    th.loading = false;
                    return
                }

                if (th.recommend == true) {
                    th.recommend = false;
                    th.searchList = true;
                }

                let isInit = true;
                th.loading = true;
                if (th.last_search_input == th.search_input) {
                    th.page++;
                    isInit = false;
                } else {
                    th.last_search_input = th.search_input;
                    th.page = 1;
                    th.items = [];
                }
                axios.get('/search', {
                    params: {
                        q: th.search_input,
                        p: th.page
                    }
                }).then(function (response) {
                    if (response.data.code == 200) {
                        let res = response.data.data.result;
                        if (res.length < 1) {
                            th.page--;
                            if (th.page < 1) {
                                th.page = 1;
                            }
                            return;
                        }

                        for (let i in res) {
                            if (res[i].coupon_info == "") {
                                res[i].coupon_info = "无";
                            } else {
                                res[i].use_coupon = res[i].zk_final_price - res[i].coupon_amount
                                res[i].use_coupon = res[i].use_coupon.toFixed(2);
                            }
                        }
                        if (isInit) {
                            th.items = res;
                        } else {
                            th.items = th.items.concat(res);
                        }
                    }
                }).catch(function (error) {
                    console.log(error);
                }).finally(function () {
                    setTimeout(function () {
                        th.loading = false;
                    }, 1500);
                });
            },
            getRecommendList: function () {
                let th = this;
                axios.get('/recommend', {
                    params: {
                        page: th.recommendPage,
                        page_size: 20
                    }
                }).then(function (resp) {
                    let data = resp.data.data.result;
                    if (data.length > 0) {
                        th.recommendPage++;
                    }
                    for (let i in data) {
                        if (data[i].coupon_amount == undefined || data[i].coupon_amount == 0) {
                            data[i].coupon_info = "无";
                        } else {
                            data[i].coupon_info = "满 " + data[i].coupon_start_fee + " 元减 "
                                + data[i].coupon_amount + " 元";
                            data[i].use_coupon = data[i].zk_final_price - data[i].coupon_amount
                            data[i].use_coupon = data[i].use_coupon.toFixed(2);
                        }
                    }
                    th.recommendItems = th.recommendItems.concat(data);

                }).catch(function (err) {
                    console.log(err)
                }).finally(function () {
                    setTimeout(function () {
                        th.recommendLoading = false;
                    }, 1500);
                });
            },
            copyShareKey: function (e) {
                let th = this;
                let itemId = e.target.getAttribute("item-id");

                for (let i in app.items) {
                    if (app.items[i].item_id == itemId) {
                        commodity = app.items[i];
                    }
                }

                for (let i in app.recommendItems) {
                    if (app.recommendItems[i].item_id == itemId) {
                        commodity = app.recommendItems[i];
                    }
                }

                if (commodity.item_id == undefined || commodity.item_id == null) {
                    return;
                }

                let title = commodity.title;
                let url = "https:" + commodity.coupon_share_url;

                axios.post('/get-share-key', {
                    title: title,
                    url: url
                }).then(function (resp) {
                    if (resp.data.code == 200) {
                        let key = resp.data.data.result.model;
                        th.shareKey = key;
                        vant.Toast.success('优惠券已复制到剪贴板，打开淘宝APP即可领取');
                        document.getElementById("copy").click();
                    } else {
                        vant.Toast.success('领取失败');
                    }
                }).catch(function (err) {
                    console.log(err)
                    vant.Toast.success('领取失败');
                });
            }
        }
    });

    Vue.use(vant.Lazyload);

    new ClipboardJS('#copy', {
        text: function () {
            return app.shareKey;
        }
    });
})();
