<template>
    <div class="cnpc-focus-main">
        <transition name="zoom">
            <el-card v-if="showUpper" style="animation-duration: 0.2s" v-loading="upperLoading">
                <el-input size="mini" v-model="material" style="width: 200px"
                    placeholder="请输入需要关注的商品编码" clearable/>
                <el-date-picker size="mini" style="padding-top: 2px; margin-left: 8px;"
                    type="daterange" v-model="dateRange" range-separator="-"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="yyyy-MM-dd" />
                <el-button size="mini" type="primary" style="margin-left: 8px; margin-bottom: 8px"
                    @click="doFocus" icon="el-icon-view">
                    查看关注
                </el-button>

                <p class="recent-focus">
                    最近关注 ({{ focuses.length }})
                    <el-input size="mini" v-model="focusFilter"
                        style="width: 250px; margin-left: 8px; margin-right: 8px"
                        prefix-icon="el-icon-search"
                        placeholder="输入商品编码或名称查询历史关注"
                        clearable>
                    </el-input>
                    <el-button size="mini" type="danger" icon="el-icon-delete"
                        @click="clearFocus">清空历史</el-button>
                </p>
                <div class="top-scroll-bar">
                    <el-tooltip v-for="v in filteredFocuses" :key="'brick-' + v.material"  placement="top"
                        :hide-after="500" :content="v.material">
                        <div class="brick">
                            <span @click="material = v.material">
                                {{ v.materialTxt }}
                            </span>
                            <i class="ri-close-line close-button" @click="delFocus(v)"/>
                        </div>

                    </el-tooltip>
                    <div class="brick" v-if="filteredFocuses.hasMore">
                        <el-tooltip placement="top" :hide-after="2000" content="更多历史关注请过滤查找">
                            <p>...</p>
                        </el-tooltip>
                    </div>
                </div>
                <i class="ri-arrow-up-s-line collapse-button" @click="showUpper = false"></i>
            </el-card>
        </transition>
        <div style="margin-top: 8px; position: relative" >
            <i v-if="!showUpper" class="ri-arrow-down-s-line expand-button" @click="showUpper = true"></i>
            <el-tabs type="card" closable v-model="activeTab" @tab-remove="removeTab" @tab-click="clickTab">
                <el-tab-pane v-for="v in tabs" :key="'tab-pane-' + v.name" :name="v.name">
                    <el-tooltip slot="label" :hide-after="2000" placement="top"
                        :content="`${v.material} (${v.sdate} ~ ${v.edate})`">
                        <span>{{ v.label }}</span>
                    </el-tooltip>
                    <el-card class="lower-card" v-loading="v.loading">
                        <el-row :gutter="4">
                            <el-col :span="8">
                                <el-card>
                                    <el-row :gutter="4">
                                        <el-col :span="18">
                                            <div class="focus-main-number">
                                                <p class="title">销售收入</p>
                                                <p class="amount-number">
                                                    {{ v.net ? beautyNum(v.net.value.toFixed(2)) : '' }}
                                                </p>
                                            </div>
                                        </el-col>
                                        <el-col :span="6">
                                            <el-row :gutter="0">
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">同比</p>
                                                        <p v-if="v.net && v.net.yoy > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.net ? v.net.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.net ? v.net.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">环比</p>
                                                        <p v-if="v.net && v.net.mom > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.net ? v.net.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.net ? v.net.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                            </el-row>
                                        </el-col>
                                    </el-row>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card>
                                    <el-row :gutter="4">
                                        <el-col :span="18">
                                            <div class="focus-main-number">
                                                <p class="title">销售毛利</p>
                                                <p class="amount-number">
                                                    {{ v.profit ? beautyNum(v.profit.value.toFixed(2)) : '' }}
                                                </p>
                                            </div>
                                        </el-col>
                                        <el-col :span="6">
                                            <el-row :gutter="0">
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">同比</p>
                                                        <p v-if="v.profit && v.profit.yoy > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.profit ? v.profit.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.profit ? v.profit.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">环比</p>
                                                        <p v-if="v.profit && v.profit.mom > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.profit ? v.profit.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.profit ? v.profit.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                            </el-row>
                                        </el-col>
                                    </el-row>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card>
                                    <el-row :gutter="4">
                                        <el-col :span="18">
                                            <div class="focus-main-number">
                                                <p class="title">销售毛利率</p>
                                                <p class="rate-number">
                                                    {{ v.margin ? beautyNum(v.margin.value.toFixed(2)) : '-' }}%
                                                </p>
                                            </div>
                                        </el-col>
                                        <el-col :span="6">
                                            <el-row :gutter="0">
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">同比</p>
                                                        <p v-if="v.margin && v.margin.yoy > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.margin ? v.margin.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.margin ? v.margin.yoy.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                                <el-col :span="24">
                                                    <div class="focus-sub-number">
                                                        <p class="title">环比</p>
                                                        <p v-if="v.margin && v.margin.mom > 0" class="percent-number">
                                                            <font-awesome-icon icon="arrow-up" />
                                                            {{ v.margin ? v.margin.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                        <p v-else class="percent-number neg-number">
                                                            <font-awesome-icon icon="arrow-down" />
                                                            {{ v.margin ? v.margin.mom.toFixed(2) + '%' : '' }}
                                                        </p>
                                                    </div>
                                                </el-col>
                                            </el-row>
                                        </el-col>
                                    </el-row>
                                </el-card>
                            </el-col>
                            <el-col :span="24">
                                <el-card class="chart-card">
                                    <v-chart ref="focus-line-chart" theme="macarons" :options="v.lineOptions"
                                        :autoresize="true"/>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="margin-top: 4px">
                                    <p class="plant-list-title">销售排名 TOP10</p>
                                    <x-table :data="v.top" :minus="400" show-overflow="tooltip" size="mini">
                                        <vxe-table-column title="油站编码" width="80" field="OrgCode" />
                                        <vxe-table-column title="油站名称" field="OrgText" />
                                        <vxe-table-column title="销售收入" width="100" field="Metric.NetvalInv">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.NetvalInv.toFixed(2)) }}
                                            </template>
                                        </vxe-table-column>
                                        <vxe-table-column title="毛利" width="70" field="Metric.GrossProfit">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.GrossProfit.toFixed(2)) }}
                                            </template>
                                        </vxe-table-column>
                                        <vxe-table-column title="毛利率" width="70" field="Metric.GrossMargin">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.GrossMargin.toFixed(2)) }}%
                                            </template>
                                        </vxe-table-column>
                                    </x-table>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="margin-top: 4px">
                                    <p class="plant-list-title">销售排名 反向TOP10</p>
                                    <x-table :data="v.rop" :minus="400" show-overflow="tooltip" size="mini">
                                        <vxe-table-column title="油站编码" width="80" field="OrgCode" />
                                        <vxe-table-column title="油站名称" field="OrgText" />
                                        <vxe-table-column title="销售收入" width="100" field="Metric.NetvalInv">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.NetvalInv.toFixed(2)) }}
                                            </template>
                                        </vxe-table-column>
                                        <vxe-table-column title="毛利" width="70" field="Metric.GrossProfit">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.GrossProfit.toFixed(2)) }}
                                            </template>
                                        </vxe-table-column>
                                        <vxe-table-column title="毛利率" width="70" field="Metric.GrossMargin">
                                            <template slot-scope="r">
                                                {{ beautyNum(r.row.Metric.GrossMargin.toFixed(2)) }}%
                                            </template>
                                        </vxe-table-column>
                                    </x-table>
                                </el-card>
                            </el-col>
                        </el-row>
                    </el-card>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>

