(function () {
    let main = new Vue({
        el: "#app",
        data: function () {
            return {
                items: [],
                visible: false,
                search_input: ""
            }
        },
        methods: {
            searchCommodity: function (event) {
                let th = this;
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
                });
            }
        }
    });
})();
