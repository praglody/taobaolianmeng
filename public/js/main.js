let main;
(function () {
    main = new Vue({
        el: "#app",
        data: function () {
            return {
                items: [],
                visible: false,
                search_input: "",
                v_loading: false
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
                            }
                        }
                        th.items = res;
                    }
                }).catch(function (error) {
                    console.log(error);
                }).finally(function () {
                    th.v_loading = false;
                });
            }
        }
    });

    Vue.use(vant.Lazyload);
})();