<style lang="scss">
.cnpc-focus-main {
    background-color: #FFF;
    height: calc(100vh - 10px);
    overflow: auto;
    padding: 8px 8px 8px 8px;
    width: 100%;
    position: relative;

    .top-scroll-bar {
        padding: 10px 0;
        display: block;
        width: 100%;
        overflow-x: auto;
        white-space: nowrap;
    }

    .brick {
        font-size: 14px;
        padding: 4px 12px;
        background-color: #FFF;
        border: 1px solid #E4E7ED;
        border-radius: 20px;
        margin-right: 10px;
        display: inline-block;
        color: #606266;

        span {
            cursor: pointer;
        }
    }

    .close-button {
        color: #909399;
        font-size: 14px;
        float: right;
        cursor: pointer;
        padding-left: 6px;
    }

    .recent-focus {
        font-size: 13px;
        color: #909399;
        margin-top: 20px;
    }

    .collapse-button {
        position: absolute;
        right: 16px;
        top: 12px;
        font-size: 18px;
        color: #909399;
        cursor: pointer;
    }

    .expand-button {
        position: absolute;
        right: 16px;
        top: 0px;
        font-size: 18px;
        color: #909399;
        opacity: 0.15;
        cursor: pointer;
        z-index: 10;

        &:hover {
            opacity: 1;
        }
    }

    .lower-card {
        .el-card__body {
            padding: 8px 8px;
        }
    }

    .chart-card {
        margin-top: 4px;

        > .el-card__body {
            padding: 8px 8px;
            width: 100%;
            height: 50vh;
        }

        width: 100%;

        .echarts {
            width: 100%;
            height: 100%;
        }
    }

    .el-tabs__item {
        height: 20px;
        font-size: 12px;
        line-height: 0;
    }

    .focus-main-number {
        position: relative;
        width: 100%;
        height: 160px;
        padding: 4px 8px;

        p.title {
            font-size: 14px;
            color: #909399;
        }

        p.amount-number {
            font-family: 'Manrope';
            font-weight: bold;
            position: absolute;
            top: 30%;
            height: 80px;

            font-size: 45px;
            font-weight: bold;
            text-align: center;
        }

        p.rate-number {
            font-family: 'Manrope';
            font-weight: bold;
            position: absolute;
            top: 30%;
            height: 80px;

            font-size: 45px;
            font-weight: bold;
            text-align: center;
        }
    }

    .focus-sub-number {
        position: relative;
        width: 100%;
        height: 80px;
        padding: 4px 4px;

        p.title {
            font-size: 12px;
            color: #909399;
        }

        p.percent-number {
            font-family: 'Manrope';
            position: absolute;
            top: 35%;

            font-size: 16px;
            font-weight: bold;
            text-align: center;
            color: #67C23A;
        }

        p.percent-number.neg-number {
            position: absolute;
            top: 35%;

            font-size: 16px;
            font-weight: bold;
            text-align: center;
            color: #F56C6C;
        }
    }

    .plant-list-title {
        font-size: 14px;
        color: #909399;
        width: 100%;
        text-align: center;
        padding: 0 0 4px 0;
    }
}
</style>

