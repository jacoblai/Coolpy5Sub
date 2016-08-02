if ($.session.get('pass') == undefined) {
    window.location.href = "login.html";
}

window.onbeforeunload = function () {
    
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
                        self.main_uid(result.data.UserName)
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

        self.logout = function () {
            $.session.clear();
            window.location.href = "login.html";
        }

        return BaseVM;
    })();
    viewModels.BaseVM = BaseVM;

    var IndexVM = (function (_super) {
        __extends(IndexVM, _super);
        var self = this;
        function IndexVM() { _super.call(self); }

        self.npwd = ko.observable().extend({ required: true, minLength: 3, maxLength: 18 });
        self.npwd1 = ko.observable().extend({ required: true, minLength: 3, maxLength: 18 });
        self.index_npwd_validationModel = ko.validatedObservable({
            npwd: self.npwd,
            npwd1: self.npwd1
        });

        self.changepwdsave = function () {
            var d = self.npwd();
            var dd = self.npwd1();
            if (self.index_npwd_validationModel.isValid()) {
                if (self.npwd() !== self.npwd1()) {
                    bootbox.alert("两次新密码不一致！");
                } else {
                    if (/.*[\u4e00-\u9fa5]+.*$/.test(self.npwd()) || /.*[\u4e00-\u9fa5]+.*$/.test(self.npwd1())) {
                        bootbox.alert("登陆错误，密码中不能含有中文！");
                        return
                    }
                    App.blockUI({
                        boxed: true
                    });
                    $.ajax({
                        method: "PUT",
                        url: basicurl + '/api/user/' + $.session.get('uid'),
                        data: JSON.stringify({ Uid: $.session.get('uid'), Pwd: self.npwd()}),
                        success: function (result) {
                            App.unblockUI();
                            if (result.ok == 1) {
                                self.npwd(''); self.npwd1('');
                                $.session.clear();
                                $.session.set('uid', result.data.Uid);
                                $.session.set('pass', "Basic " + base64.encode(result.data.Uid + ":" + result.data.Pwd));
                                bootbox.alert("修改密码成功，您的新密码是：" + result.data.Pwd);
                            } else {
                                bootbox.alert(result.err);
                            }
                        },
                        error: function (xhr, status, error) {
                            App.unblockUI();
                            bootbox.alert("错误提示：" + xhr.responseText);
                        }
                    })
                }
            } else {
                bootbox.alert("用户名密码长度错误！");
            }
        }

        IndexVM.prototype.onNavigatedTo = function () {

        }
        return IndexVM;
    })(BaseVM);
    viewModels.IndexVM = IndexVM;

    var usersVM = (function (_super) {
        __extends(usersVM, _super);
        var self = this;
        function usersVM() { _super.call(self); }
        usersVM.prototype.onNavigatedTo = function () {
            self.LoadUsers()
        }

        self.mg_user_users = ko.observableArray();
        self.LoadUsers = function () {
            App.blockUI({
                boxed: true
            });
            $.ajax({
                method: "GET",
                url: basicurl + '/api/um/all',
                success: function (result) {
                    App.unblockUI();
                    if (result.ok == 1) {
                        self.mg_user_users(result.data);
                    } else {
                        bootbox.alert(result.err);
                    }
                },
                error: function (xhr, status, error) {
                    App.unblockUI();
                    bootbox.alert(xhr,err);
                }
            })
        }

        self.addUser = function () {
            um_UserName("");
            um_UserId("")
            um_PassWord("")
            um_Email("")
            um_ukey("")
            $('#m_user_add').modal('toggle');
        }

        self.um_UserName = ko.observable().extend({ minLength: 2, maxLength: 18 });
        self.um_UserId = ko.observable().extend({ required: true, minLength: 3, maxLength: 18 });
        self.um_PassWord = ko.observable().extend({ required: true, minLength: 3, maxLength: 18 });
        self.um_Email = ko.observable().extend({ minLength: 3, maxLength: 128 });
        self.um_ukey = ko.observable();
        self.users_newuser_validationModel = ko.validatedObservable({
            um_UserName: self.um_UserName,
            um_UserId: self.um_UserId,
            um_PassWord: self.um_PassWord,
            um_Email: self.um_Email
        });
        self.um_ev_adduser = function () {
            if (self.users_newuser_validationModel.isValid()) {
                if (/.*[\u4e00-\u9fa5]+.*$/.test(self.um_UserId()) || /.*[\u4e00-\u9fa5]+.*$/.test(self.um_PassWord())) {
                    bootbox.alert("登陆错误，用户ID或密码中不能含有中文！");
                    return
                }
                App.blockUI({
                    boxed: true
                });
                var nuser = new Object();
                nuser.UserName = self.um_UserName();
                nuser.Uid = self.um_UserId();
                nuser.Pwd = self.um_PassWord();
                nuser.Email = self.um_Email();
                $.ajax({
                    method: "POST",
                    url: basicurl + '/api/user',
                    data: JSON.stringify(nuser),
                    success: function (result) {
                        App.unblockUI();
                        if (result.ok == 1) {
                            $('#m_user_add').modal('toggle');
                            self.LoadUsers();
                        } else {
                            bootbox.alert(result.err);
                        }
                    },
                    error: function (xhr, status, error) {
                        App.unblockUI();
                        bootbox.alert("错误，用户身份验证失败！");
                    }
                })
            } else {
                bootbox.alert("用户名密码长度错误！");
            }
        }

        self.editUser = function (u) {
            um_UserName(u.UserName);
            um_UserId(u.Uid)
            um_PassWord(u.Pwd)
            um_Email(u.Email)
            um_ukey(u.Ukey)
            $('#m_user_edit').modal('toggle');
        }

        self.um_ev_showedituser = function () {
            App.blockUI({
                boxed: true
            });
            var euser = new Object();
            euser.UserName = self.um_UserName();
            euser.Uid = self.um_UserId();
            euser.Pwd = self.um_PassWord();
            euser.Email = self.um_Email();
            euser.Ukey = self.um_ukey();
            $.ajax({
                method: "PUT",
                url: basicurl + '/api/user/' + euser.Uid,
                data: JSON.stringify(euser),
                success: function (result) {
                    App.unblockUI();
                    if (result.ok == 1) {
                        $('#m_user_edit').modal('toggle');
                        self.LoadUsers();
                    } else {
                        bootbox.alert(result.err);
                    }
                },
                error: function (xhr, status, error) {
                    App.unblockUI();
                    bootbox.alert("错误，用户身份验证失败！");
                }
            })
        }

        self.um_ev_delclick = function (u) {
            bootbox.confirm("你确定想删除吗？", function (result) {
                if (result) {
                    App.blockUI({
                        boxed: true
                    });
                    $.ajax({
                        method: "DELETE",
                        url: basicurl + '/api/user/' + u.Uid,
                        success: function (result) {
                            App.unblockUI();
                            if (result.ok == 1) {
                                self.LoadUsers();
                            } else {
                                bootbox.alert("错误：" + result.err);
                            }
                        },
                        error: function (xhr, status, error) {
                            App.unblockUI();
                            bootbox.alert("错误：" +xhr.responseText);
                        }
                    })
                }
            });
        };

        return usersVM;
    })(BaseVM);
    viewModels.usersVM = usersVM;

})(viewModels || (viewModels = {}));