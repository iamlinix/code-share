<template>
    <div class="kpi-analyse-div">
        <el-card>
            <el-form inline size="mini" label-width="80px">
                <el-form-item style="margin-bottom: 0" label="统计组织">
                    <el-select v-model="selectedOrg" filterable clearable style="width: 250px">
                        <el-option v-for="(v, i) in orgs" :key="v.code" :value="v.code" :label="v.name"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="统计日期" style="margin-bottom: 0">
                    <el-date-picker
                        v-model="dateRange"
                        type="daterange"
                        range-separator="-"
                        format="yyyy-MM-dd"
                        value-format="yyyy-MM-dd"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期">
                    </el-date-picker>
                </el-form-item>
                <el-form-item label="同比时段" style="margin-bottom: 0">
                    <el-date-picker v-model="startYear" type="year" value-format="yyyy" style="width: 100px"/>
                    -
                    <el-date-picker v-model="endYear" type="year" value-format="yyyy" style="width: 100px"/>
                </el-form-item>
                <el-button type="primary" size="mini" @click="analyse">开始分析</el-button>
            </el-form>
        </el-card>
        <el-card style="margin-top: 8px">
            <p style="font-size: 22px; font-weight: bold; color: #409EFF; margin-bottom: 10px">{{ selectedOrgName }}</p>
            <el-table :border="false" style="margin-bottom: 16px" :data="models" :show-header="false">
                <el-table-column width="130">
                    <template slot-scope="r">
                        <p style="font-size: 17px; font-weight: 100">{{ r.row.name }}</p>
                    </template>
                </el-table-column>
                <el-table-column width="300">
                    <template slot-scope="r">
                        <div v-if="r.$index == 0">
                            <el-tag effect="dark" :type="upndowns[0] ? 'success' : 'danger'">进站人数</el-tag>
                            <el-tag effect="dark" :type="upndowns[2] ? 'success' : 'danger'" style="margin-left: 6px">客单价</el-tag>
                            <el-tag effect="dark" :type="upndowns[3] ? 'success' : 'danger'" style="margin-left: 6px">油非转化</el-tag>
                        </div>
                        <div v-if="r.$index == 1">
                            <el-tag effect="dark" :type="upndowns[4] ? 'success' : 'danger'">销售收入</el-tag>
                            <el-tag effect="dark" :type="upndowns[5] ? 'success' : 'danger'" style="margin-left: 6px">销售毛利</el-tag>
                        </div>
                        <div v-if="r.$index == 2">
                            <el-tag effect="dark" :type="upndowns[1] ? 'success' : 'danger'">进店人数</el-tag>
                            <el-tag effect="dark" :type="upndowns[2] ? 'success' : 'danger'" style="margin-left: 6px">客单价</el-tag>
                            <el-tag effect="dark" :type="upndowns[4] ? 'success' : 'danger'" style="margin-left: 6px">销售收入</el-tag>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column>
                    <template slot-scope="r">
                        <div v-if="r.$index == 0">
                            <p style="font-size: 15px; font-weight: 100">{{ advice1[upndowns[0] * 4 + upndowns[2] * 2 + upndowns[3] * 1] }}</p>
                        </div>
                        <div v-if="r.$index == 1">
                            <p style="font-size: 15px; font-weight: 100">{{ advice2[upndowns[4] * 2 + upndowns[5] * 1] }}</p>
                        </div>
                        <div v-if="r.$index == 2">
                            <p style="font-size: 15px; font-weight: 100">{{ advice3[upndowns[1] * 4 + upndowns[2] * 2 + upndowns[4] * 1] }}</p>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-row :gutter="8">
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">进站人数</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.fuelCount && currentYear.fuelCount.value >= customerCountAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.fuelCount && currentYear.fuelCount.value >= customerCountAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.fuelCount ? ((currentYear.fuelCount.value - customerCountAvg) * 100 / customerCountAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="custom-count-line-chart" theme="macarons" :options="customerCountOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">进店人数</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.nonFuelCount && currentYear.nonFuelCount.value >= nfCustomerCountAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.nonFuelCount && currentYear.nonFuelCount.value >= nfCustomerCountAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.nonFuelCount ? ((currentYear.nonFuelCount.value - nfCustomerCountAvg) * 100 / nfCustomerCountAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="custom-count-line-chart" theme="macarons" :options="nfCustomerCountOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">客单价</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.avgTrxValue && currentYear.avgTrxValue.value >= customerAvgAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.avgTrxValue && currentYear.avgTrxValue.value >= customerAvgAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.avgTrxValue ? ((currentYear.avgTrxValue.value - customerAvgAvg) * 100 / customerAvgAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="custom-avg-line-chart" theme="macarons" :options="customerAvgOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">油非转化</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.fnConversionRate && currentYear.fnConversionRate.value >= customerTransferAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.fnConversionRate && currentYear.fnConversionRate.value >= customerTransferAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.fnConversionRate ? ((currentYear.fnConversionRate.value - customerTransferAvg) * 100 / customerTransferAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="custom-transfer-line-chart" theme="macarons" :options="customerTransferOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">销售收入</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.netIncome && currentYear.netIncome.value >= incomeAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.netIncome && currentYear.netIncome.value >= incomeAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.netIncome ? ((currentYear.netIncome.value - incomeAvg) * 100 / incomeAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="income-line-chart" theme="macarons" :options="incomeOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
                <el-col :span="12" style="margin-bottom: 8px">
                    <el-card>
                        <el-row :gutter="6" style="padding-bottom: 8px">
                            <el-col :span="5">
                                <p style="padding-left: 8px; color: #66B1FF; font-size: 18px; font-weight: bold">销售毛利</p>
                            </el-col>
                            <el-col :span="18" :style="'font-size: 18px; font-weight: bold;' + ((currentYear.grossMargin && currentYear.grossMargin.value >= profitAvg) ? 'color: #67C23A' : 'color: #F56C6C')">
                                <font-awesome-icon v-if="currentYear.grossMargin && currentYear.grossMargin.value >= profitAvg" icon="arrow-up" style="color: #67C23A"/>
                                <font-awesome-icon v-else icon="arrow-down" style="color: #F56C6C"/>
                                {{ currentYear.grossMargin ? ((currentYear.grossMargin.value - profitAvg) * 100 / profitAvg).toFixed(2) : 0 }}%
                            </el-col>
                        </el-row>
                        <el-card :body-style="{padding: '8px 8px', height: '350px'}">
                            <v-chart ref="profit-line-chart" theme="macarons" :options="profitOption" :autoresize="true"/>
                        </el-card>
                    </el-card>
                </el-col>
            </el-row>
        </el-card>
    </div>
