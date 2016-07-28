var routes = [
        new routing.routes.NavigationRoute("pages/total", "mainpage.html", {
            cacheView: true, vmFactory: function (callback) {
            callback(new viewModels.IndexVM());
            }, title: "数据模型"
        }),
        new routing.routes.NavigationRoute("pages/index", "firstpage.html", {
            cacheView: true, isDefault: true, vmFactory: function (callback) {
                callback(new viewModels.firstpageVM());
            }, title: "一门式自然人库"
        }),
        new routing.routes.NavigationRoute("pages/pl/{uid}/{?value}", "personlabel.html?uid={uid}", {
            cacheView: false, vmFactory: function (callback) {
                callback(new viewModels.plVM());
            }, title: "自然人标签"
        }),
        new routing.routes.NavigationRoute("pages/pf/{uid}/{?value}", "profile.html?uid={uid}", {
            cacheView: false, vmFactory: function (callback) {
                callback(new viewModels.pfVM());
            }, title: "自然人祥细资料"
        }),
        new routing.routes.NavigationRoute("pages/more", "morepage.html", {
            cacheView: true, vmFactory: function (callback) {
                callback(new viewModels.moreVM());
            }, title: "自人库展示"
        }),
        new routing.routes.NavigationRoute("pages/detailes/{skin}/{?index}/{?ps}", "viewdetailes.html?skin={skin}", {
            cacheView: false, vmFactory: function (callback) {
                callback(new viewModels.detailesVM());
            }, title: "数据浏览"
        }),
        new routing.routes.NavigationRoute("pages/logs/{?skin}/{?index}/{?ps}", "viewlogs.html", {
            cacheView: false, vmFactory: function (callback) {
                callback(new viewModels.logsVM());
            }, title: "日志浏览"
        }),
        new routing.routes.NavigationRoute("pages/dm", "datamodelpage.html", {
            cacheView: false, vmFactory: function (callback) {
                callback(new viewModels.dmVM());
            }, title: "数据模型应用"
        })
];
var router = new routing.Router("views-placeholder", // ID of element in which will be loaded views.
{
    beforeNavigation: function () { },    // Global before navigation handler.
    afterNavigation: function () {
        //if (history.length === 0) {
        //    $("#bt_goback").attr("visibility", "hidden");
        //} else {
        //    $("#bt_goback").attr("visibility", "visible");
        //}
    },     // Global after navigation handler.
    navigationError: function () {
        router.navigateBack();
    }, // Global navigation error handler.
    enableLogging: false
},
routes);        // Routes are described below.
// This is the array of Route objects.
routing.knockout.setCurrentRouter(router);
ko.applyBindings({}); // This requered to allow ko bindings to work ewrywhere on the page.
// You can put here root level view model of application.
router.run(); // Starting of router.