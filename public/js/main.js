$(function ($) {
    $("#search_button").click(function () {
        let q = $("#search_input").val();
        $.getJSON("/search", {"q": q}, function (json) {
            $.each(json, function (i, v) {
                console.log(i, v);
                let container = $("<div>").addClass("container-lg commodity-item");
                let link = $("<a>").addClass("commodity-link").attr("target", "_blank").attr("href", v.item_url);
                let item_pic = $("<div>").addClass("commodity-pic").html('<img src="' + v.pict_url + '">');
                let item_cont = $("<div>").addClass("commodity-content")
                    .html('<p class="comd-title">' + v.title + '</p>'
                        + '<p class="comd-price"><span>￥' + v.reserve_price + '</span>￥' + v.zk_final_price + '</p>'
                        + '<p class="comd-coupon">优惠券：' + v.coupon_info + '</p>'
                    );
                link.append(item_pic).append(item_cont).appendTo(container);
                container.appendTo("#commodity_list");
            });
        });
    });
});
