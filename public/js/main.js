(function () {
    let main = new Vue({
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
                    th.recommendItems = th.recommendItems.concat(data);
                }).catch(function (err) {

                }).finally(function () {
                    th.recommendLoading = false;
                });
            }
        }
    });

    Vue.use(vant.Lazyload);
})();
