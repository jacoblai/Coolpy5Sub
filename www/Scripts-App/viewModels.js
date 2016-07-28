if (localStorage.getItem('token') == null || localStorage.getItem('username') == null) {
    window.location.href = "login.html";
}

window.onbeforeunload = function () {
    localStorage.clear();
}

$.ajaxSetup({
    cache: false,
    beforeSend: function (xhr) {
        xhr.setRequestHeader("Authorization", "Basic " + localStorage.getItem('token'));
    }
});

function myalert(msg) {
    bootbox.alert(msg);
}

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
        self.uid = ko.observable();
        self.Message = ko.observable();
        self.pps = ko.observableArray();

        self.opwd = ko.observable().extend({ required: true, minLength: 3, maxLength: 64 });
        self.npwd = ko.observable().extend({ required: true, minLength: 3, maxLength: 64 });
        self.npwd1 = ko.observable().extend({ required: true, minLength: 3, maxLength: 64 });
        self.validationModel = ko.validatedObservable({
            opwd: self.opwd,
            npwd: self.opwd,
            npwd1: self.opwd
        });

        self.username = ko.observable(localStorage.getItem('username'));
        self.logout = function () {
            localStorage.clear();
            window.location.href = "login.html";
        }
        self.changepwdshowwin = function () {
            self.opwd('');
            self.npwd('');
            self.npwd1('');
            $('#m_changepwd').modal('toggle');
        }
        self.changepwdsave = function () {
            if (self.validationModel.isValid()) {
                if (self.npwd() !== self.npwd1()) {
                    myalert("两次新密码不一致！");
                } else {
                    Metronic.blockUI({
                        boxed: true
                    });
                    $.ajax({
                        url: url + "/total/changepwd/" + self.username() + "/" + self.opwd() + "/" + self.npwd()
                    }).done(function (data) {
                        if (data === "密码修改成功") {
                            self.opwd('');
                            self.npwd('');
                            self.npwd1('');
                            self.logout();
                        } else {
                            myalert(data);
                        }
                        Metronic.unblockUI();
                    });
                }
            } else {
                myalert("用户名密码长度错误！");
            }
        }

        self.doview = function (dv) {
            gotoview("pages/pl/", dv.p_身份证号, new Array());
            self.uid(dv.p_身份证号);
            $('#m_basic').modal('toggle');
        };
        self.more = function () {
            gotoview("pages/more", "", new Array());
        };
        self.serach = function () {
            if (self.uid() !== undefined && self.uid() !== "") {
                Metronic.blockUI({
                    boxed: true
                });
                $.ajax({
                    url: url + "/api/persons/pname/" + encodeURIComponent(self.uid())
                }).done(function (data) {
                    if (data === "0") {
                        myalert("没有查找到相关自然人信息！");
                    } else {
                        self.pps.removeAll();
                        ko.utils.arrayForEach(data, function (ss) {
                            self.pps.push(ss);
                        });
                        $('#m_basic').modal('toggle');
                    }
                    Metronic.unblockUI();
                }).fail(function (err) {
                    Metronic.unblockUI();
                    myalert(err);
                });
            }
            else {
                myalert("请输入身份证号码或姓名！");
            }
        };

        self.loadData = function () {
            $.ajax({
                url: url + "/total/gszrsl"
            }).done(function (data) {
                self.Message(data[0].Key + ":" + data[0].Value + "," + data[1].Key + ":" + data[1].Value + "," + data[2].Key + ":" + data[2].Value);
            });
        };
        self.loadData();
        return BaseVM;
    })();
    viewModels.BaseVM = BaseVM;

    var IndexVM = (function (_super) {
        __extends(IndexVM, _super);
        var self = this;
        function IndexVM() { _super.call(self); }
        IndexVM.prototype.onNavigatedTo = function () {
            var initChartAge = function () {
                var myChart = echarts.init(document.getElementById('chart_Age'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/glldbzsl"
                }).done(function (data) {
                    var itemLengths = [];

                    itemLengths.push(data["18岁以下"]);
                    itemLengths.push(data["18岁至30岁"]);
                    itemLengths.push(data["30岁至40岁"]);
                    itemLengths.push(data["40岁至50岁"]);
                    itemLengths.push(data["50岁至65岁"]);
                    itemLengths.push(data["65岁以上"]);

                    var option = {
                        title: {
                            show: false,
                            text: '办事频率',
                            subtext: '纯属虚构'
                        },
                        tooltip: {
                            trigger: 'axis'
                        },
                        legend: {
                            show: false,
                            data: ['办事频率', '最低气温']
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                mark: { show: false },
                                dataView: { show: true, readOnly: false },
                                magicType: { show: false, type: ['line', 'bar'] },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        calculable: true,
                        xAxis: [
                            {
                                type: 'category',
                                boundaryGap: false,
                                data: ['18岁以下', '18岁至30岁', '30岁至40岁', '40岁至50岁', '50岁至65岁', '65岁以上']
                            }
                        ],
                        yAxis: [
                            {
                                type: 'value',
                                axisLabel: {
                                    formatter: '{value}'
                                }
                            }
                        ],
                        series: [
                            {
                                name: '办事频率',
                                type: 'line',
                                //data: itemLengths,
                                data: [42787, 151282, 200205, 203328, 290042, 134037],
                                markPoint: {
                                    data: [
                                        { type: 'max', name: '最大值' },
                                        { type: 'min', name: '最小值' }
                                    ]
                                },
                                markLine: {
                                    data: [
                                        { type: 'average', name: '平均值' }
                                    ]
                                }
                            }
                        ]
                    };

                    myChart.hideLoading();
                    myChart.setOption(option);
                });
            }
            var initChartServiceFrequency = function () {
                var myChart = echarts.init(document.getElementById('chart_ServiceFrequency'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/bspl"
                }).done(function (data) {

                    var items;

                    items = data;
                    document.getElementById('A').innerHTML = 'A.' + items[0].Key;
                    document.getElementById('B').innerHTML = 'B.' + items[1].Key;
                    document.getElementById('C').innerHTML = 'C.' + items[2].Key;
                    document.getElementById('D').innerHTML = 'D.' + items[3].Key;
                    document.getElementById('E').innerHTML = 'E.' + items[4].Key;
                    document.getElementById('F').innerHTML = 'F.' + items[5].Key;
                    document.getElementById('G').innerHTML = 'G.' + items[6].Key;
                    document.getElementById('H').innerHTML = 'H.' + items[7].Key;
                    document.getElementById('I').innerHTML = 'I.' + items[8].Key;
                    document.getElementById('J').innerHTML = 'J.' + items[9].Key;
                    //获取总数
                    var allCount = 0;
                    $.each(items, function () {
                        allCount += this.Value;
                    });

                    //获取其他总数
                    var otherCount = allCount;
                    for (var i = 0; i < 9; i++) {
                        otherCount -= items[i].Value;
                    }

                    var labelTop = {
                        normal: {
                            label: {
                                show: true,
                                position: 'center',
                                formatter: '{b}',
                                textStyle: {
                                    baseline: 'bottom'
                                }
                            },
                            labelLine: {
                                show: false
                            }
                        }
                    };
                    var labelFromatter = {
                        normal: {
                            label: {
                                formatter: function (params) {
                                    return Math.round(params.series.value / allCount * 10000) / 100.00 + "%";
                                },
                                textStyle: {
                                    baseline: 'top'
                                }
                            }
                        }
                    }
                    var labelBottom = {
                        normal: {
                            color: '#ccc',

                            label: {
                                show: true,
                                position: 'center'
                            },
                            labelLine: {
                                show: false
                            }
                        },
                        emphasis: {
                            color: 'rgba(0,0,0,0)'
                        }
                    };
                    var radius = [40, 55];
                    var option = {
                        legend: {
                            show: false,
                            x: 0,
                            y: 55,
                            data: [
                                "A", "B", "C", "D", "E",
                                "F", "G", "H", "I", "J"
                            ]
                        },
                        title: {
                            show: false,
                            text: 'The App World',
                            subtext: 'from global web index',
                            x: 'center'
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                dataView: { show: true, readOnly: false },
                                magicType: {
                                    show: false,
                                    type: ['pie', 'funnel'],
                                    option: {
                                        funnel: {
                                            width: '20%',
                                            height: '30%',
                                            itemStyle: {
                                                normal: {
                                                    label: {
                                                        formatter: function (params) {
                                                            return 'other\n' + params.value + '%\n'
                                                        },
                                                        textStyle: {
                                                            baseline: 'middle'
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        series: [
                            {
                                type: 'pie',
                                center: ['10%', '30%'],
                                radius: radius,
                                x: '0%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[0].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "A", value: items[0].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['30%', '30%'],
                                radius: radius,
                                x: '20%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[1].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "B", value: items[1].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['50%', '30%'],
                                radius: radius,
                                x: '40%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[2].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "C", value: items[2].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['70%', '30%'],
                                radius: radius,
                                x: '60%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[3].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "D", value: items[3].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['90%', '30%'],
                                radius: radius,
                                x: '80%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[4].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "E", value: items[4].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['10%', '70%'],
                                radius: radius,
                                y: '55%', // for funnel
                                x: '0%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[5].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "F", value: items[5].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['30%', '70%'],
                                radius: radius,
                                y: '55%', // for funnel
                                x: '20%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[6].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "G", value: items[6].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['50%', '70%'],
                                radius: radius,
                                y: '55%', // for funnel
                                x: '40%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[7].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "H", value: items[7].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['70%', '70%'],
                                radius: radius,
                                y: '55%', // for funnel
                                x: '60%', // for funnel
                                itemStyle: labelFromatter,
                                value: items[8].Value,
                                data: [
                                    { name: '总办理量', value: allCount, itemStyle: labelBottom },
                                    { name: "I", value: items[8].Value, itemStyle: labelTop }
                                ]
                            },
                            {
                                type: 'pie',
                                center: ['90%', '70%'],
                                radius: radius,
                                y: '55%', // for funnel
                                x: '80%', // for funnel
                                itemStyle: labelFromatter,
                                value: otherCount,
                                data: [
                                    { name: "总办理量", value: allCount, itemStyle: labelBottom },
                                    { name: "J", value: otherCount, itemStyle: labelTop }
                                ]
                            }
                        ]
                    };
                    myChart.hideLoading();
                    myChart.setOption(option);
                });
            }
            var initChart1 = function () {
                var myChart = echarts.init(document.getElementById('chart_cl'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/szrkzzck"
                }).done(function (data) {
                    var itemLengths = [];
                    var itemLengthskeys = [];
                    Object.keys(data).forEach(function (key) {
                        itemLengths.push(data[key]);
                        itemLengthskeys.push(key);
                    });
                    var option = {
                        title: {
                            show: false,
                            x: 'center',
                            text: '',
                            subtext: 'Rainbow bar example',
                            link: 'http://echarts.baidu.com/doc/example.html'
                        },
                        tooltip: {
                            trigger: 'item'
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                dataView: { show: true, readOnly: false },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        calculable: true,
                        grid: {
                            borderWidth: 0,
                            y: 80,
                            y2: 60
                        },
                        xAxis: [
                            {
                                type: 'category',
                                show: false,
                                data: itemLengthskeys
                            }
                        ],
                        yAxis: [
                            {
                                type: 'value',
                                show: false
                            }
                        ],
                        series: [
                            {
                                name: '发证残疾人数 ',
                                type: 'bar',
                                itemStyle: {
                                    normal: {
                                        color: function (params) {
                                            // build a color map as your need.
                                            var colorList = [
                                              '#C1232B', '#B5C334', '#FCCE10', '#E87C25', '#27727B',
                                               '#FE8463', '#9BCA63', '#FAD860', '#F3A43B', '#60C0DD',
                                               '#D7504B', '#C6E579', '#F4E001', '#F0805A', '#26C0C0'
                                            ];
                                            return colorList[params.dataIndex];
                                        },
                                        label: {
                                            show: true,
                                            position: 'top',
                                            formatter: '{b}\n{c}'
                                        }
                                    }
                                },
                                data: itemLengths
                            }
                        ]
                    };
                    myChart.hideLoading();
                    myChart.setOption(option);
                }
                );
            }
            var initChart2 = function () {
                var myChart = echarts.init(document.getElementById('chart_2'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/bsrrhzdfb"
                }).done(function (data) {
                    //var xx = data[0].Key;
                    var sss = data[0].Key;
                    var sss1 = data[0].Value["s_本市户籍受理数"];
                    var sss2 = data[0].Value["s_非本市户籍受理数"];

                    var option = {
                        title: {
                            show: false,
                            text: '某地区蒸发量和降水量',
                            subtext: '纯属虚构'
                        },
                        tooltip: {
                            trigger: 'axis'
                        },
                        legend: {
                            data: ['禅城户籍人口', '流动人口']
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                mark: { show: false },
                                dataView: { show: true, readOnly: false },
                                magicType: { show: false, type: ['line', 'bar'] },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        calculable: true,
                        xAxis: [
                            {
                                type: 'category',
                                data: [data[0].Key, data[1].Key, data[2].Key, data[3].Key, data[4].Key]
                                //data: ['祖庙街道办事处', '石湾街道办事处', '张槎街道办事处', '魁奇路', '南庄镇']
                            }
                        ],
                        yAxis: [
                            {
                                type: 'value'
                            }
                        ],
                        series: [
                            {
                                name: '禅城户籍人口',
                                type: 'bar',
                                itemStyle: {
                                    normal: {
                                        label: {
                                            show: true,
                                            position: 'top',
                                            formatter: '{c}'
                                        }
                                    }
                                },
                                data: [data[0].Value["s_本市户籍受理数"], data[1].Value["s_本市户籍受理数"], data[2].Value["s_本市户籍受理数"], data[3].Value["s_本市户籍受理数"], data[4].Value["s_本市户籍受理数"]]
                                //data: [75642, 77791, 158854, 293885, 40264]
                            },
                            {
                                name: '流动人口',
                                type: 'bar',
                                itemStyle: {
                                    normal: {
                                        label: {
                                            show: true,
                                            position: 'top',
                                            formatter: '{c}'
                                        }
                                    }
                                },
                                data: [data[0].Value["s_非本市户籍受理数"], data[1].Value["s_非本市户籍受理数"], data[2].Value["s_非本市户籍受理数"], data[3].Value["s_非本市户籍受理数"], data[4].Value["s_非本市户籍受理数"]]
                                //data: [44569, 78518, 59994, 73261, 21083]
                            }
                        ]
                    };
                    myChart.hideLoading();
                    myChart.setOption(option);
                }
                );
            }
            var initChart4 = function () {
                var myChart = echarts.init(document.getElementById('chart_4'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/ccbzzlzzzbzl"
                }).done(function (data) {

                    var itemLengths = [];
                    var itemLengths1 = [];

                    var ss = data["全区办件总量"];
                    var ss1 = data["魁奇路"];
                    var ss2 = data["南庄"];
                    var ss3 = data["石湾"];
                    var ss4 = data["祖庙"];
                    var ss5 = data["张槎"];

                    var temp = 0;

                    itemLengths.push(ss);
                    itemLengths1.push(0);
                    itemLengths.push(ss1);
                    itemLengths1.push(ss5 + ss4 + ss3 + ss2);
                    temp += ss1;
                    itemLengths.push(ss2);
                    itemLengths1.push(ss5 + ss4 + ss3);
                    temp += ss2;
                    itemLengths.push(ss3);
                    itemLengths1.push(ss5 + ss4);
                    temp += ss3;
                    itemLengths.push(ss4);
                    itemLengths1.push(ss5);
                    temp += ss4;
                    itemLengths.push(ss5);
                    itemLengths1.push(0);

                    var option = {
                        title: {
                            show: false,
                            text: '深圳月最低生活费组成（单位:元）',
                            subtext: 'From ExcelHome',
                            sublink: 'http://e.weibo.com/1341556070/AjQH99che'
                        },
                        tooltip: {
                            trigger: 'axis',
                            axisPointer: {            // 坐标轴指示器，坐标轴触发有效
                                type: 'shadow'        // 默认为直线，可选为：'line' | 'shadow'
                            },
                            formatter: function (params) {
                                var tar = params[0];
                                return tar.name + '<br/>' + tar.seriesName + ' : ' + tar.value;
                            }
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                mark: { show: false },
                                dataView: { show: true, readOnly: false },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        xAxis: [
                            {  
                                type: 'category',
                                splitLine: { show: false },
                                data: ["全区办件总量", "魁奇路", "南庄", "石湾", "祖庙", "张槎"]
                            }
                        ],
                        yAxis: [
                            {
                                type: 'value'
                            }
                        ],
                        series: [
                            {
                                name: '辅助',
                                type: 'bar',
                                stack: '总量',
                                itemStyle: {
                                    normal: {
                                        barBorderColor: 'rgba(0,0,0,0)',
                                        color: 'rgba(0,0,0,0)'
                                    },
                                    emphasis: {
                                        barBorderColor: 'rgba(0,0,0,0)',
                                        color: 'rgba(0,0,0,0)'
                                    }
                                },
                                data: itemLengths1
                                //data: [0, 913016, 766374, 540204, 160751, 0]
                            },
                            {
                                name: '办事量',
                                type: 'bar',
                                stack: '总量',
                                itemStyle: { normal: { label: { show: true, position: 'inside' } } },
                                data: itemLengths
                                //data: [1023037, 110221, 146642, 226170, 379453, 160751]
                            }
                        ]
                    };
                    myChart.hideLoading();
                    myChart.setOption(option);
                }
                );
            }
            var initChart5 = function () {
                var myChart = echarts.init(document.getElementById('chart_5'));
                myChart.showLoading({
                    text: "加载中.......",
                    effect: 'whirling',
                    textStyle: {
                        fontSize: 20
                    }
                });
                $.ajax({
                    url: url + "/total/hzsrrrfb"
                }).done(function (data) {
                    var ss1 = data["南庄"];
                    var ss2 = data["石湾"];
                    var ss3 = data["祖庙"];
                    var ss4 = data["张槎"];

                    var option = {
                        title: {
                            show: false,
                            text: '某站点用户访问来源',
                            subtext: '纯属虚构',
                            x: 'center'
                        },
                        tooltip: {
                            trigger: 'item',
                            formatter: "{a} <br/>{b} : {c} ({d}%)"
                        },
                        legend: {
                            show: false,
                            orient: 'vertical',
                            x: 'left',
                            data: ["南庄", "石湾", "祖庙", "张槎"]
                        },
                        toolbox: {
                            show: true,
                            feature: {
                                mark: { show: false },
                                dataView: { show: true, readOnly: false },
                                magicType: {
                                    show: false,
                                    type: ['pie', 'funnel'],
                                    option: {
                                        funnel: {
                                            x: '25%',
                                            width: '50%',
                                            funnelAlign: 'left',
                                            max: 1548
                                        }
                                    }
                                },
                                restore: { show: true },
                                saveAsImage: { show: true }
                            }
                        },
                        calculable: true,
                        series: [
                            {
                                name: '户籍失业人员分布',
                                type: 'pie',
                                radius: '60%',
                                center: ['50%', '50%'],
                                itemStyle: {
                                    normal: {
                                        label: {
                                            formatter: function (params) {
                                                return params.name + '\n' + params.percent + '%\n'
                                            },
                                            textStyle: {
                                                baseline: 'middle'
                                            }
                                        }
                                    }
                                },
                                data: [
                                    { value: ss1, name: "南庄" },
                                    { value: ss2, name: "石湾" },
                                    { value: ss3, name: "祖庙" },
                                    { value: ss4, name: "张槎" }
                                ]
                            }
                        ]
                    };
                    myChart.hideLoading();
                    myChart.setOption(option);
                }
                );
            }
            initChartAge();
            initChartServiceFrequency();
            initChart1();
            initChart2();
            initChart4();
            initChart5();
        }
        return IndexVM;
    })(BaseVM);
    viewModels.IndexVM = IndexVM;


    var firstpageVM = (function (_super) {
        __extends(firstpageVM, _super);
        var self = this;
        function firstpageVM() {
            _super.call(self);
        }

        self.gotott = function () {
            gotoview("pages/total", "", "");
        }

        self.gotots = function () {
            gotoview("pages/more", "", "");
        }

        self.gotosr = function () {
            gotoview("pages/detailes", "/sr/1/25", "");
        }

        self.gotolog = function () {
            gotoview("pages/logs/1", "/1/25", "");
        }

        self.gotomd = function () {
            gotoview("pages/dm", "", "");
        }

        firstpageVM.prototype.onNavigatedTo = function (params, payload) {
        };
        firstpageVM.prototype.canNavigateFrom = function (callback) {
            callback(true);
        };
        return firstpageVM;
    })(BaseVM);
    viewModels.firstpageVM = firstpageVM;


    var dmVM = (function (_super) {
        __extends(dmVM, _super);
        var self = this;
        function dmVM() {
            _super.call(self);
        }

        self.goto_cl_datamodel = function () {
            window.open("http://19.133.60.104:5656/");
        }

        self.goto_rs_datamodel = function () {
            window.open("http://19.133.60.104:8010/");
        }

        dmVM.prototype.onNavigatedTo = function (params, payload) {
        };
        dmVM.prototype.canNavigateFrom = function (callback) {
            callback(true);
        };
        return dmVM;
    })(BaseVM);
    viewModels.dmVM = dmVM;

    var plVM = (function (_super) {
        __extends(plVM, _super);
        var self = this;
        function plVM() {
            _super.call(self);
            self.nid = ko.observable();
            self.p_姓名 = ko.observable();
            self.p_性别 = ko.observable();
            self.p_生日 = ko.observable();
            self.p_固定电话 = ko.observable();
            self.p_移动电话 = ko.observable();
            self.p_现居住地址 = ko.observable();
            self.p_户籍所在地 = ko.observable();

            self.lgStatus = ko.observable();
            self.clStatus = ko.observable();
            self.jsStatus = ko.observable();
            self.rsStatus = ko.observable();
            self.ssStatus = ko.observable();
        }

        self.clickme = function (pr) {
            if (self.nid() !== undefined) {
                if (self.nid().length === 18) {
                    gotoview("pages/pf/", self.nid(), { k: pr });
                } else if (self.nid().length === 24) {
                    gotoview("pages/pf/", self.nid() + "/oid", { k: pr });
                } else {
                    myalert("号码意外！");
                }
            } else {
                myalert("请输入身份证号码");
            }
        }

        function loaddata(u) {
            $.ajax({
                url: u
            }).done(function (data) {
                self.p_姓名(data.p_姓名);
                self.p_性别(data.p_性别);
                self.p_生日((data.p_生日 === undefined || data.p_生日 === null) ? "" : data.p_生日.substring(0, 10));
                self.p_固定电话(data.p_固定电话 === "BsonNull" ? "" : data.p_固定电话);
                self.p_移动电话(data.p_移动电话);
                self.p_现居住地址(data.p_现居住地址);
                self.p_户籍所在地(data.p_户籍所在地);

                if (data.p_是否流管 === true) {
                    self.lgStatus('btn red');
                }
                if (data.p_是否残疾 === true) {
                    self.clStatus('btn green');
                }
                if (data.p_是否计生 === true) {
                    self.jsStatus('btn yellow');
                }
                if (data.p_是否人社 === true) {
                    self.rsStatus('btn purple');
                }
                if (data.p_是否失信 === true) {
                    self.ssStatus('btn blue-hoki');
                }
            });
        }

        plVM.prototype.onNavigatedTo = function (params, payload) {
            if (params.value === "oid") {
                self.nid(params.uid);
                var myurl = url + "/api/persons/oid/" + self.nid();
                loaddata(myurl);
            } else {
                self.nid(params.uid);
                var myurl = url + "/api/persons/uidcard/" + self.nid();
                loaddata(myurl);
            }

            var initview = function () {
                var myChart = echarts.init(document.getElementById('chart_pl'));
                var option = {
                    title: {
                        show: false,
                        text: '预算 vs 开销（Budget vs spending）',
                        subtext: '纯属虚构'
                    },
                    tooltip: {
                        trigger: 'axis'
                    },
                    legend: {
                        show: false,
                        orient: 'vertical',
                        x: 'right',
                        y: 'bottom',
                        data: ['预算分配', '实际开销']
                    },
                    toolbox: {
                        show: false,
                        feature: {
                            mark: { show: true },
                            dataView: { show: true, readOnly: false },
                            restore: { show: true },
                            saveAsImage: { show: true }
                        }
                    },
                    polar: [
                       {
                           indicator: [
                               { text: '个人信息', max: 30 },
                               { text: '流管', max: 10 },
                               { text: '残联', max: 10 },
                               { text: '信用', max: 10 },
                               { text: '人社', max: 10 },
                               { text: '计生', max: 10 }
                           ]
                       }
                    ],
                    calculable: true,
                    series: [
                        {
                            name: '预算 vs 开销（Budget vs spending）',
                            type: 'radar',
                            data: [
                                {
                                    value: [30, 4, 3, 1, 8, 4],
                                    name: '实际标签数'
                                }
                            ]
                        }
                    ]
                };
                myChart.setOption(option);
            }
            initview();
        };
        return plVM;
    })(BaseVM);
    viewModels.plVM = plVM;

    var pfVM = (function (_super) {
        __extends(pfVM, _super);
        var self = this;
        function pfVM() {
            _super.call(self);
            self.imgurl = ko.observable("http://19.133.103.22:8010/materialTake/");
            self.p_姓名 = ko.observable();
            self.p_性别 = ko.observable();
            self.p_生日 = ko.observable();
            self.p_身高 = ko.observable();
            self.p_视力 = ko.observable();
            self.p_固定电话 = ko.observable();
            self.p_移动电话 = ko.observable();
            self.p_户籍所在地 = ko.observable();
            self.p_现居住地址 = ko.observable();
            self.p_身份证号 = ko.observable();
            self.p_护照号 = ko.observable();
            self.p_居住证号 = ko.observable();
            self.p_军官证号 = ko.observable();
            self.p_出生证号 = ko.observable();
            self.p_绿卡号 = ko.observable();
            self.p_港澳通行证号 = ko.observable();
            self.p_社保卡号 = ko.observable();
            self.p_学籍号 = ko.observable();
            self.p_其他证件号 = ko.observable();
            self.p_就业状态 = ko.observable();
            self.p_工种名称 = ko.observable();
            self.p_期望工资 = ko.observable();
            self.p_工作地域要求 = ko.observable();
            self.p_残疾证号 = ko.observable();
            self.p_残疾证发证日期 = ko.observable();
            self.p_低保状态 = ko.observable();
            self.c_残疾类型 = ko.observable();
            self.c_残疾等级 = ko.observable();
            self.p_最小子女出生日期 = ko.observable();
            self.p_子女数 = ko.observable();
            self.p_女孩数 = ko.observable();
            self.p_当前有效避孕节育措施 = ko.observable();
            self.p_落实节育措施时间 = ko.observable();
            self.p_所属镇街 = ko.observable();
            self.p_所属村居 = ko.observable();
            self.p_婚姻状况 = ko.observable();
            self.p_户口性质 = ko.observable();
            self.p_政治面貌 = ko.observable();
            self.p_配偶姓名 = ko.observable();
            self.p_户籍归属 = ko.observable();
            self.p_初婚日期 = ko.observable();
            self.p_婚姻变动日期 = ko.observable();
            self.p_配偶证件号码 = ko.observable();
            self.p_户籍口径统计地 = ko.observable();
            self.p_常住口径统计地 = ko.observable();
            self.p_流动口径统计地 = ko.observable();
            self.p_民族 = ko.observable();
            self.p_职业 = ko.observable();
            self.p_本人是否独生 = ko.observable();
            self.p_工作单位 = ko.observable();
            self.p_单位性质 = ko.observable();
            self.p_健康状况 = ko.observable();
            self.p_文化程度 = ko.observable();
            self.p_迁出 = ko.observable();
            self.p_离开 = ko.observable();
            self.p_身高 = ko.observable();
            self.p_视力 = ko.observable();

            self.dvs = ko.observableArray();
            self.cls = ko.observableArray();
            self.pxs = ko.observableArray();
            self.jzs = ko.observableArray();
            self.lgs = ko.observableArray();
            self.srs = ko.observableArray();
            self.nid = ko.observable();
            self.clcount = ko.observable();
            self.bscount = ko.observable();
            self.cjzzcount = ko.observable();
            self.cjpscount = ko.observable();
            self.lgcount = ko.observable();
            self.sscount = ko.observable();

            self.policetext = ko.observable();
            self.policetitle = ko.observable();
        }

        self.serachpolice = function () {
            if (self.p_身份证号() !== undefined) {
                Metronic.blockUI({
                    boxed: true
                });
                $.ajax({
                    url: url + "/api/persons/uidcard/police/" + self.p_身份证号()
                }).done(function (data) {
                    self.policetitle(self.p_姓名() + "，" + self.p_性别());
                    self.policetext(JSON.stringify(data));
                    Metronic.unblockUI();
                    $('#m_police').modal('toggle');
                });
            }
        }

        function loaddata(u) {
            $.ajax({
                url: u
            }).done(function (data) {
                self.p_姓名(data.p_姓名);
                self.p_性别(data.p_性别);
                self.p_生日(data.p_生日 === null ? "" : data.p_生日.substring(0, 10));
                self.p_固定电话(data.p_固定电话 === "BsonNull" ? "" : data.p_固定电话);
                self.p_移动电话(data.p_移动电话);
                self.p_户籍所在地(data.p_户籍所在地);
                self.p_现居住地址(data.p_现居住地址);
                self.p_身份证号(data.p_身份证号);
                self.p_护照号(data.p_护照号);
                self.p_居住证号(data.p_居住证号);
                self.p_军官证号(data.p_军官证号);
                self.p_出生证号(data.p_出生证号);
                self.p_绿卡号(data.p_绿卡号);
                self.p_港澳通行证号(data.p_港澳通行证号);
                self.p_社保卡号(data.p_社保卡号);
                self.p_学籍号(data.p_学籍号);
                self.p_其他证件号(data.p_其他证件号);
                self.p_所属镇街(data.p_所属镇街);
                self.p_所属村居(data.p_所属村居);
                self.p_婚姻状况(data.p_婚姻状况);
                self.p_户口性质(data.p_户口性质);
                self.p_政治面貌(data.p_政治面貌);
                self.p_配偶姓名(data.p_配偶姓名);
                self.p_户籍归属(data.p_户籍归属);
                self.p_初婚日期(data.p_初婚日期);
                self.p_婚姻变动日期(data.p_婚姻变动日期);
                self.p_配偶证件号码(data.p_配偶证件号码);
                self.p_户籍口径统计地(data.p_户籍口径统计地);
                self.p_常住口径统计地(data.p_常住口径统计地);
                self.p_流动口径统计地(data.p_流动口径统计地);
                self.p_民族(data.p_民族);
                self.p_职业(data.p_职业);
                self.p_本人是否独生(data.p_本人是否独生);
                self.p_工作单位(data.p_工作单位);
                self.p_单位性质(data.p_单位性质);
                self.p_健康状况(data.p_健康状况);
                self.p_文化程度(data.p_文化程度);
                self.p_迁出(data.p_迁出);
                self.p_离开(data.p_离开);
                self.p_身高(data.p_身高);
                self.p_视力(data.p_视力);

                if (data.p_是否残疾) {
                    self.p_残疾证号(data.p_残疾证号);
                    if (data.p_残疾证发证日期 !== null && data.p_残疾证发证日期 !== undefined) {
                        self.p_残疾证发证日期(data.p_残疾证发证日期.split("T")[0]);
                    }
                    self.c_残疾类型(data.FK_残疾情况[0].c_残疾类型);
                    self.c_残疾等级(data.FK_残疾情况[0].c_残疾等级);

                    if (data.FK_残疾救助情况 !== null) {
                        ko.utils.arrayForEach(data.FK_残疾救助情况, function (dv) {
                            self.jzs.push(dv);
                        });
                        self.cjzzcount(data.FK_残疾救助情况.length);
                    }

                    if (data.FK_残疾培训情况 !== null) {
                        ko.utils.arrayForEach(data.FK_残疾培训情况, function (dv) {
                            dv.c_开始时间 = dv.c_开始时间.split("T")[0];
                            dv.c_结束时间 = dv.c_结束时间.split("T")[0];
                            self.pxs.push(dv);
                        });
                        self.cjpscount(data.FK_残疾培训情况.length);
                    }
                }

                if (data.p_是否计生) {
                    if (data.p_最小子女出生日期 !== null) {
                        self.p_最小子女出生日期(data.p_最小子女出生日期.split("T")[0]);
                    }

                    if (data.p_子女数) {
                        self.p_子女数(data.p_子女数);
                    }

                    if (data.p_女孩数) {
                        self.p_女孩数(data.p_女孩数);
                    }

                    if (data.p_当前有效避孕节育措施 !== null) {
                        self.p_当前有效避孕节育措施(data.p_当前有效避孕节育措施);
                    }

                    if (data.p_落实节育措施时间 !== null)
                        self.p_落实节育措施时间(data.p_落实节育措施时间.split("T")[0]);
                }

                if (data.p_是否人社) {
                    self.p_就业状态(data.p_就业状态);
                    self.p_工种名称(data.p_工种名称);
                    self.p_期望工资(data.p_期望工资);
                    self.p_工作地域要求(data.p_工作地域要求);
                }

                if (data.FK_办事经历 !== null) {
                    ko.utils.arrayForEach(data.FK_办事经历, function (dv) {
                        dv.a_受理时间 = dv.a_受理时间.split("T")[0] + " " + dv.a_受理时间.split("T")[1].split("Z")[0];
                        dv.a_办结时间 = dv.a_办结时间.split("T")[0] + " " + dv.a_办结时间.split("T")[1].split("Z")[0];
                        self.dvs.push(dv);
                    });
                    self.bscount(data.FK_办事经历.length);
                }

                if (data.FK_流管经历 !== null) {
                    ko.utils.arrayForEach(data.FK_流管经历, function (dv) {
                        self.lgs.push(dv);
                    });
                    self.lgcount(data.FK_流管经历.length);
                }

                if (data.FK_办事经历 !== null) {
                    ko.utils.arrayForEach(data.FK_材料记录, function (dv) {
                        if (dv.c_文件路径 !== "") {
                            self.cls.push(dv);
                        }
                    });
                    self.clcount(data.FK_材料记录.length);
                }

                if (data.FK_失信记录 !== null) {
                    ko.utils.arrayForEach(data.FK_失信记录, function (dv) {
                        self.srs.push(dv);
                    });
                    self.sscount(data.FK_失信记录.length);
                }
            });
        }

        pfVM.prototype.onNavigatedTo = function (params, payload) {
            if (params.value === "oid") {
                loaddata(url + "/api/persons/oid/" + params.uid);
            } else {
                loaddata(url + "/api/persons/uidcard/" + params.uid);
            }
            if (payload !== null && payload !== undefined && payload.k !== null) {
                if (payload.k !== '') {
                    if (payload.k === 'lgview') {
                        $('html, body').animate({ scrollTop: 2000 }, 500);
                    } else if (payload.k === 'clview') {
                        $('html, body').animate({ scrollTop: 1550 }, 500);
                    } else if (payload.k === 'zsview') {
                        $('html, body').animate({ scrollTop: 2300 }, 500);
                    } else if (payload.k === 'rsview') {
                        $('html, body').animate({ scrollTop: 1100 }, 500);
                    } else if (payload.k === 'ssview') {
                        $('html, body').animate({ scrollTop: 1100 }, 500);
                    }
                }
            }
        };
        pfVM.prototype.canNavigateFrom = function (callback) {
            callback(true);
        };
        return pfVM;
    })(BaseVM);
    viewModels.pfVM = pfVM;

    var moreVM = (function (_super) {
        __extends(moreVM, _super);
        var self = this;
        function moreVM() {
            _super.call(self);
            self.lgcount = ko.observable();
            self.jscount = ko.observable();
            self.clcount = ko.observable();
        }

        self.gotozdrc = function () {
            window.open('http://xy.fspc.gov.cn/focusgroups/index.html?leftnum=2&num=4');
        }

        self.gotocl = function () {
            gotoview("pages/detailes", "/cl/1/25", "");
        }

        self.gotojs = function () {
            gotoview("pages/detailes", "/js/1/25", "");
        }

        self.gotolg = function () {
            gotoview("pages/detailes", "/lg/1/25", "");
        }

        moreVM.prototype.onNavigatedTo = function (params, payload) {
            $.ajax({
                url: url + "/total/gszrsl"
            }).done(function (data) {
                self.clcount(data[0].Value);
                self.lgcount(data[1].Value);
                self.jscount(data[2].Value);
            });
        };
        return moreVM;
    })(BaseVM);
    viewModels.moreVM = moreVM;

    var detailesVM = (function (_super) {
        __extends(detailesVM, _super);
        var self = this;
        function detailesVM() {
            _super.call(self);
            self.vpps = ko.observableArray();
            self.index = ko.observable(1);
            self.pagesize = ko.observable(25);
            self.skin = ko.observable();
            self.currentlabel = ko.observable();
            self.totalpage = ko.observable(0);
        }

        self.viewdetail = function (dv) {
            gotoview("pages/pl/", dv.Id + "/oid", new Array());
        }

        self.up = function () {
            if (self.index() > 1) {
                self.index(self.index() - 1);
                myRefresh("/api/persons/list/" + self.skin());
                $('html, body').animate({ scrollTop: 0 }, 500);
            }
        }
        self.down = function () {
            self.index(self.index() + 1);
            var d = self.skin();
            myRefresh("/api/persons/list/" + self.skin());
            $('html, body').animate({ scrollTop: 0 }, 500);
        }

        function myRefresh(path) {
            var ddd = url + path + "/" + self.pagesize() + "/" + self.index();
            $.ajax({
                url: url + path + "/" + self.pagesize() + "/" + self.index()
            }).done(function (data) {
                if (data === "0") {
                    myalert("没有更多数据了！");
                } else {
                    self.vpps.removeAll();
                    ko.utils.arrayForEach(data, function (ss) {
                        if (ss.p_生日 === null) {
                            ss.p_生日 = "";
                        }
                        self.vpps.push(ss);
                    });
                }
            });
        }

        function myRefreshCount(path) {
            $.ajax({
                url: url + path + "/pagect/" + self.pagesize()
            }).done(function (data) {
                self.totalpage(data);
            });
        }

        detailesVM.prototype.onNavigatedTo = function (params, payload) {
            self.pagesize(parseInt(params.ps) !== NaN ? parseInt(params.ps) : 25);
            self.index(parseInt(params.index) !== NaN ? parseInt(params.index) : 1);
            self.skin(params.skin);
            if (self.skin() === 'cl') {
                self.currentlabel('残疾人员数据')
                myRefresh("/api/persons/list/" + self.skin());
                myRefreshCount("/api/persons/list/" + self.skin());
            } else if (self.skin() === "js") {
                self.currentlabel('常住人员数据')
                myRefresh("/api/persons/list/" + self.skin());
                myRefreshCount("/api/persons/list/" + self.skin());
            } else if (self.skin() === "lg") {
                self.currentlabel('流动人员数据')
                myRefresh("/api/persons/list/" + self.skin());
                myRefreshCount("/api/persons/list/" + self.skin());
            } else if (self.skin() === "sr") {
                self.currentlabel('失信人员数据')
                myRefresh("/api/persons/list/" + self.skin());
                myRefreshCount("/api/persons/list/" + self.skin());
            } else {
                myalert("非法请求页面！");
            }
        };
        return detailesVM;
    })(BaseVM);
    viewModels.detailesVM = detailesVM;

    var logsVM = (function (_super) {
        __extends(logsVM, _super);
        var self = this;
        function logsVM() {
            _super.call(self);
            self.vpps = ko.observableArray();
            self.index = ko.observable(1);
            self.pagesize = ko.observable(25);
            self.currentlabel = ko.observable();
            self.totalpage = ko.observable(0);
            self.skin = ko.observable(1);
        }

        self.logsetskin = function (sk) {
            self.skin(sk);
            myRefresh("/api/persons/list/log/" + self.skin());
            myRefreshCount("/api/persons/list/log/" + self.skin());
        }

        self.logup = function () {
            if (self.index() > 1) {
                self.index(self.index() - 1);
                myRefresh("/api/persons/list/log/" + self.skin());
                $('html, body').animate({ scrollTop: 0 }, 500);
            }
        }
        self.logdown = function () {
            self.index(self.index() + 1);
            myRefresh("/api/persons/list/log/" + self.skin());
            $('html, body').animate({ scrollTop: 0 }, 500);
        }

        function myRefresh(path) {
            var d = url + path + "/" + self.pagesize() + "/" + self.index();
            $.ajax({
                url: url + path + "/" + self.pagesize() + "/" + self.index()
            }).done(function (data) {
                if (data === "0") {
                    myalert("没有更多数据了！");
                } else if (data === "1") {
                    myalert("数据库数据源出错，请清空日志表");
                } else {
                    self.vpps.removeAll();
                    ko.utils.arrayForEach(data, function (ss) {
                        self.vpps.push(ss);
                    });
                }
            });
        }

        function myRefreshCount(path) {
            $.ajax({
                url: url + path + "/pagect/" + self.pagesize()
            }).done(function (data) {
                self.totalpage(data);
            });
        }

        logsVM.prototype.onNavigatedTo = function (params, payload) {
            self.skin(parseInt(params.skin) !== NaN ? parseInt(params.skin) : 1);
            self.pagesize(parseInt(params.ps) !== NaN ? parseInt(params.ps) : 25);
            self.index(parseInt(params.index) !== NaN ? parseInt(params.index) : 1);
            self.currentlabel('调用日志数据')
            myRefresh("/api/persons/list/log/" + self.skin());
            myRefreshCount("/api/persons/list/log/" + self.skin());
        };
        return logsVM;
    })(BaseVM);
    viewModels.logsVM = logsVM;

})(viewModels || (viewModels = {}));