<script>
import 'remixicon/fonts/remixicon.css'
import 'vue2-animate/dist/vue2-animate.min.css'
import XTable from '../mixins/SmartVxeTable'
import moment from 'moment'
import { doRequest, message, confirm, allRequests, beautifyNumber, deepCopy } from '../../utils/utils'
import ECharts from 'vue-echarts'
import 'echarts/lib/chart/line'
import 'echarts/lib/component/title'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/toolbox'
import 'echarts/theme/macarons'
import fitty from 'fitty'

export default {
    components: {
        'x-table': XTable,
        'v-chart': ECharts
    },
    mounted() {
        let earlier = new Date();
        earlier.setDate(earlier.getDate() - 30);
        this.dateRange.push(moment(earlier).format('YYYY-MM-DD'));
        this.dateRange.push(moment(Date.now()).format('YYYY-MM-DD'));

        this.getFocuses();
    },
    computed: {
        filteredFocuses() {
            let filtered = this.focuses.filter(data => {
                return !this.focusFilter || data.material.includes(this.focusFilter) ||
                    data.materialTxt.toLowerCase().includes(this.focusFilter.toLowerCase())
            });
            if (filtered.length > 20) {
                filtered = filtered.slice(0, 20);
                filtered.hasMore = true;
            }
            return filtered;
        }
    },
    data() {
        return {
            ff: null,
            dateRange: [],
            material: '',
            focuses: [],
            focusFilter: '',
            focusMap: {},
            showUpper: true,
            upperLoading: false,
            tabs: [],
            tabMap: {},
            activeTab: '',
            lineOptions: {
                title: {
                    text: '销售趋势',
                    textStyle: {
                        fontSize: 16,
                    }
                },
                tooltip: {
                    trigger: 'axis',
                    formatter: function(params) {
                        console.log(params);
                        return 'hellpo'
                    },
                    axisPointer: {
                        type: 'cross',
                        label: {
                            backgroundColor: '#6a7985'
                        }
                    }
                },
                legend: {
                    data: ['销售收入', '销售毛利', '销售毛利率']
                },
                toolbox: {
                    feature: {
                        saveAsImage: {
                            name: '销售趋势'
                        }
                    }
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                xAxis: [
                    {
                        type: 'category',
                        boundaryGap: false,
                        data: []
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '收入',
                        axisLine:{
                          lineStyle:{
                            color:'#ff7c7c',
                          }
                        }
                    }, {
                        type: 'value',
                        name: '毛利',
                        position: 'left',
                        offset: 60,
                        axisLine:{
                          lineStyle:{
                            color:'#60acfc',
                          }
                        }
                    }, {
                        type: 'value',
                        name: '毛利率',
                        axisLine:{
                          lineStyle:{
                            color:'#5bc49f',
                          }
                        }
                    }
                ],
                series: [
                    {
                        name: '销售收入',
                        type: 'line',
                        symbol: 'rectangle',
                        symbolSize: 7,
                        lineStyle: {
                            // color: '#D87A80',
                            color: '#ff7c7c',
                        },
                        itemStyle: {
                            // color: '#D87A80',
                            color: '#ff7c7c',
                            borderWidth: 2,
                            borderColor: '#efa666',
                        },
                        data: []
                    },
                    {
                        name: '销售毛利',
                        type: 'line',
                        symbol: 'circle',
                        symbolSize: 7,
                        lineStyle: {
                            // color: '#FFB980'
                            color: '#60acfc'
                        },
                        itemStyle: {
                            // color: '#FFB980'
                            color: '#60acfc',
                            borderWidth: 2,
                            borderColor: '#9192ab',
                        },
                        yAxisIndex: 1,
                        data: []
                    },
                    {
                        name: '销售毛利率',
                        type: 'line',
                        symbol: 'triangle',
                        symbolSize: 10,
                        lineStyle: {
                            // color: '#5AB1EF'
                            color: '#5bc49f',
                            type:'dashed'
                        },
                        itemStyle: {
                            // color: '#5AB1EF'
                            color: '#5bc49f',
                            borderWidth: 2,
                            borderColor: '#f8cb7f',
                        },
                        yAxisIndex: 2,
                        data: []
                    }
                ]
            }
        }
    },
    methods: {
        copyOptions() {
            let op = deepCopy(this.lineOptions);
            op.tooltip.formatter = function(args) {
                let tooltip = `<p>${args[0].axisValue}</p> `;
                if (args[0])
                    tooltip += `<p>${args[0].marker} ${args[0].seriesName}:
                        ${beautifyNumber(args[0].value.toFixed(2))}</p>`;
                if (args[1])
                    tooltip += `<p>${args[1].marker} ${args[1].seriesName}:
                        ${beautifyNumber(args[1].value.toFixed(2))}</p>`;
                if (args[2])
                    tooltip += `<p>${args[2].marker} ${args[2].seriesName}:
                        ${args[2].value.toFixed(2)}%</p>`;
                return tooltip;
            }
            return op;
        },
        getFocuses() {
            doRequest({
                url: '/v1/web/sales/material/config/list',
                method: 'GET',
                loading: true
            }, {
                success: res => {
                    if (res.cfgs) {
                        this.focuses = res.cfgs.sort((a, b) => {
                            return -a.createdTime.localeCompare(b.createdTime)
                        })
                    }
                    if (this.focuses) {
                        this.focuses.forEach(e => {
                            e.material = e.material.replace(/^0+/g, '')
                            this.focusMap[e.material] = e;
                        })
                    }
                },
                fail: err => {
                    message('error', '获取历史关注失败')
                }
            })
        },
        doFocus() {
            if (this.material.length <= 0 || !this.dateRange ||
                this.dateRange[0].length <= 0 || this.dateRange[1].length <= 0) {
                message('warning', '请输入正确的商品编码及时间范围');
                return;
            }

            this.material = this.material.replace(/^0+/g, '')
            let materials = this.material.split(' ')
            // if (!this.focusMap[this.material]) {
            if (!this.focusMap[materials[0]]) {
                this.upperLoading = true;
                // if (this.focuses.length >= 20) {
                //     doRequest({
                //         url: `/v1/web/sales/material/config/del/${this.focuses[this.focuses.length - 1].material}`
                //     })
                // }
                // doRequest({
                //     url: `/v1/web/basic/material/${this.material}`,
                // }, {
                //     success: res => {
                for(let i=0; i<materials.length; i++){
                    doRequest({
                      url: '/v1/web/sales/material/config/add',
                      method: 'POST',
                      data: {
                        // material: this.material,
                        material: materials[i],
                        materialTxt: '',
                        tag: ''
                      }
                    }, {
                      success: r => {
                        let brick = {
                          // material: this.material,
                          material: materials[i],
                          materialTxt: r.cfg.materialTxt
                        }
                        // this.focusMap[this.material] = brick;
                        this.focusMap[materials[i]] = brick;
                        this.focuses.splice(0, 0, brick)
                        this.viewFocus(brick);
                      },
                      finally: _ => {
                        this.upperLoading = false;
                      }
                    })
                }
                    // },
                    // fail: err => {
                    //     message('error', '关注商品失败，请稍后再试')
                    // },
                    // finally: _ => {
                    //     this.upperLoading = false;
                    // }
                // });
            } else {
              // this.viewFocus(this.focusMap[this.material]);
              for(let i=0; i<materials.length; i++){
                this.viewFocus(this.focusMap[materials[i]]);
              }
            }
        },
        addFocus(v) {
        },
        delFocus(v) {
            delete this.focusMap[v.material];
            for (let i = 0; i < this.focuses.length; i ++) {
                if (this.focuses[i].material == v.material) {
                    this.focuses.splice(i, 1);
                    break;
                }
            }
            doRequest({
                url: `/v1/web/sales/material/config/del/${v.material}`
            }, {
                fail: err => {
                    message('error', '删除历史失败')
                }
            })
        },
        viewFocus(v) {
            if (!this.dateRange ||
                this.dateRange[0].length <= 0 || this.dateRange[1].length <= 0) {
                message('warning', '请输入正确的时间范围');
                return;
            }

            let tabName = `${v.material}-${this.dateRange[0]}-${this.dateRange[1]}`;
            if (this.tabMap[tabName]) {
                this.activeTab = tabName;
                return;
            }

            doRequest({
                url: '/v1/web/sales/material/config/update',
                method: 'POST',
                data: {
                    material: v.material,
                    materialTxt: v.materialTxt,
                    tag: v.tag
                }
            })

            this.material = v.material;
            let tab = {
                material: v.material,
                sdate: this.dateRange[0],
                edate: this.dateRange[1],
                name: tabName,
                label: v.materialTxt,
                loading: true,
                i: this.tabs.length
            }
            this.tabs.push(tab);
            this.tabMap[tabName] = tab;
            this.$nextTick(() => {
                this.activeTab = tabName;
            })
            let self = this;

            allRequests({
                success: res => {
                    // console.log(res);
                    res.forEach(e => {
                        if (e.data.kpi) {
                            tab.net = e.data.kpi.netIncome;
                            tab.margin = e.data.kpi.grossMargin;
                            tab.profit = e.data.kpi.grossProfit;
                            fitty('.amount-number', {
                                maxSize: 60
                            })
                            fitty('.rate-number', {
                                maxSize: 60
                            })
                            fitty('.percent-number', {
                                minSize: 6
                            })
                        } else if (e.data.dateList) {
                            tab.history = e.data.dateList;
                            tab.lineOptions = this.copyOptions(this.lineOptions);
                            tab.history.forEach(e => {
                                tab.lineOptions.title.text = `销售趋势`
                                tab.lineOptions.xAxis[0].data.push(e.Date);
                                tab.lineOptions.series[0].data.push(e.Metric.NetvalInv);
                                tab.lineOptions.series[1].data.push(e.Metric.GrossProfit);
                                tab.lineOptions.series[2].data.push(e.Metric.GrossMargin);
                            })
                        } else if (e.data.plantList) {
                            if (e.config.data.includes('DESC')) {
                                tab.top = e.data.plantList;
                            } else {
                                tab.rop = e.data.plantList;
                            }
                        }
                    })

                    this.$refs['focus-line-chart'][tab.i].mergeOptions(tab.lineOptions, true, false);
                },
                fail: err => {
                    console.log(err);
                },
                finally: _ => {
                    tab.loading = false;
                }
            }, [], [{
                url: '/v1/web/sales/material/kpi',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1],
                    material: v.material
                }
            }, {
                url: '/v1/web/sales/material/date',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1],
                    material: v.material
                }
            }, {
                url: '/v1/web/sales/material/plant/rank',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1],
                    material: v.material,
                    limit: 10,
                    sortBy: 'DESC'
                }
            }, {
                url: '/v1/web/sales/material/plant/rank',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1],
                    material: v.material,
                    limit: 10,
                    sortBy: 'ASC'
                }
            }])
        },
        clearFocus() {
            doRequest({
                url: '/v1/web/sales/material/config/clear',
                loading: true
            }, {
                success: res => {
                    this.focuses = [];
                    this.focusMap = {};
                },
                fail: err => {
                    message('error', '清空历史失败, 请稍后再试')
                }
            })
        },
        clickTab(tab) {
            setTimeout(() => {
                fitty.fitAll();
            }, 200);
        },
        removeTab(name) {
            this.tabs.splice(this.tabMap[name].i, 1);
            delete this.tabMap[name];
            if (this.activeTab == name) {
                if (this.tabs.length > 0) {
                    this.activeTab = this.tabs[0].name;
                }
            }

            let i = 0;
            this.tabs.forEach(e => {
                e.i = i;
                i += 1;
            })
        },
        beautyNum(n) {
            return beautifyNumber(n)
        }
    }
}
</script>
