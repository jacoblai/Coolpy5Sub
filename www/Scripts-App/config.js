var routes = [
        new routing.routes.NavigationRoute("pages/index", "mg_users.html", {
            cacheView: true, isDefault: true, vmFactory: function (callback) {
                callback(new viewModels.IndexVM());
            }, title: "用户管理"
        }),
        new routing.routes.NavigationRoute("pages/cpwd", "mg_changepwd.html", {
            cacheView: true, vmFactory: function (callback) {
                callback(new viewModels.cpwdVM());
            }, title: "更新个人信息"
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