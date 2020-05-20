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
                    for (let i in response.data) {
                        if (response.data[i].coupon_info == "") {
                            response.data[i].coupon_info = "无";
                        }
                    }
                    th.items = response.data;
                }).catch(function (error) {
                    console.log(error);
                });
            }
        }
    });
})();
