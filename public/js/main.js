(function () {

    function parsePrice(num) {
        return (Math.round(num * 100) / 100).toFixed(2);
    }

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
                materialId: "",
                recommendItems: [],
                recommendLoading: false,
                recommendFinished: false,
                shareKey: "",
                shareUrl: "",
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
                            th.finished = true;
                            th.loading = false;
                            return;
                        }

                        for (let i in res) {
                            if (res[i].coupon_info == "") {
                                res[i].coupon_info = "无";
                            } else {
                                res[i].use_coupon = parsePrice(res[i].zk_final_price - res[i].coupon_amount);
                            }
                            res[i].zk_final_price = parsePrice(res[i].zk_final_price);
                        }
                        if (isInit) {
                            th.items = res;
                        } else {
                            th.items = th.items.concat(res);
                        }

                        setTimeout(function () {
                            th.loading = false;
                        }, 1500);
                    } else {
                        th.finished = true;
                        th.loading = false;
                    }
                }).catch(function (error) {
                    console.log(error);
                    th.finished = true;
                    th.loading = false;
                })
            },
            getRecommendList: function () {
                let th = this;
                axios.get('/recommend', {
                    params: {
                        page: th.recommendPage,
                        material_id: th.materialId,
                        page_size: 20
                    }
                }).then(function (resp) {
                    if (resp.data.code == 200) {
                        let data = resp.data.data.result;
                        if (data.length < 1) {
                            th.recommendLoading = false;
                            th.recommendFinished = true;
                            return;
                        }

                        th.materialId = resp.data.data.material_id;
                        th.recommendPage++;
                        for (let i in data) {
                            data[i].zk_final_price = parsePrice(data[i].zk_final_price);
                            if (data[i].coupon_amount == undefined || data[i].coupon_amount == 0) {
                                data[i].coupon_info = "无";
                            } else {
                                data[i].coupon_info = "满 " + parsePrice(data[i].coupon_start_fee) + " 元减 "
                                    + parsePrice(data[i].coupon_amount) + " 元";
                                data[i].use_coupon = parsePrice(data[i].zk_final_price - data[i].coupon_amount);
                            }
                            data[i].item_url = data[i].click_url;
                            if (data[i].shop_title == undefined || data[i].shop_title == "") {
                                data[i].shop_title = data[i].nick;
                            }
                            data[i].item_url = data[i].click_url;
                        }
                        th.recommendItems = th.recommendItems.concat(data);

                        setTimeout(function () {
                            th.recommendLoading = false;
                            th.recommendFinished = false;
                        }, 1500);
                    } else {
                        th.recommendLoading = false;
                        th.recommendFinished = true;
                    }
                }).catch(function (err) {
                    console.log(err);
                    th.recommendLoading = false;
                    th.recommendFinished = true;
                });
            },
            copyShareKey: function (e) {
                let th = this;
                let itemId = e.target.getAttribute("item-id");
                let commodity = {};

                for (let i in app.recommendItems) {
                    if (app.recommendItems[i].item_id == itemId) {
                        commodity = app.recommendItems[i];
                    }
                }

                for (let i in app.items) {
                    if (app.items[i].item_id == itemId) {
                        commodity = app.items[i];
                    }
                }

                if (commodity.item_id == undefined || commodity.item_id == null) {
                    return;
                }

                let title = commodity.title;
                let url = "https:";

                if (commodity.coupon_share_url) {
                    url += commodity.coupon_share_url;
                } else {
                    if (commodity.url) {
                        url += commodity.url;
                    } else if (commodity.item_url) {
                        url += commodity.item_url;
                    } else {
                        vant.Toast.success('获取商品链接失败');
                        return;
                    }
                }

                th.shareUrl = url;
                axios.post('/get-share-key', {
                    title: title,
                    url: url
                }).then(function (resp) {
                    if (resp.data.code == 200) {
                        let key = resp.data.data.result.model;
                        if (commodity.use_coupon) {
                            th.shareKey = title + "\n【在售价】" + parsePrice(commodity.zk_final_price) + "元\n";
                            th.shareKey += "【券后价】" + commodity.use_coupon + "元\n";
                        } else {
                            th.shareKey = title + "\n【折扣价】" + parsePrice(commodity.zk_final_price) + "元\n";
                        }
                        th.shareKey += "-----------------\n" +
                            "注意，请完整复制这条信息，" + key + "，到【手机淘宝】即可查看";
                        document.getElementById("copy").click();
                    } else {
                        //vant.Toast.success('获取淘口令失败');
                        window.open(url,"_blank");
                    }
                }).catch(function (err) {
                    console.log(err)
                    //vant.Toast.success('服务器错误');
                    window.open(th.shareUrl,"_blank");
                });
            },
            viewCommodity: function(e) {
                let th = this;
                let itemId = e.target.getAttribute("item-id");
                let commodity = {};

                for (let i in app.recommendItems) {
                    if (app.recommendItems[i].item_id == itemId) {
                        commodity = app.recommendItems[i];
                    }
                }

                for (let i in app.items) {
                    if (app.items[i].item_id == itemId) {
                        commodity = app.items[i];
                    }
                }

                if (commodity.item_id == undefined || commodity.item_id == null) {
                    return;
                }

                let images = [];
                if (commodity.small_images.string.length > 0) {
                    images = commodity.small_images.string;
                }
                images = images.concat(commodity.pict_url);

                vant.ImagePreview({
                    images: images,
                    closeable: true,
                    closeOnPopstate: true
                });
            }
        }
    });

    Vue.use(vant.Lazyload);

    let clipboard = new ClipboardJS('#copy', {
        text: function () {
            return app.shareKey;
        }
    });

    clipboard.on('success', function (e) {
        vant.Toast.success('优惠券已复制到剪贴板，打开淘宝APP即可领取');
    });

    clipboard.on('error', function (e) {
        //vant.Toast.success('您的手机不支持自动领取，请手动复制淘口令');
        axios.post('/report-error', {
            code: 10060,
            msg: "用户设备不支持复制"
        }).then(function (resp) {
        }).catch(function (err) {
            console.log(err)
        });
        window.open(th.shareUrl,"_blank");
    });
})();
