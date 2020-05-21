let main;
(function () {
    main = new Vue({
        el: "#app",
        data: function () {
            return {
                items: [],
                search_input: "",
                last_search_input: "",
                page: 1,
                v_loading: false,
                finished: false,
                loading: false,
            }
        },
        methods: {
            searchCommodity: function (event) {
                if (this.search_input == "") {
                    this.loading = false;
                    return
                }
                let th = this;
                let isInit = true;
                th.loading = true;
                if (th.last_search_input == th.search_input) {
                    th.page++;
                    isInit = false;
                } else {
                    th.last_search_input = th.search_input;
                    th.page = 1;
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
            }
        }
    });

    Vue.use(vant.Lazyload);
})();
