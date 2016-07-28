$.ajaxSetup({
    cache: false,
    beforeSend: function (xhr) {
        xhr.setRequestHeader("Authorization", "Basic " + localStorage.getItem('token'));
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

})(viewModels || (viewModels = {}));