</template>

<script>
import { doRequest, message, deepCopy, beautifyNumber } from '../../utils/utils'
import { getUserInfo } from '../../utils/dataStorage'
import ECharts from 'vue-echarts'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/bar'
import 'echarts/lib/chart/pie'
import 'echarts/lib/component/title'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/toolbox'
import 'echarts/theme/macarons'

export default {
    components: {
        'v-chart': ECharts
    },
    mounted() {
        doRequest({
            method: 'GET',
            url: `/v1/web/user/org-perm/${getUserInfo().name}`,
            loading: true
        }, {
            success: res => {
                let perms = []
                if (res.orgPerm.length > 0) {
                    perms = JSON.parse(res.orgPerm)
                }

                doRequest({
                    method: 'GET',
                    url: '/v1/web/org/plant',
                    loading: true
                }, {
                    success: res => {
                        let orgs = [{
                            level: 0,
                            code: '000',
                            name: '北京公司'
                        }, {
                            level: 1,
                            code: 'A13',
                            name: '一分公司'
                        }, {
                            level: 1,
                            code: 'A14',
                            name: '二分公司'
                        }, {
                            level: 1,
                            code: 'A15',
                            name: '三分公司'
                        }, {
                            level: 1,
                            code: 'A16',
                            name: '四分公司'
                        }]
                        res.orgList.forEach(e => {
                            orgs.push({
                                level: 2,
                                code: e.orgCode,
                                name: e.orgText
                            })
                        })

                        if (perms.length > 0) {
                            orgs.forEach(e => {
                                perms.forEach(p => {
                                    if (p.orgCode == e.code) {
                                        if (p.show || (p.show == null))
                                            this.orgs.push(e)
                                    }
                                })
                            })
                        } else {
                            this.orgs = orgs
                        }
                    },
                    fail: _ => {
                        message('error', '获取油站列表失败，请稍后再试')
                    }
                })
            },
            fail: _ => {
                message('error', '获取用户机构权限失败, 请稍后再试')
            }
        })
    },
    data() {
        return {
            models: [{
                name: '客户流量模型',
            }, {
                name: '基础销售模型',
            }, {
                name: '基础非油模型',
            }],
            advice1: [
                '开展促销活动，提升人流量',
                '更改促销方案，替换为高利润商品',
                '更改促销方案，替换为低价格商品',
                '加大优惠力度或者扩大优惠种类',
                '清理滞销商品，上架畅销新品',
                '上架更多高利润商品',
                '加大推广推销力度',
                '请保持'
            ],
            advice2: [
                '开展促销活动，提升人流量',
                '加大促销活动优惠力度',
                '替换低利润率商品',
                '请保持'
            ],
            advice3: [
                '开展促销活动，提升人流量',
                '更改促销方案，替换为高利润商品',
                '对畅销商品进行促销',
                '加大推广推销力度',
                '选用更高利润率商品',
                '选用更高利润率商品',
                '适当降低促销力度',
                '请保持'
            ],
            upndowns: [
               false, false, false, false, false, false
            ],
            dateRange: [],
            startYear: '',
            endYear: '',
            currentYear: {},
            pastYears: [],
            selectedOrg: '',
            selectedOrgName: '',
            orgs: [],
            customerCountAvg: 0,
            nfCustomerCountAvg: 0,
            customerAvgAvg: 0,
            customerTransferAvg: 0,
            incomeAvg: 0,
            profitAvg: 0,
            customerCountOption: {},
            nfCustomerCountOption: {},
            customerAvgOption: {},
            customerTransferOption: {},
            incomeOption: {},
            profitOption: {},
            metaBarOption: {
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        label: {
                            backgroundColor: '#6a7985'
                        }
                    }
                },
                grid: {
                    left: '3%',
                    right: '90',
                    bottom: '3%',
                    containLabel: true
                },
                legend: {
                    data: ['今年', '往年'],
                    icon: 'roundRect',
                    textStyle: {
                        color: '#000'
                    }
                },
                xAxis: [
                    {
                        type: 'category',
                        data: [],
                        axisTick: {
                            alignWithLabel: true
                        }
                    }
                ],
                yAxis: [
                    {
                        type: 'value'
                    }
                ],
                series: [
                    {
                        name: '往年',
                        type: 'bar',
                        barWidth: '40%',
                        data: [],
                        markLine: {
                            lineStyle: {
                                type: 'dashed'
                            },
                            data: [{
                                name: '平均',
                                type: 'average'
                            }]
                        }
                    }, {
                        name: '今年',
                        type: 'line',
                        data: [],
                    }
                ]
            }
        }
    },
    methods: {
        xxx(r){
            console.log(r)
        },
        formatToolTip(args) {
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
        },
        analyse() {
            if (this.dateRange.length == 0 || this.startYear.length == 0 || this.endYear.length == 0 || this.selectedOrg.length == 0) {
                message('warning', '所有选项均不能为空')
                return
            }

            let org = null;
            this.orgs.forEach(e => {
                if (e.code == this.selectedOrg) {
                    org = e
                    return
                }
            })

            let syear = parseInt(this.startYear), eyear = parseInt(this.endYear), year = parseInt(this.dateRange[0].substring(0, 5))
            if (syear > eyear) {
                message('warning', '同比开始时间大于结束时间')
                return
            }

            if (eyear >= year) {
                message('warning', '同比时段结束时间应小于统计日期开始时间')
                return
            }

            if (syear < 2012) {
                message('warning', '同比时段不能早于 2012 年')
                return
            }

            doRequest({
                method: 'POST',
                url: '/v1/web/sales/kpi/analyse',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1],
                    orgLevel: org.level,
                    orgCode: org.code,
                    startYear: syear,
                    endYear: eyear
                },
                loading: true
            }, {
                success: res => {
                    this.customerCountAvg = this.customerAvgAvg = this.customerTransferAvg = this.incomeAvg = this.profitAvg = this.nfCustomerCountAvg = 0

                    this.selectedOrgName = org.name
                    this.currentYear = res.currentKpi
                    this.pastYears = res.pastYears

                    let xAxisData = []
                    for (let i = syear; i <= eyear; i ++) {
                        xAxisData.push(i + '')
                    }

                    this.customerCountOption = deepCopy(this.metaBarOption)
                    this.nfCustomerCountOption = deepCopy(this.metaBarOption)
                    this.customerAvgOption = deepCopy(this.metaBarOption)
                    this.customerTransferOption = deepCopy(this.metaBarOption)
                    this.incomeOption = deepCopy(this.metaBarOption)
                    this.profitOption = deepCopy(this.metaBarOption)

                    this.customerCountOption.tooltip.formatter = this.formatToolTip
                    this.nfCustomerCountOption.tooltip.formatter = this.formatToolTip
                    this.customerAvgOption.tooltip.formatter = this.formatToolTip
                    this.customerTransferOption.tooltip.formatter = this.formatToolTip
                    this.incomeOption.tooltip.formatter = this.formatToolTip
                    this.profitOption.tooltip.formatter = this.formatToolTip

                    this.customerCountOption.xAxis[0].data = xAxisData
                    this.nfCustomerCountOption.xAxis[0].data = xAxisData
                    this.customerAvgOption.xAxis[0].data = xAxisData
                    this.customerTransferOption.xAxis[0].data = xAxisData
                    this.incomeOption.xAxis[0].data = xAxisData
                    this.profitOption.xAxis[0].data = xAxisData

                    res.pastYears.forEach(e => {
                        this.customerCountAvg += e.fuelCount.value
                        this.customerCountOption.series[0].data.push(e.fuelCount.value)
                        this.customerCountOption.series[1].data.push(res.currentKpi.fuelCount.value)

                        this.nfCustomerCountAvg += e.nonFuelCount.value
                        this.nfCustomerCountOption.series[0].data.push(e.nonFuelCount.value)
                        this.nfCustomerCountOption.series[1].data.push(res.currentKpi.nonFuelCount.value)

                        this.customerAvgAvg += e.avgTrxValue.value
                        this.customerAvgOption.series[0].data.push(e.avgTrxValue.value)
                        this.customerAvgOption.series[1].data.push(res.currentKpi.avgTrxValue.value)

                        this.customerTransferAvg += e.fnConversionRate.value
                        this.customerTransferOption.series[0].data.push(e.fnConversionRate.value)
                        this.customerTransferOption.series[1].data.push(res.currentKpi.fnConversionRate.value)

                        this.incomeAvg += e.netIncome.value
                        this.incomeOption.series[0].data.push(e.netIncome.value)
                        this.incomeOption.series[1].data.push(res.currentKpi.netIncome.value)

                        this.profitAvg += e.grossMargin.value
                        this.profitOption.series[0].data.push(e.grossMargin.value)
                        this.profitOption.series[1].data.push(res.currentKpi.grossMargin.value)
                    })

                    this.customerCountAvg /= res.pastYears.length
                    this.nfCustomerCountAvg /= res.pastYears.length
                    this.customerAvgAvg /= res.pastYears.length
                    this.customerTransferAvg /= res.pastYears.length
                    this.incomeAvg /= res.pastYears.length
                    this.profitAvg /= res.pastYears.length

                    this.$set(this.upndowns, 0, this.currentYear.fuelCount.value >= this.customerCountAvg)
                    this.$set(this.upndowns, 1, this.currentYear.nonFuelCount.value >= this.nfCustomerCountAvg)
                    this.$set(this.upndowns, 2, this.currentYear.avgTrxValue.value >= this.customerAvgAvg)
                    this.$set(this.upndowns, 3, this.currentYear.fnConversionRate.value >= this.customerTransferAvg)
                    this.$set(this.upndowns, 4, this.currentYear.netIncome.value >= this.incomeAvg)
                    this.$set(this.upndowns, 5, this.currentYear.grossMargin.value >= this.profitAvg)
                },
                fail: err => {
                    console.log(err)
                    message('error', '分析指标失败，请稍后再试')
                }
            })
        }
    }
}
</script>

<style lang="scss">
.kpi-analyse-div {
    .echarts {
        width: 100%;
        height: 100%;
    }

    table {
        border: none
    }
    td {
        border: none
    }

    max-height: calc(100vh - 20px);
    overflow-y: auto;
}
</style>
