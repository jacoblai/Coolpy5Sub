if ($.session.get('pass') == undefined) {
    window.location.href = "login.html";
}

window.onbeforeunload = function () {
    $.session.clear();
}

$.ajaxSetup({
    cache: false,
    beforeSend: function (xhr) {
        xhr.setRequestHeader("Authorization", $.session.get('pass'));
    }
});

var __extends = this.__extends || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    __.prototype = b.prototype;
    d.prototype = new __();
};
function gotoview(path, params, payload) {
    var navigateOptions = {
        payload: payload,
        forceReloadOnNavigation: true,
        forceNavigationInCache: false,
    };
    router.navigateTo(path + params, navigateOptions);
}
var viewModels;
(function (viewModels) {
    var BaseVM = (function () {
        function BaseVM() {
            for (var i in this) {
                if (!this.hasOwnProperty(i) && typeof (this[i]) === "function" && i != "constructor") {
                    this[i] = this[i].bind(this);
                }
            }
            ko.validation.locale('zh-CN');
        }
        var self = this;
        self.main_uid = ko.observable()
        self.LoadUserInfo = function () {
            $.ajax({
                method: "GET",
                url: basicurl + '/api/user/' + $.session.get('uid'),
                success: function (result) {
                    if (result.ok == 1) {
                        self.main_uid(result.data.Uid)
                    }
                },
                error: function (xhr, status, error) {
                }
            })
        }
        self.LoadUserInfo()

        self.users = function () {
            gotoview("pages/index", "", new Array());
        };

        return BaseVM;
    })();
    viewModels.BaseVM = BaseVM;

    var IndexVM = (function (_super) {
        __extends(IndexVM, _super);
        var self = this;
        function IndexVM() { _super.call(self); }

        IndexVM.prototype.onNavigatedTo = function () {

        }
        return IndexVM;
    })(BaseVM);
    viewModels.IndexVM = IndexVM;

    var cpwdVM = (function (_super) {
        __extends(cpwdVM, _super);
        var self = this;
        function cpwdVM() { _super.call(self); }
        cpwdVM.prototype.onNavigatedTo = function () {

        }
        return cpwdVM;
    })(BaseVM);
    viewModels.cpwdVM = cpwdVM;

})(viewModels || (viewModels = {}));