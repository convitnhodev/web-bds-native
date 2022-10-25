var Select2Selects = function () {
    var _componentSelect2 = function () {
        if (!$().select2) {
            console.warn('Warning - select2.min.js is not loaded.');
            return;
        }
        $('.select').select2({
            minimumResultsForSearch: Infinity
        });
        $('.provinces-select-remote').select2({
            placeholder: "Chọn tỉnh/thành phố",
            ajax: {
                url: '/api/v1/locations/provinces',
                dataType: 'json',
                delay: 250,
                processResults: function (data, params) {
                    return {
                        results: $.map(data, function (item) {
                            return {
                                text: item.name,
                                id: item.code,
                                data: item
                            };
                        })
                    };
                },
                cache: true
            },
        });
        $('.cities-select-remote').select2({placeholder: "Chọn Quận/huyện"});
        $('.wards-select-remote').select2({placeholder: "Chọn Xã/phường"});
        $('.select-remote-data').select2({
            ajax: {
                url: 'https://api.github.com/search/repositories',
                dataType: 'json',
                delay: 250,
                data: function (params) {
                    return {
                        q: params.term, // search term
                        page: params.page
                    };
                },
                processResults: function (data, params) {
                    params.page = params.page || 1;
                    return {
                        results: data.items,
                        pagination: {
                            more: (params.page * 30) < data.total_count
                        }
                    };
                },
                cache: true
            },
            escapeMarkup: function (markup) {
                return markup;
            }, // let our custom formatter work
            minimumInputLength: 1,
        });
    };

    return {
        init: function () {
            _componentSelect2();
        }
    }
}();

var InputGroups = function () {
    var _componentTouchspin = function () {
        $('.touchspin-step').TouchSpin({
            min: 0,
            max: 100,
        });
    }
    return {
        init: function () {
            _componentTouchspin();
        }
    }
}();

document.addEventListener('DOMContentLoaded', function () {
    Select2Selects.init();
    // Plupload.init();
    InputGroups.init();
});

$(document).ready(function () {
    $(".provinces-select-remote").change(function () {
        var province = $(".provinces-select-remote").val();
        $('.cities-select-remote').select2({
            placeholder: "Chọn Quận/huyện",
            ajax: {
                url: '/api/v1/locations/provinces/' + province + '/districts',
                dataType: 'json',
                delay: 250,
                processResults: function (data, params) {
                    return {
                        results: $.map(data, function (item) {
                            return {
                                text: item.full_name,
                                id: item.code,
                                data: item
                            };
                        })
                    };
                },
                cache: true
            },
        });
    });
    $(".cities-select-remote").change(function () {
        var districts = $(".cities-select-remote").val();
        $('.wards-select-remote').select2({
            placeholder: "Chọn Xã/phường",
            ajax: {
                url: '/api/v1/locations/districts/' + districts + '/wards',
                dataType: 'json',
                delay: 250,
                processResults: function (data, params) {
                    return {
                        results: $.map(data, function (item) {
                            return {
                                text: item.full_name,
                                id: item.code,
                                data: item
                            };
                        })
                    };
                },
                cache: true
            },
        });
    });

    $("input[name=general_purchase_terms_select]").on("change", function () {
        // alert($("input[name=general_purchase_terms_select]:checked").val());
        if ($("input[name=general_purchase_terms_select]:checked").val() == "url") {
            $("div.general_purchase_terms").html('<input type="text" name="general_purchase_terms" class="form-control" placeholder="đường dẫn file">');
        } else {
            $("div.general_purchase_terms").html('<input type="file" class="form-control h-auto">');
        }
    });

    var uploader = "";

    $('.html5-uploader').pluploadQueue({
        runtimes: 'html5',
        url: '/admin/products/images/upload',
        multipart: true,
        file_data_name: 'file',
        filters: {
            max_file_size: '10mb',
            mime_types: [
                {title: "Image files", extensions: "jpg,jpeg"}
            ]
        },
        init: {
            FilesAdded: function (up, files) {
                uploader = up
            },
        }
    });

    $("#submit").click(function () {
        uploader.start();
    })
});