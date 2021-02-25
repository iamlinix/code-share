<template>
    <div class="price-chart-block">
        <el-card style="margin-bottom: 20px">
            <div class="chart-title">
                <span >{{ title }}</span>
                <el-tag size="mini" effect="dark">{{ `总销量：${beauty(countTotal)}` }}</el-tag>
                <el-tag size="mini" effect="dark">{{ `总销售额：${beauty(saleTotal.toFixed(2))}` }}</el-tag>
                <el-tag size="mini" effect="dark">{{ `平均售价：${beauty((saleTotal / countTotal).toFixed(2))}` }}</el-tag>
            </div>
            <p v-if="!rawData" style="width: 100%; text-align: center; font-size: 12px; color: #909399">无可用数据</p>
            <el-row v-else :gutter="8">
                <el-col :span="countSpan">
                    <el-card :id="chartKey + 'ResizePivot'" style="height: 40vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                        <i v-show="!showSaleChart" class="el-icon-more" style="right: 8px; color: #909399; cursor: pointer; position: absolute" @click="toggleSaleChart(true)"></i>
                        <el-dropdown @command="changeChartType" size="mini" trigger="click" style="position: absolute; right: 30px; z-index: 10">
                            <span class="el-dropdown-link">
                                图表类型<i class="el-icon-arrow-down el-icon--right"></i>
                            </span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item icon="el-icon-pie-chart" command="cnt-doughnut">空心饼图</el-dropdown-item>
                                <el-dropdown-item icon="el-icon-pie-chart" command="cnt-pie">实心饼图</el-dropdown-item>
                                <el-dropdown-item icon="el-icon-data-analysis" command="cnt-bar">柱状图</el-dropdown-item>
                                <el-dropdown-item icon="el-icon-data-line" command="cnt-line">折线图</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                        <div style="width: 100%; height: 100%;" :id="chartKey + '-count-chart'" />
                    </el-card>
                </el-col>
                <el-col :span="saleSpan">
                    <el-card style="height: 40vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                        <div v-show="showSaleChart" style="width: 100%; height: 100%">
                            <i class="el-icon-close" style="right: 8px; color: #E6A23C; cursor: pointer; position: absolute" @click="toggleSaleChart(false)"></i>
                            <el-dropdown @command="changeChartType" size="mini" trigger="click" style="position: absolute; right: 30px; z-index: 10">
                                <span class="el-dropdown-link">
                                    图表类型<i class="el-icon-arrow-down el-icon--right"></i>
                                </span>
                                <el-dropdown-menu slot="dropdown">
                                    <el-dropdown-item icon="el-icon-pie-chart" command="sale-doughnut">空心饼图</el-dropdown-item>
                                    <el-dropdown-item icon="el-icon-pie-chart" command="sale-pie">实心饼图</el-dropdown-item>
                                    <el-dropdown-item icon="el-icon-data-analysis" command="sale-bar">柱状图</el-dropdown-item>
                                    <el-dropdown-item icon="el-icon-data-line" command="sale-line">折线图</el-dropdown-item>
                                </el-dropdown-menu>
                            </el-dropdown>
                            <div style="width: 100%; height: 100%;" :id="chartKey + '-sale-chart'" />
                        </div>
                    </el-card>
                </el-col>
                <el-col :span="8">
                    <el-table :key="tableKey" :ref="chartKey" :id="chartKey + 'CountTable'" :stripe="true" max-height="40vh" :data="rawData" :border="true" size="small">
                        <!--el-table-column type="expand">
                            <template slot-scope="scope">
                                <el-table class="price-chart-expand" :row-class-name="expandRowClass" :key="expandKey" size="mini" :stripe="true" :data="rawData[scope.$index].expand" :border="false">
                                    <el-table-column label="排名" width="45px">
                                        <template slot-scope="sscope">
                                            <span>{{ sscope.$index + 1 }}</span>
                                        </template>
                                    </el-table-column>
                                    <el-table-column label="商品名称" prop="materialTxt" />
                                    <el-table-column label="销售数量" prop="salesCount" width="70px"/>
                                    <el-table-column label="占比" width="70px">
                                        <template slot-scope="sscope">
                                            <span>{{ (sscope.row.salesCount * 100 / scope.row.salesCount).toFixed(2) + '%' }} </span>
                                        </template>
                                    </el-table-column>
                                </el-table>
                            </template>
                        </el-table-column!-->
                        <el-table-column label="价格区间" prop="priceZone" />
                        <el-table-column label="销售数量" prop="salesCount" />
                        <el-table-column label="占比">
                            <template slot-scope="scope">
                                <span>{{ (scope.row.salesCount * 100 / countTotal).toFixed(2) + '%' }} </span>
                            </template>
                        </el-table-column>
                        <el-table-column label="Top10" width="60px">
                            <template slot-scope="props">
                                <el-popover trigger="click" placement="left" width=500 propper-class="chart-block-proper">
                                    <p style="font-size: 12px; width: 100%; padding: 4px 8px 4px 0px; margin-bottom: 8px">
                                        {{ '区间：' + props.row.priceZone }}
                                    </p>
                                    <div v-if="!rawData[props.$index].expand" class="gj-center-div" :key="expandKeyLoading">
                                        <i class="el-icon-loading" />
                                    </div>
                                    <el-table v-else max-height="300px" class="price-chart-expand" :row-class-name="expandRowClass" :key="expandKey" size="mini" :stripe="true" :data="rawData[props.$index].expand" :border="false">
                                        <el-table-column label="排名" width="45px">
                                            <template slot-scope="sscope">
                                                <span>{{ sscope.$index + 1 }}</span>
                                            </template>
                                        </el-table-column>
                                        <el-table-column label="商品名称" prop="materialTxt" />
                                        <el-table-column label="销售数量" prop="salesCount" width="70px"/>
                                        <el-table-column label="占比" width="70px">
                                            <template slot-scope="sscope">
                                                <span>{{ (sscope.row.salesCount * 100 / props.row.salesCount).toFixed(2) + '%' }} </span>
                                            </template>
                                        </el-table-column>
                                    </el-table>
                                    <div slot="reference" style="width: 100%; text-align: center; cursor: pointer" @click="expandChange(props.row)">
                                        <i class="el-icon-trophy"/>
                                    </div>
                                </el-popover>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-col>
            </el-row>
        </el-card>
    </div>
