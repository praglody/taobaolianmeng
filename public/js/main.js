let commodity;
(function () {
    let main = new Vue({
        el: "#main",
        data: function() {
            return { visible: false }
        }
    });

    commodity = new Vue({
        el: "#commodity_list",
        data: {
            items: []
        }
    });

    let search = new Vue({
        el: "#search_group",
        data: {
            search_input: ""
        },
        methods: {
            searchCommodity: function (event) {
                let th = this;
                axios.get('/search', {
                    params: {
                        q: th.search_input
                    }
                }).then(function (response) {
                    for (let i in response.data) {
                        if (response.data[i].coupon_info == "") {
                            response.data[i].coupon_info = "无";
                        }
                    }
                    commodity.items = response.data;
                }).catch(function (error) {
                    console.log(error);
                });
            }
        }
    });
})();
