<template>
    <div class="cnpc-stock-alarm-div">
        <el-row :gutter="8" style="margin-bottom: 8px">
            <el-col :span="12">
                <el-card body-style="padding: 25px 10px; height: 80px" >
                    <el-date-picker size="mini" style="padding-top: 2px; margin-right: 4px"
                        v-model="dateRange"
                        type="daterange"
                        range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd" />
                    <el-button size="mini" type="primary" icon="el-icon-view" @click="getAlarms"
                        :disabled="tableLoading" v-loading="tableLoading">查看预警</el-button>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card body-style="padding: 10px 10px">
                    <el-input-number size="mini" v-model="threshDay" style="margin-right: 4px; margin-bottom: 4px; width: 50%"/>
                    <el-button size="mini" type="success" icon="el-icon-upload2" 
                        @click="updateThreshold('eff_dsi_th', threshDay)"
                        :disabled="threshLoading" v-loading="threshLoading">修改周转天数阈值</el-button>
                    <el-input-number size="mini" v-model="threshCost" style="margin-right: 4px; width: 50%"/>
                    <el-button size="mini" type="success" icon="el-icon-upload2" 
                        @click="updateThreshold('eff_invcost_th', threshCost)"
                        :disabled="threshLoading" v-loading="threshLoading">修改库存金额阈值</el-button>
                </el-card>
            </el-col>
        </el-row>
        <div class="middle-row">
            <el-tag type="primary" size="medium">
                {{ `${sdate} ~ ${edate}` }}
            </el-tag>
            <el-button size="mini" icon="el-icon-download" type="warning" 
                style="margin-left: 8px" @click="exportReport" 
                :disabled="exportLoading || !sdate || !edate" v-loading="exportLoading">导出预警</el-button>
        </div>
        <x-table :minus="155" :data="alarms" v-loading="tableLoading" size="mini" show-overflow="tooltip">
            <vxe-table-column type="seq" width="80"/>
            <vxe-table-column title="商品编码" field="material" width="160" sortable/>
            <vxe-table-column title="商品名称" field="materialTxt" sortable/>
            <vxe-table-column title="销售成本" field="metric.cost" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.cost) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="期初库存" field="metric.openInvCost" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.openInvCost) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="期末库存" field="metric.closeInvCost" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.closeInvCost) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="周转天数" field="metric.daysSalesInv" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.daysSalesInv) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="平均库存金额" field="metric.avgInvCost" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.avgInvCost) }}
                </template>
            </vxe-table-column>
        </x-table>
    </div>
</template>

<style lang="scss">
.cnpc-stock-alarm-div {
    width: 100%;
    max-height: calc(100vh - 12px);
    background-color: #FFF;
    padding: 8px 8px;
    overflow: auto;

    .middle-row {
        margin: 8px 0;
    }
}
</style>

<script>
import {doRequest, message, confirm, beautifyNumber, downloadFile } from '../../utils/utils'
import XTable from '../mixins/SmartVxeTable'
import moment from 'moment'

export default {
    components: {
        'x-table': XTable
    },
    data() {
        return {
            threshDay: 0,
            threshCost: 0,
            threshLoading: false,
            dateRange: [],
            alarms: [],
            tableLoading: false,
            sdate: '',
            edate: '',
            exportLoading: false,
        }
    },
    mounted() {
        let earlier = new Date();
        earlier.setDate(earlier.getDate() - 30);
        this.dateRange.push(moment(earlier).format('YYYY-MM-DD'));
        this.dateRange.push(moment(Date.now()).format('YYYY-MM-DD'));

        //this.getAlarms()
        this.getThresholds()
    },
    methods: {
        getAlarms() {
            if (!this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的时间范围')
                return
            }

            this.tableLoading = true;
            this.sdate = this.dateRange[0]
            this.edate = this.dateRange[1]
            doRequest({
                url: '/v1/web/eff/inventory/material',
                method: 'POST',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1]
                }
            }, {
                success: res => {
                    if (res.matlList) {
                        this.alarms = res.matlList.sort((a, b) => {
                            return - (a.metric.avgInvCost - b.metric.avgInvCost)
                        })
                    }
                },
                fail: err => {
                    message('error', '获取预警列表失败, 请稍后再试')
                },
                finally: _ => {
                    this.tableLoading = false
                }
            })
        },
        getThresholds() {
            this.threshLoading = true
            doRequest({
                url: '/v1/web/eff/metric/config/list/EFF002'
            }, {
                success: res => {
                    if (res.cfgs && res.cfgs.length >= 2) {
                        res.cfgs.forEach(e => {
                            if (e.name == 'eff_dsi_th') {
                                this.threshDay = e.value
                            } else {
                                this.threshCost = e.value
                            }
                        })
                    }
                },
                fail: _ => {
                    message('error', '获取阈值失败, 请稍后再试')
                },
                finally: _ => {
                    this.threshLoading = false;
                }
            })
        },
        updateThreshold(name, value) {
            this.threshLoading = true
            doRequest({
                url: '/v1/web/eff/metric/config/update',
                method: 'POST',
                data: {
                    name: name,
                    value: value
                }
            }, {
                success: _ => {
                    message('success', '修改阈值成功')
                },
                fail: _ => {
                    message('error', '修改阈值失败, 请稍后再试')
                },
                finally: _ => {
                    this.threshLoading = false;
                }
            })
        },
        exportReport() {
            if (!this.sdate || !this.edate) {
                message('warning', '请选择有效的时间范围, 并生成报告后导出')
                return
            }
            this.exportLoading = true
            downloadFile({
                url: '/v1/web/eff/inventory/material/export',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate
                }
            }, {
                fail: err => {
                    message('error', '预警导出失败, 请稍后再试')
                },
                finally: _ => {
                    this.exportLoading = false
                }
            })
        },
        beautyNum(n) {
            if (isNaN(n))
                return n
            return beautifyNumber(n.toFixed(2))
        }
    }
}
</script>