</template>

<style lang="scss">
    .price-chart-block {
        .chart-title {
            padding: 0px 0px 12px 0px;
            font-size: 13px;

            span {
                margin-right: 12px;
            }
        }

        .el-dropdown-link {
            cursor: pointer;
            color: #409EFF;
            font-size: 10px;
            opacity: 0.2;
        }

        .el-dropdown-link:hover {
            opacity: 1;
        }

        .chart-table-max-height {
            max-height: 40vh
        }

        .el-table__body-wrapper {
            max-height: calc(40vh - 40px);
        }

        .el-table__body {
            max-height: calc(40vh - 40px);
        }

        .price-chart-expand {
            td {
                border: none;
            }

            th {
                border: none;
            }
        }

        .expand-odd-row {
            background-color: #C6E2FF
        }

        .expand-even-row {
            background-color: #ECF5FF
        }
    }
</style>

<script>
import echarts from 'echarts';
import 'echarts/theme/macarons'
import { doRequest, beautifyNumber } from '../../utils/utils';
import SmartResizeChart from '../mixins/SmartResizeChart'

export default {
    mixins: [ SmartResizeChart ],
    name: 'ChartBlock',
    props: {
        rawData: Array,
        title: String,
        chartKey: String,
        beginDate: String,
        endDate: String,
        orgLevel: Number,
        orgCode: String,
        classLevel: Number,
        classCode: String
    },
    data() {
        return {
            countChart: undefined,
            saleChart: undefined,
            countChartOption: {},
            saleChartOption: {},
            countLegend: [],
            countTotal: 0,
            saleTotal: 0,
            saleLegend: [],
            countData: [],
            saleData: [],
            expandKey: this.chartKey + 'ExpandKey',
            expandKeyLoading: this.chartKey + 'ExpandKeyLoading',
            tableKey: this.chartKey + 'Table',
            showSaleChart: true,
            countSpan: 8,
            saleSpan: 8,
            smartResizeCounter: 0,
            smartResizeWidth: 0,
            smartResizeHeight: 0,

            countChartType: 'doughnut',
            saleChartType: 'doughnut',
            doughnutOption: {
                title: {
                    text: '',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item',
                    formatter: '{a} <br/>{b} : {c} ({d}%)'
                },
                legend: {
                    orient: 'vertical',
                    left: 'left',
                    data: []
                },
                series: [
                    {
                        name: '',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        center: ['60%', '60%'],
                        data: [],
                        avoidLabelOverlap: false,
                        label: {
                            normal: {
                                show: false,
                                position: 'center'
                            },
                            emphasis: {
                                show: true,
                                textStyle: {
                                    fontSize: '20',
                                    fontWeight: 'bold'
                                }
                            }
                        },
                        labelLine: {
                            normal: {
                                show: false
                            }
                        }
                    }
                ]
            },
            pieOption: {
                title: {
                    text: '',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item',
                    formatter: '{a} <br/>{b} : {c} ({d}%)'
                },
                legend: {
                    orient: 'vertical',
                    left: 'left',
                    data: []
                },
                series: [
                    {
                        name: '',
                        type: 'pie',
                        radius: '70%',
                        center: ['60%', '60%'],
                        data: [],
                        label: {
                            normal: {
                                show: false,
                                position: 'center'
                            }
                        },
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)'
                            }
                        }
                    }
                ]
            },
            barOption: {
                title: {
                    text: '',
                    left: 'center'
                },
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                    right: 10,
                    top: 20,
                    bottom: 20,
                    data: []
                },
                tooltip: {
                    trigger: 'axis',
                    formatter: '{a} <br/>{b} : {c}',
                    axisPointer: {            // 坐标轴指示器，坐标轴触发有效
                        type: 'shadow'        // 默认为直线，可选为：'line' | 'shadow'
                    }
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                xAxis: {
                    type: 'category',
                    data: [],
                    axisTick: {
                        alignWithLabel: true
                    }
                },
                yAxis: {
                    type: 'value'
                },
                series: [{
                    data: [],
                    type: 'bar',
                    barWidth: '40%'
                }]
            },
            lineOption: {
                title: {
                    text: '',
                    left: 'center'
                },
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                    right: 10,
                    top: 20,
                    bottom: 20,
                    data: []
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    data: []
                },
                yAxis: {
                    type: 'value'
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        label: {
                            backgroundColor: '#6a7985'
                        }
                    }
                },
                series: [{
                    data: [],
                    type: 'line',
                    smooth: true,
                    areaStyle: {}
                }]
            }
        }
    },
    computed: {
        tableMaxHeight: function() {
            return parseInt(document.documentElement.clientHeight * 0.4) + 'px';
        }
    },
    mounted: function() {
        if (this.rawData) {
            for (let i = 0; i < this.rawData.length; i ++) {
                this.rawData[i].index = i;
                this.countTotal += this.rawData[i].salesCount;
                this.saleTotal += this.rawData[i].salesIncome;
            }

            this.countChartOptions = this.generateChartOption('cnt')
            this.saleChartOptions = this.generateChartOption('sale');

            this.countChart = echarts.init(document.getElementById(`${this.chartKey}-count-chart`), 'macarons');
            this.countChart.setOption(this.countChartOptions, true);

            this.saleChart = echarts.init(document.getElementById(`${this.chartKey}-sale-chart`), 'macarons');
            this.saleChart.setOption(this.saleChartOptions, true);

            this.smrschtAddChart(this.countChart);
            this.smrschtAddChart(this.saleChart);
            this.smrschtMonitorResizeById(`${this.chartKey}-count-chart`);

            // window.addEventListener('resize', this.delayedResize);
        }
    },
    updated: function() {
        this.$nextTick(function () {
            
        })
    },
    methods: {
        deepCopy(obj) {
            return JSON.parse(JSON.stringify(obj));
        },
        generateChartOption(which) {
            let legend = [], sortedLegend = [], data = [], type = '', option = {}, title = '', ele = null, i = 0, seriesName = '';
            if (this.rawData) {
                if (which === 'cnt') {
                    for (i = 0; i < this.rawData.length; i ++) {
                        ele = this.rawData[i];
                        legend.push(ele.priceZone);
                        sortedLegend.push(ele.priceZone);
                        data.push({
                            value: ele.salesCount,
                            name: ele.priceZone
                        });
                        this.saleData.push({
                            value: ele.salesIncome,
                            name: ele.priceZone
                        });
                    }
                    data.sort((a, b) => {
                        return -(a.value - b.value)
                    })
                    sortedLegend.sort((a, b) => {
                        if (a.startsWith('>='))
                            return 1;
                        if (b.startsWith('>='))
                            return -1;
                        let pa = a.split('-'), pb = b.split('-');
                        let fa = parseFloat(pa[0]), fb = parseFloat(pb[0]);
                        return fa - fb;
                    })
                    type = this.countChartType;
                    title = '销量分布';
                    seriesName = '销售数量';
                } else {
                    for (i = 0; i < this.rawData.length; i ++) {
                        ele = this.rawData[i];
                        legend.push(ele.priceZone);
                        sortedLegend .push(ele.priceZone);
                        data.push({
                            value: ele.salesIncome,
                            name: ele.priceZone
                        });
                    }
                    data.sort((a, b) => {
                        return -(a.value - b.value)
                    })
                    sortedLegend.sort((a, b) => {
                        if (a.startsWith('>='))
                            return 1;
                        if (b.startsWith('>='))
                            return -1;
                        let pa = a.split('-'), pb = b.split('-');
                        let fa = parseFloat(pa[0]), fb = parseFloat(pb[0]);
                        return fa - fb;
                    })
                    type = this.saleChartType;
                    title = '销售额分布';
                    seriesName = '销售额';
                }

                if (type === 'bar') {
                    option = this.deepCopy(this.barOption);
                    option.title.text = title;
                    option.xAxis.data = legend;
                    option.legend.data = sortedLegend;
                    option.series[0].data = data;
                    option.series[0].name = seriesName;
                } else if (type === 'doughnut') {
                    option = this.deepCopy(this.doughnutOption);
                    option.title.text = title;
                    option.legend.data = legend;
                    option.series[0].data = data;
                    option.series[0].name = seriesName;
                } else if (type === 'pie') {
                    option = this.deepCopy(this.pieOption);
                    option.title.text = title;
                    option.legend.data = legend;
                    option.series[0].data = data;
                    option.series[0].name = seriesName;
                } else if (type === 'line') {
                    option = this.deepCopy(this.lineOption);
                    option.title.text = title;
                    option.xAxis.data = legend;
                    option.legend.data = sortedLegend;
                    option.series[0].data = data;
                    option.series[0].name = seriesName;
                }
            }
            return option 
        },
        expandRowClass({row, rowIndex}) {
            if (rowIndex % 2 === 0)
                return 'expand-even-row';
            else
                return 'expand-odd-row';
        },
        toggleKey(key) {
            let suffix = '-toggle';
            let val = this[key];
            if (val) {
                if (val.endsWith(suffix))
                    this[key] = val.substring(0, val.length - suffix.length);
                else
                    this[key] = val + suffix;
            }
            // if (key.endsWith(suffix)) 
            //     return key.substring(0, key.length - suffix.length);
            // else
            //     return key + suffix;
        },
        doChartResize() {
            this.countChart.resize();
            this.saleChart.resize();
            // let ele = document.getElementById(`${this.chartKey}ResizePivot`);
            // if (ele) {
            //     document.getElementById(`${this.chartKey}CountTable`).style.maxHeight = ele.offsetHeight + 'px';
                // this.$refs[this.chartKey].doLayout();
                // this.expandKey = this.toggleKey(this.expandKey);
            // }
        },
        delayedResize() {
            let ele = document.getElementById(`${this.chartKey}-count-chart`);
            if (ele) {
                if (ele.clientWidth !== this.smartResizeWidth || ele.clientHeight !== this.smartResizeHeight) {
                    this.smartResizeWidth = ele.clientWidth;
                    this.smartResizeHeight = ele.clientHeight;
                    setTimeout(() => {
                        this.delayedResize();
                    }, 50)
                } else {
                    this.smartResizeCounter += 1;
                    if (this.smartResizeCounter >= 3) {
                        this.smartResizeWidth = 0;
                        this.smartResizeHeight = 0;
                        this.smartResizeCounter = 0;
                        this.doChartResize()
                    } else {
                        setTimeout(() => {
                            this.delayedResize();
                        }, 50)
                    }
                }
            }
            // console.log(document.getElementById(`${this.chartKey}-count-chart`).clientHeight);
            // setTimeout(() => {
            //     this.doChartResize();
            //     // setTimeout(() => {
            //     //     this.doChartResize();
            //     //     setTimeout(() => {
            //     //         this.doChartResize();
            //     //     }, 50)
            //     // }, 50)
            // }, 50)
        },
        expandChange(row) {
            if (!this.rawData[row.index].expand) {
                doRequest({
                    url: '/v1/web/pricezone/rank',
                    method: 'post',
                    data: {
                        beginDate: this.beginDate,
                        endDate: this.endDate,
                        orgLevel: this.orgLevel,
                        orgCode: this.orgCode,
                        classLevel: this.classLevel,
                        classCode: this.classCode,
                        zoneTxt: row.priceZone
                    }
                }, {
                    success: res => {
                        this.rawData[row.index].expand = res.mList;
                        this.toggleKey('expandKey');
                        this.toggleKey('expandKeyLoading');
                    },
                    fail: err => {
                        console.log(err);
                    }
                })
            }
        },
        toggleSaleChart(show) {
            if (show) {
                this.countSpan = 8;
                this.saleSpan = 8;
                this.showSaleChart = true;
            } else {
                this.countSpan = 16;
                this.saleSpan = 0;
                this.showSaleChart = false;
            }
            // this.delayedResize();
            this.smrschtDoSmartDelayedResize();
        },
        changeChartType(cmd) {
            let p = cmd.split('-');
            let chartType = p[1];
            if (p[0] === 'cnt') {
                if (chartType !== this.countChartType) {
                    this.countChartType = chartType;
                    this.countChartOption = this.generateChartOption('cnt');
                    this.countChart.setOption(this.countChartOption, true);
                }
            } else {
                if (chartType !== this.saleChartType) {
                    this.saleChartType = chartType;
                    this.saleChartOption = this.generateChartOption('sale');
                    this.saleChart.setOption(this.saleChartOption, true);
                }
            }
        },
        beauty(num) {
            return beautifyNumber(num);
        }
    }
}
</script>