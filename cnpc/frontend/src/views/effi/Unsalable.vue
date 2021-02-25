<template>
    <div class="cnpn-unsalable-div">
        <el-row :gutter="8" style="margin-bottom: 8px">
            <el-col :span="12">
                <el-card body-style="padding: 10px 10px" >
                    <el-date-picker size="mini" style="padding-top: 2px; margin-right: 4px"
                        v-model="dateRange"
                        type="daterange"
                        range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd" />
                    <el-button size="mini" type="primary" icon="el-icon-view" @click="getStaleGoods"
                        :disabled="tableLoading" v-loading="tableLoading">查看清单</el-button>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card body-style="padding: 10px 10px" >
                    <el-input-number size="mini" v-model="threshold" style="margin-right: 4px; width: 50%"/>
                    <el-button size="mini" type="success" icon="el-icon-upload2" @click="updateThreshold"
                        :disabled="threshLoading" v-loading="threshLoading">修改动销率阈值</el-button>
                </el-card>
            </el-col>
        </el-row>
        
        <div class="middle-row">
            <el-tag type="primary" size="medium">
                {{ `${sdate} ~ ${edate}` }}
            </el-tag>
            <el-button size="mini" icon="el-icon-download" type="success" 
                style="margin-left: 8px" @click="exportReport" 
                :disabled="exportLoading || !sdate || !edate" v-loading="exportLoading">导出清单</el-button>
        </div>
        <x-table :minus="125" :data="goods" v-loading="tableLoading" size="mini" show-overflow="tooltip"
            :highlight-current-row="true">
            <vxe-table-column type="seq" fixed="left" width="80"/>
            <vxe-table-column title="商品编码" field="material" fixed="left" width="160" sortable/>
            <vxe-table-column title="商品名称" field="materialTxt" sortable/>
            <vxe-table-column title="销售额" field="metric.netSales" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.netSales) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="库存金额" field="metric.invCost" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.invCost) }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="动销率" field="metric.salesActRate" width="120" sortable>
                <template slot-scope="r">
                    {{ beautyNum(r.row.metric.salesActRate * 100) }}%
                </template>
            </vxe-table-column>
        </x-table>
    </div>
</template>

<style lang="scss">
.cnpn-unsalable-div {
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
import { doRequest, allRequests, message, confirm, beautifyNumber, downloadFile } from '../../utils/utils'
import XTable from '../mixins/SmartVxeTable'
import moment from 'moment'

export default {
    components: {
        'x-table': XTable
    },
    data() {
        return {
            dateRange: [],
            goods: [],
            tableLoading: false,
            exportLoading: false,
            sdate: '',
            edate: '',
            threshName: '',
            threshold: 0,
            threshLoading: false,
        }
    },
    mounted() {
        let earlier = new Date();
        earlier.setDate(earlier.getDate() - 30);
        this.dateRange.push(moment(earlier).format('YYYY-MM-DD'));
        this.dateRange.push(moment(Date.now()).format('YYYY-MM-DD'));

        //this.getStaleGoods();
        this.getThreshold();
    },
    methods: {
        getThreshold() {
            this.threshLoading = true
            doRequest({
                url: '/v1/web/eff/metric/config/list/EFF001'
            }, {
                success: res => {
                    if (res.cfgs && res.cfgs.length > 0) {
                        this.threshold = res.cfgs[0].value
                        this.threshName = res.cfgs[0].name
                    }
                },
                fail: err => {
                    message('error', '获取动销率阈值失败, 请稍后再试')
                },
                finally: _ => {
                    this.threshLoading = false;
                }
            })
        },
        updateThreshold() {
            this.threshLoading = true
            doRequest({
                url: '/v1/web/eff/metric/config/update',
                method: 'POST',
                data: {
                    name: this.threshName,
                    value: this.threshold
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
        getStaleGoods() {
            if (!this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('error', '请选择有效的时间范围')
                return
            }

            this.tableLoading = true
            this.sdate = this.dateRange[0]
            this.edate = this.dateRange[1]
            doRequest({
                url: '/v1/web/eff/unsalable/material',
                method: 'POST',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1]
                }
            }, {
                success: res => {
                    if (res.matlList) {
                        this.goods = res.matlList.sort((a, b) => {
                            return a.metric.salesActRate - b.metric.salesActRate
                        })
                    }
                },
                fail: err => {
                    message('error', '获取清单失败, 请稍后再试')
                },
                finally: _ => {
                    this.tableLoading = false;
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
                url: '/v1/web/eff/unsalable/material/export',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate
                }
            }, {
                fail: err => {
                    message('error', '清单导出失败, 请稍后再试')
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