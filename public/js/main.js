let main;
(function () {
    main = new Vue({
        el: "#app",
        data: function () {
            return {
                items: [],
                search_input: "",
                v_loading: false,
                finished: false,
                loading: false
            }
        },
        methods: {
            searchCommodity: function (event) {
                if (this.search_input == "") {
                    return
                }
                let th = this;
                th.v_loading = true;
                axios.get('/search', {
                    params: {
                        q: th.search_input
                    }
                }).then(function (response) {
                    if (response.data.code == 200) {
                        let res = response.data.data.result;
                        for (let i in res) {
                            if (res[i].coupon_info == "") {
                                res[i].coupon_info = "无";
                            }else{
                                res[i].use_coupon = res[i].zk_final_price - res[i].coupon_amount
                                res[i].use_coupon = res[i].use_coupon.toFixed(2);
                            }
                            console.log(res[i])
                        }
                        th.items = res;
                    }
                }).catch(function (error) {
                    console.log(error);
                }).finally(function () {
                    th.v_loading = false;
                });
            },
            getMoreCommodity: function () {
                let th = this;
                console.log("fadsfas")
                th.loading = false;
            }
        }
    });

    Vue.use(vant.Lazyload);
})();
