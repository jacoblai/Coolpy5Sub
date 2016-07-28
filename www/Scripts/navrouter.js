var routing;
(function (routing) {
    var HashService = (function () {
        function HashService() {
            this.preventNextEvent = false;
            this.pending = false;
            this.cancellingPrev = false;
            this.callCount = 0;
            this.forwardingCount = 2;
            this.hash = "";
            this.on_changing = null;
            this.on_changed = null;
            this.on_cancelledByUrl = null;
        }
        HashService.prototype.setHash = function (hash) {
            window.location.hash = hash;
        };

        HashService.prototype.setHashAsReplace = function (hash) {
            window.location.replace(hash);
        };

        HashService.prototype.start = function () {
            this.storedHash = window.location.hash;
            if ("onhashchange" in window) {
                window.onhashchange = this.onHashChangedEventHandler.bind(this);
            } else {
                this.storedHash = window.location.hash;
                window.setInterval(function () {
                    if (window.location.hash != this.storedHash) {
                        this.onHashChangedEventHandler();
                    }
                }, 77);
            }

            if (window.location.hash) {
                this.hashChanged(window.location.hash);
            }
        };

        HashService.prototype.lock = function () {
            this.pending = true;
            this.cancellingPrev = false;
        };

        HashService.prototype.release = function () {
            this.pending = false;
            this.cancellingPrev = false;
        };

        HashService.prototype.changingCallback = function (cancelNavigation) {
            if (cancelNavigation) {
                this.preventNextEvent = true;
            }
        };

        HashService.prototype.hashChanged = function (newHash) {
            var _this = this;
            var continueHashChanged = function () {
                if (!_this.hash) {
                    _this.prevHash = newHash;
                } else {
                    _this.prevHash = _this.hash;
                }

                _this.hash = newHash;

                if (_this.on_changed) {
                    _this.on_changed(newHash);
                }

                _this.release();
            };

            this.lock();
            this.callCount++;
            if (this.on_changing) {
                var currentCount = this.callCount;
                this.on_changing(newHash, function (cancelNavigation) {
                    if (currentCount != _this.callCount) {
                        return;
                    }

                    if (cancelNavigation) {
                        _this.release();
                        _this.preventNextEvent = true;
                        history.back();

                        return;
                    }

                    continueHashChanged();
                });
                return;
            }

            continueHashChanged();
        };

        HashService.prototype.onHashChangedEventHandler = function () {
            if (this.pending) {
                if (this.prevHash == window.location.hash) {
                    this.release();
                    this.storedHash = window.location.hash;
                    this.preventNextEvent = false;
                    if (this.on_cancelledByUrl) {
                        this.on_cancelledByUrl();
                    }
                } else {
                    if (this.cancellingPrev) {
                        this.cancellingPrev = false;
                    } else {
                        this.cancellingPrev = true;
                        history.back();
                    }
                }

                return;
            }

            if (this.preventNextEvent) {
                if (this.hash != window.location.hash) {
                    if (this.forwardingCount == 0) {
                        throw new Error("History was broken, please reload page.");
                    }

                    this.forwardingCount--;
                    history.forward();
                } else {
                    this.storedHash = window.location.hash;
                    this.preventNextEvent = false;
                    this.forwardingCount = 2;
                }

                return;
            }

            if (this.storedHash == window.location.hash) {
                return;
            }

            this.storedHash = window.location.hash;
            this.hashChanged(window.location.hash);
        };
        return HashService;
    })();
    routing.HashService = HashService;
    ;
})(routing || (routing = {}));
var routing;
(function (routing) {
    var DefaultRouterLogger = (function () {
        function DefaultRouterLogger() {
        }
        DefaultRouterLogger.prototype.warning = function (message) {
            this.write("Router [Warning] >> " + message);
        };

        DefaultRouterLogger.prototype.error = function (message) {
            this.write("Router [Error]!  >> " + message);
        };

        DefaultRouterLogger.prototype.info = function (message) {
            this.write("Router [Info]    >> " + message);
        };

        DefaultRouterLogger.prototype.write = function (message) {
            if (typeof console == "undefined") {
                return;
            }

            console.log(message);
        };
        return DefaultRouterLogger;
    })();
    routing.DefaultRouterLogger = DefaultRouterLogger;

    var SilentLogger = (function () {
        function SilentLogger() {
        }
        SilentLogger.prototype.warning = function (message) {
        };

        SilentLogger.prototype.error = function (message) {
        };

        SilentLogger.prototype.info = function (message) {
        };
        return SilentLogger;
    })();
    routing.SilentLogger = SilentLogger;
})(routing || (routing = {}));
var __extends = this.__extends || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    __.prototype = b.prototype;
    d.prototype = new __();
};
var routing;
(function (routing) {
    (function (routes) {
        var Route = (function () {
            function Route(routePattern, options) {
                if (!routePattern) {
                    throw new Error("Route pattern should be specified!");
                }

                this.parrentRoute = null;
                this.pattern = routePattern || null;
                if (!options) {
                    this.isDefault = false;
                    this.canLeave = function (callback) {
                        callback(true);
                    };
                    return;
                }

                this.isDefault = options.isDefault || false;
                this.canLeave = options.canLeave || (function (callback, navOptions) {
                    callback(true);
                });
            }
            return Route;
        })();
        routes.Route = Route;

        var VirtualRoute = (function (_super) {
            __extends(VirtualRoute, _super);
            function VirtualRoute(routePattern, childRoutes, options) {
                this.childRoutes = childRoutes || new Array();
                _super.call(this, routePattern, options);
            }
            return VirtualRoute;
        })(Route);
        routes.VirtualRoute = VirtualRoute;
        ;

        (function (LoadingState) {
            LoadingState[LoadingState["canceled"] = 0] = "canceled";
            LoadingState[LoadingState["complete"] = 1] = "complete";
            LoadingState[LoadingState["loading"] = 2] = "loading";
        })(routes.LoadingState || (routes.LoadingState = {}));
        var LoadingState = routes.LoadingState;

        var NavigationRoute = (function (_super) {
            __extends(NavigationRoute, _super);
            function NavigationRoute(routePattern, viewPath, options) {
                var _this = this;
                if (!viewPath) {
                    throw new Error("Route view path should be specified!");
                }

                this.viewPath = viewPath;
                this.currentVM = null;

                if (!options) {
                    this.cacheView = true;
                    this.vmFactory = null;
                    this.onNavigatedTo = null;
                    this.title = null;
                    this.toolbarId = null;
                    this.state = 1 /* complete */;
                } else {
                    if (options.vmFactory) {
                        var factory = eval(options.vmFactory);
                        this.vmFactory = function (callback) {
                            factory(function (instance) {
                                _this.currentVM = instance;
                                _this.onNavigatedTo = instance.onNavigatedTo || null;
                                _this.onNavigatedFrom = instance.onNavigatedFrom || null;
                                _this.canLeave = function (callback, navOptions) {
                                    if (instance.canNavigateFrom) {
                                        instance.canNavigateFrom(callback, navOptions);
                                        return;
                                    }

                                    callback(true);
                                };

                                callback(instance);
                            });
                        };
                    }

                    this.toolbarId = options.toolbarId;
                    this.cacheView = options.cacheView == undefined ? true : options.cacheView;
                    this.title = options.title || null;
                    this.state = 1 /* complete */;
                }

                _super.call(this, routePattern, options);
            }
            return NavigationRoute;
        })(Route);
        routes.NavigationRoute = NavigationRoute;
        ;
    })(routing.routes || (routing.routes = {}));
    var routes = routing.routes;
})(routing || (routing = {}));
var routing;
(function (routing) {
    (function (utils) {
        function getType(obj) {
            var funcNameRegex = /function (.+)\(/;
            var results = (funcNameRegex).exec((obj).constructor.toString());
            return (results && results.length > 1) ? results[1] : "";
        }
        utils.getType = getType;
        ;

        function getHash(path) {
            if (typeof path != "String" && path.toString != "undefined") {
                path = path.toString();
            }

            var matches = path.match(/^[^#]*(#.+)$/);
            var hash = matches ? matches[1] : '';
            return hash;
        }
        utils.getHash = getHash;
        ;
    })(routing.utils || (routing.utils = {}));
    var utils = routing.utils;
})(routing || (routing = {}));

var routing;
(function (routing) {
    var RouteHandler = (function () {
        function RouteHandler(pattern, handler) {
            this.pattern = pattern || null;
            this.handler = handler || null;
        }
        return RouteHandler;
    })();
    ;

    var Router = (function () {
        function Router(mainContainerId, options, routesToMap) {
            var _this = this;
            this.initialized = false;
            this.routes = new Array();
            this.currentRoute = this.createCurrentRoute();
            this.history = new Array();
            this.hashSymbol = "#!/";
            this.defaultPath = "";
            this.currentHash = "";
            this.startupUrl = "";
            this.containerId = "";
            this.defaultRoute = null;
            this.fresh = true;
            this.allRoutes = new Array();
            this.handlers = new Array();
            this.hashService = new routing.HashService();
            this.currentLogger = new routing.SilentLogger();
            this.currentPayload = null;
            this.navigationFlags = null;
            this.forceReloadOnNavigation = false;
            this.forceNavigationInCache = false;
            this.forceCaching = false;
            this.backNavigation = false;
            this.isRedirecting = false;
            this.preventRaisingNavigateTo = false;
            this.beforeNavigationHandler = null;
            this.afterNavigationHandler = null;
            this.navigationErrorHandler = null;
            this.cancelledByUrlHandler = null;
            this.registerRoutes = function (routesToMap) {
                for (var i in routesToMap) {
                    this.registerRoute(routesToMap[i]);
                }

                this.defaultPath = this.hashSymbol + this.getPathForRoute(this.defaultRoute);
            };
            this.init = function (routes, mainContainerId, options) {
                var enableLogging;
                this.hashService.on_changing = this.hashChanginHandler;
                this.hashService.on_changed = this.hashChangedHandler;
                this.hashService.on_cancelledByUrl = this.hashChangeCancelledHandler;

                if (options) {
                    this.forceCaching = options.preloadEnabled || false;
                    this.onPreloadComplete = options.preloadComplete || null;
                    this.beforeNavigationHandler = options.beforeNavigation || null;
                    this.afterNavigationHandler = options.afterNavigation || null;
                    this.navigationErrorHandler = options.navigationError || null;
                    enableLogging = options.enableLogging || true;
                }

                this.currentLogger = enableLogging ? new routing.DefaultRouterLogger() : new routing.SilentLogger();
                this.containerId = mainContainerId;
                this.initialized = true;
                this.registerRoutes(routes);
                this.currentLogger.info("Initialized.");
                return this;
            };
            this.hashChanginHandler = function (hash, callback) {
                if (!_this.getRoute(hash)) {
                    callback(true);
                    _this.currentLogger.error("Navigation to '" + hash + "' was prevented. The route to this pattern was not found.");
                    return;
                }

                if (_this.currentRoute() == null) {
                    callback(false);
                    return;
                }

                if (!_this.isRedirecting) {
                    _this.currentRoute().canLeave(function (accept) {
                        if (accept) {
                            _this.isRedirecting = false;
                            _this.currentHash = hash;
                            callback(false);
                            _this.preventRaisingNavigateTo = false;
                        } else {
                            _this.backNavigation = false;
                            callback(true);
                        }

                        _this.forceReloadOnNavigation = false;
                        _this.forceNavigationInCache = false;
                    }, {
                        targetRoute: _this.getRoute(hash),
                        forceReloadOnNavigation: _this.forceReloadOnNavigation,
                        forceNavigationInCache: _this.forceNavigationInCache
                    });
                } else {
                    _this.isRedirecting = false;
                    _this.currentHash = hash;
                    callback(false);
                    _this.preventRaisingNavigateTo = false;
                }
            };
            this.hashChangedHandler = function (hash) {
                var route = _this.getRoute(hash);
                var context = _this.getContext(route, hash);
                var routeHandler;
                var delegate = function (x) {
                    return x.pattern == route.pattern;
                };
                for (var i = 0; i < _this.handlers.length; i++) {
                    if (delegate(_this.handlers[i])) {
                        routeHandler = _this.handlers[i];
                    }
                }

                var croute = _this.currentRoute();
                if (croute && croute instanceof routing.routes.NavigationRoute) {
                    croute.state = 0 /* canceled */;
                    if (croute.onNavigatedFrom) {
                        croute.onNavigatedFrom();
                    }
                }

                routeHandler.handler(context);

                if (!_this.preventRaisingNavigateTo) {
                    _this.currentLogger.info("Navigated to '" + hash + "'.");
                } else {
                    _this.currentLogger.info("Navigion was prevented.");
                }

                _this.refreshCurrentRoute();
            };
            this.hashChangeCancelledHandler = function () {
                _this.forceReloadOnNavigation = false;
                _this.forceNavigationInCache = false;
                if (_this.cancelledByUrlHandler) {
                    _this.cancelledByUrlHandler();
                }
            };
            var enableLogging;
            this.hashService.on_changing = this.hashChanginHandler;
            this.hashService.on_changed = this.hashChangedHandler;
            this.hashService.on_cancelledByUrl = this.hashChangeCancelledHandler;

            if (options) {
                this.beforeNavigationHandler = options.beforeNavigationHandler || null;
                this.afterNavigationHandler = options.afterNavigationHandler || null;
                this.navigationErrorHandler = options.navigationErrorHandler || null;
                enableLogging = options.enableLogging || true;
            }

            this.currentLogger = enableLogging ? new routing.DefaultRouterLogger() : new routing.SilentLogger();
            this.containerId = mainContainerId;
            this.initialized = true;
            if (routesToMap) {
                this.registerRoutes(routesToMap);
            }
        }
        Router.prototype.navigateTo = function (path, options) {
            var actualPath = path, relRoute = this.getRoute(path), removeCurrentHistory = false;

            if (options) {
                this.currentPayload = options.payload || null;
                this.forceReloadOnNavigation = options.forceReload || false;
                this.forceNavigationInCache = options.forceNavigationInCache || false;
                removeCurrentHistory = options.removeCurrentHistory || false;
            }

            if (relRoute && relRoute instanceof routing.routes.VirtualRoute) {
                actualPath = this.getPathForRoute(relRoute);
            }

            if (!(actualPath == this.currentHash || this.hashSymbol + actualPath == this.currentHash)) {
                actualPath = this.fixPath(actualPath);
                if (removeCurrentHistory) {
                    this.hashService.setHashAsReplace(actualPath);
                } else {
                    this.hashService.setHash(actualPath);
                }
            }
        };

        Router.prototype.navigateBack = function () {
            history.back();
        };

        Router.prototype.navigateBackInCache = function () {
            this.forceNavigationInCache = true;
            this.navigateBack();
        };

        Router.prototype.navigateHome = function () {
            this.navigateTo(this.startupUrl);
        };

        Router.prototype.getHashSymbol = function () {
            return this.hashSymbol;
        };

        Router.prototype.cancelledByUrl = function (handler) {
            this.cancelledByUrlHandler = handler;
        };

        Router.prototype.refreshCurrentRoute = function () {
            var pureHash = routing.utils.getHash(window.location.toString()).replace(this.hashSymbol, "");
            var route = this.getRoute(pureHash);
            if (route != null) {
                this.currentRoute(route);
            }
        };

        Router.prototype.registerRoute = function (routeToMap) {
            this.routes.push(routeToMap);
            if (routeToMap.isDefault) {
                this.defaultRoute = routeToMap;
                this.defaultPath = this.hashSymbol + this.getPathForRoute(routeToMap);
            }

            this.initRoute(routeToMap);
            return this;
        };

        Router.prototype.setLogger = function (logger) {
            if (!logger) {
                throw new Error("Parameter 'logger' is null or undefined!");
            }

            this.currentLogger = logger;
            return this;
        };

        Router.prototype.run = function () {
            if (!this.initialized) {
                throw new Error("Router is not initialized. Router should be initialized first!");
                return;
            }

            if (this.forceCaching) {
            }

            this.defaultTitle = document.title;
            this.currentLogger.info("Successfully started.");
            this.hashService.start();
            this.startupUrl = this.hashService.hash || this.defaultPath;
            if (this.startupUrl == this.defaultPath) {
                this.hashService.setHash(this.startupUrl);
            }

            this.currentHash = this.startupUrl;
            return this;
        };

        Router.prototype.getRoute = function (routeLink) {
            var _this = this;
            var delegate = function (x) {
                var path2 = routeLink.toString().replace(_this.hashSymbol, "");
                var result = _this.isMatchesV2(x.pattern, path2);
                return result;
            };
            var res = null;
            for (var i = 0; i < this.allRoutes.length; i++) {
                if (delegate(this.allRoutes[i])) {
                    res = this.allRoutes[i];
                    break;
                }
            }

            return res;
        };

        Router.prototype.isMatches = function (path1, path2) {
            var result = true, path1Parts = path1.toString().split("/"), path2Parts = path2.toString().split("/");
            if (path1Parts.length == path2Parts.length) {
                for (var i = 0; i < path1Parts.length; i++) {
                    if (!path1Parts[i].match(/^:.+/) && path1Parts[i] != path2Parts[i]) {
                        return false;
                    }
                }
            } else {
                result = false;
            }

            return result;
        };

        Router.prototype.isMatchesV2 = function (path1, path2) {
            var path1Parts = this.cleanPath(path1).split("/"), path2Parts = this.cleanPath(path2).split("/");
            if (path1Parts.length < path2Parts.length) {
                return false;
            }

            for (var i = 0; i < path1Parts.length; i++) {
                if (path1Parts[i].match(/^\{\?[^\?]+\}$/)) {
                    continue;
                }

                if (path1Parts[i].match(/^\{([^\?])+\}$/) && path2Parts[i]) {
                    continue;
                }

                if (path1Parts[i] == path2Parts[i]) {
                    continue;
                }

                return false;
            }

            return true;
        };

        Router.prototype.cleanPath = function (path) {
            return path.replace(/(\/\/+)/, "/").replace(/(\/+)$/, "").replace(/^(\/+)/, "");
        };

        Router.prototype.getPathForRoute = function (route) {
            if (route) {
                if (route instanceof routing.routes.VirtualRoute) {
                    var vroute = route;
                    var defaultChild = null;
                    for (var i = 0; i < vroute.childRoutes.length; i++) {
                        if (vroute.childRoutes[i].isDefault) {
                            defaultChild = vroute.childRoutes[i];
                            break;
                        }
                    }

                    if (defaultChild == null) {
                        throw new Error("Route '" + route.pattern + "' has invalid configuration of child elements.");
                    }

                    return this.getPathForRoute(defaultChild);
                }

                return route.pattern;
            }

            return null;
        };

        Router.prototype.getCompletePath = function (path, params) {
            var matches = path.toString().match(/\{.+\}/);
            var completePath = path.toString();
            if (matches) {
                for (var i = 0; i < matches.length; i++) {
                    var paramName = matches[i].toString().replace("{", "").replace("}", "");
                    completePath = completePath.replace("{" + paramName + "}", params[paramName]);
                }
            }

            return completePath;
        };

        Router.prototype.fixPath = function (path) {
            if (!path.match(/^/ + this.hashSymbol + /.+/)) {
                return this.hashSymbol + path.replace("#/", "");
            }
        };

        Router.prototype.createCurrentRoute = function () {
            return ko.observable(null);
        };

        Router.prototype.raiseOnNavigatedTo = function (route, context) {
            if (route.onNavigatedTo != null && (!this.isRedirecting || !this.preventRaisingNavigateTo)) {
                var params = context.params;
                route.onNavigatedTo(params, this.currentPayload);
                this.currentPayload = null;
            }
        };

        Router.prototype.getContext = function (route, hash) {
            var context = {
                associeatedRoute: route,
                path: hash.replace(this.hashSymbol, "")
            };
            var params = {};
            var patternParts = route.pattern.split("/");
            var pathParts = hash.replace(this.hashSymbol, "").split("/");

            for (var i = 0; i < patternParts.length; i++) {
                if (patternParts[i].toString().match(/^\{[^\?]+\}$/)) {
                    var paramName = patternParts[i].toString().replace("{", "").replace("}", "");
                    params[paramName] = pathParts[i];
                }

                if (patternParts[i].toString().match(/^\{\?[^\?]+\}$/)) {
                    var paramName = patternParts[i].toString().replace("{?", "").replace("}", "");
                    params[paramName] = pathParts[i];
                }
            }

            context.params = params;
            return context;
        };

        Router.prototype.mapVirtualRoute = function (routeToMap) {
            if (routeToMap.childRoutes) {
                for (var i in routeToMap.childRoutes) {
                    routeToMap.childRoutes[i].parrentRoute = routeToMap;
                    routeToMap.childRoutes[i].pattern = routeToMap.pattern + "/" + routeToMap.childRoutes[i].pattern;
                    this.initRoute(routeToMap.childRoutes[i]);
                }
            }
        };

        Router.prototype.mapNavigationRoute = function (routeToMap) {
            var _this = this;
            this.handlers.push(new RouteHandler(routeToMap.pattern, function (context) {
                function completeNavigation() {
                    context.associeatedRoute.state = 1 /* complete */;
                    if (this.afterNavigationHandler) {
                        this.afterNavigationHandler();
                    }

                    if (this.fresh) {
                        this.fresh = false;
                    }
                }
                ;

                function onNavigationError() {
                    if (this.navigationErrorHandler) {
                        this.currentLogger.warning("Navigation error is handling...");
                        this.navigationErrorHandler();
                    }
                }
                ;

                if (_this.beforeNavigationHandler) {
                    _this.beforeNavigationHandler();
                }

                context.associeatedRoute.state = 2 /* loading */;
                var jelem = $("#" + _this.containerId);
                var completePath = _this.getCompletePath(routeToMap.viewPath, context.params);
                var existing = $("[data-view=\"" + routeToMap.pattern + "\"]", jelem);
                var preventRaisingNavigateToCache = _this.preventRaisingNavigateTo;

                if (routeToMap.title) {
                    document.title = routeToMap.title;
                } else {
                    document.title = _this.defaultTitle;
                }

                if (existing && existing.length >= 1) {
                    if ((routeToMap.cacheView || _this.forceNavigationInCache) && !_this.forceReloadOnNavigation) {
                        if (_this.forceNavigationInCache) {
                            _this.forceNavigationInCache = false;
                        }

                        jelem.children().hide();
                        existing.show();
                        if (!preventRaisingNavigateToCache) {
                            _this.raiseOnNavigatedTo(routeToMap, context);
                        }

                        completeNavigation();
                    } else if (!preventRaisingNavigateToCache) {
                        if (_this.forceReloadOnNavigation) {
                            _this.forceReloadOnNavigation = false;
                        }

                        $.ajax({
                            url: completePath,
                            data: null,
                            cache: false,
                            error: onNavigationError,
                            success: function (response) {
                                if (routeToMap.state == 0 /* canceled */) {
                                    _this.currentLogger.warning("Navigation to " + context.path + " was cancelled!");
                                    return;
                                }

                                if (routeToMap.vmFactory != null) {
                                    if (existing && existing.get(0)) {
                                        ko.cleanNode(existing.get(0));
                                    }
                                }

                                existing.html(response);
                                if (routeToMap.vmFactory != null) {
                                    var factory = routeToMap.vmFactory;
                                    factory(function (instance) {
                                        ko.applyBindings(instance, existing.get(0));
                                        routeToMap.currentVM = instance;
                                        jelem.children().hide();
                                        if (!preventRaisingNavigateToCache) {
                                            _this.raiseOnNavigatedTo(routeToMap, context);
                                        }

                                        existing.show();
                                        completeNavigation();
                                    });
                                } else {
                                    routeToMap.currentVM = null;
                                    jelem.children().hide();
                                    if (!preventRaisingNavigateToCache) {
                                        _this.raiseOnNavigatedTo(routeToMap, context);
                                    }

                                    existing.show();
                                    completeNavigation();
                                }
                            }
                        });
                    } else {
                        completeNavigation();
                    }
                } else {
                    $.ajax({
                        url: completePath,
                        data: null,
                        cache: false,
                        error: onNavigationError,
                        success: function (response) {
                            if (routeToMap.state == 0 /* canceled */) {
                                _this.currentLogger.warning("Navigation to " + context.path + " were cancelled!");
                                return;
                            }

                            jelem.children().hide();
                            jelem.append("<div data-view=\"" + routeToMap.pattern + "\">" + response + "</div>");
                            if (routeToMap.vmFactory != null) {
                                existing = $("[data-view=\"" + routeToMap.pattern + "\"]", jelem);
                                var factory = routeToMap.vmFactory;
                                factory(function (instance) {
                                    ko.applyBindings(instance, existing.get(0));
                                    routeToMap.currentVM = instance;
                                    if (!preventRaisingNavigateToCache) {
                                        _this.raiseOnNavigatedTo(routeToMap, context);
                                    }

                                    completeNavigation();
                                });
                            } else {
                                if (!preventRaisingNavigateToCache) {
                                    _this.raiseOnNavigatedTo(routeToMap, context);
                                }

                                existing.show();
                                completeNavigation();
                            }
                        }
                    });
                }
            }));
        };

        Router.prototype.initRoute = function (routeToMap) {
            this.allRoutes.push(routeToMap);
            this.currentLogger.info("Registering route '" + routeToMap.pattern + "'.");
            switch (routing.utils.getType(routeToMap)) {
                case "VirtualRoute":
                    this.mapVirtualRoute(routeToMap);
                    break;
                case "NavigationRoute":
                    this.mapNavigationRoute(routeToMap);
                    break;
            }
        };
        return Router;
    })();
    routing.Router = Router;
})(routing || (routing = {}));

var routing;
(function (routing) {
    (function (knockout) {
        var _router = null;

        function checkRouter() {
            if (_router == null || _router == undefined) {
                throw new Error("Router instance do not setted. Please set it usting 'Routing.ko.setCurrentRouter' method.");
            }
        }

        function isString(obj) {
            return typeof obj == "string" || obj instanceof String;
        }

        function setCurrentRouter(router) {
            _router = router;
        }
        knockout.setCurrentRouter = setCurrentRouter;

        ko.bindingHandlers.navigate = {
            init: function (element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
                var $elem = $(element), bindings = allBindingsAccessor(), navLink = valueAccessor(), payload = null, forceReloadOnNavigation = bindings.forceReload || false, forceNavigationInCache = bindings.forceNavigationInCache || false, oldClass;

                payload = bindings.payload || null;
                if (element.tagName == "A" && payload == null && false) {
                    $elem.attr("href", "#!/" + navLink);
                } else {
                    if (element.tagName == "A") {
                        $elem.attr("href", "#");
                    }

                    $elem.click(function (event) {
                        event.preventDefault();
                        if (_router.initialized) {
                            if (payload == null) {
                                _router.navigateTo(navLink, {
                                    removeCurrentHistory: false,
                                    forceReload: ko.utils.unwrapObservable(forceReloadOnNavigation),
                                    forceNavigationInCache: ko.utils.unwrapObservable(forceNavigationInCache)
                                });
                            } else {
                                _router.navigateTo(navLink, {
                                    payload: ko.utils.unwrapObservable(payload),
                                    removeCurrentHistory: false,
                                    forceReload: ko.utils.unwrapObservable(forceReloadOnNavigation),
                                    forceNavigationInCache: ko.utils.unwrapObservable(forceNavigationInCache)
                                });
                            }
                        }
                    });
                }

                var checkChilds = function (path, route) {
                    if (route instanceof routing.routes.VirtualRoute && route.childRoutes.length > 0) {
                        for (var i in route.childRoutes) {
                            if (_router.isMatches(route.childRoutes[i].pattern, path)) {
                                return true;
                            } else if (checkChilds(path, route.childRoutes[i])) {
                                return true;
                            }
                        }
                    }

                    return false;
                };

                _router.currentRoute.subscribe(function () {
                    var currentClass = $elem.attr("class");
                    if (bindings.activeClass) {
                        var path = routing.utils.getHash(window.location).replace(_router.getHashSymbol(), "");
                        if (path == navLink || checkChilds(path, _router.getRoute(navLink))) {
                            if (!$elem.hasClass(bindings.activeClass)) {
                                oldClass = currentClass || null;
                                $elem.addClass(bindings.activeClass);
                            }
                        } else {
                            if ($elem.hasClass(bindings.activeClass)) {
                                $elem.removeClass(bindings.activeClass);
                            }
                        }
                    }
                });
            }
        };

        ko.bindingHandlers.navigateBack = {
            init: function (elem, valueAccessor) {
                var $elem = $(elem), options = valueAccessor(), forceNavigationInCache = options.forceNavigationInCache || false;

                $elem.click(function () {
                    if (forceNavigationInCache) {
                        _router.navigateBackInCache();
                    } else {
                        _router.navigateBack();
                    }
                });
            }
        };
    })(routing.knockout || (routing.knockout = {}));
    var knockout = routing.knockout;
})(routing || (routing = {}));
