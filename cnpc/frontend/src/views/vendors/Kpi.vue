<template>
    <div class="vendor-kpi-main">
            <el-tabs type="card" v-model="activeTab" @tab-click="tabClick" @tab-remove="tabRemove">
                <el-tab-pane name="report-config" label="参数配置">
                    <el-steps :active="5" simple>
                        <el-step title="设置权重" icon="el-icon-setting" description="hellpo"></el-step>
                        <el-step title="生成预览" icon="el-icon-view"></el-step>
                        <el-step title="修改主观得分" icon="el-icon-edit"></el-step>
                        <el-step title="提交正式报告" icon="el-icon-upload"></el-step>
                        <el-step title="查看/导出报告" icon="el-icon-s-data"></el-step>
                    </el-steps>
                    <el-row :gutter="4">
                        <el-col :span="5">
                            <el-card>
                                <p class="weight-header">权重分配</p>
                                <el-form size="mini" label-position="right" label-width="90px">
                                    <el-form-item v-for="v in weights" :key="'vdr-kpi-wt-' + v.name"
                                        :label="v.name">
                                        <el-input-number size="mini" :min="0" v-model="v.weight" 
                                        style="width: 100%"/>
                                    </el-form-item>
                                    <el-button size="mini" icon="el-icon-setting" style="width: 100%"
                                        type="success" @click="updateWeight">更新权重</el-button>
                                </el-form>
                            </el-card>
                            <el-card style="margin-top: 4px">
                                <el-date-picker size="mini" 
                                    style="margin-bottom: 8px; width: 100%  "
                                    type="daterange" v-model="dateRange" range-separator="-"
                                    start-placeholder="开始日期"
                                    end-placeholder="结束日期"
                                    value-format="yyyy-MM-dd" />
                                <el-button type="primary" size="mini" @click="getPreview"
                                    icon="el-icon-view" style="width: 100%">生成预览报告</el-button>
                            </el-card>
                            <el-card style="margin-top: 4px">
                                <el-input v-model="reportName" size="mini" 
                                    style="width: 100%; margin-bottom: 8px"
                                    placeholder="请输入为新报告输入一个名称"/>
                                <el-button size="mini" type="success" style="width: 100%"
                                    icon="el-icon-upload" @click="createNewReport">提交正式报告</el-button>
                            </el-card>
                        </el-col>
                        <el-col :span="19">
                            <split-pane :min-percent='10' :default-percent='50' split="vertical">
                                <template slot="paneL">
                                    <el-card>
                                        <p class="weight-header">预览报告</p>
                                        <x-table :minus="195" :data="preview" size="mini" v-loading="previewLoading"
                                            show-overflow="tooltip" highlight-current-row ref="preview-report-list">
                                            <vxe-table-column fixed="left" type="seq" width="45"/>
                                            <vxe-table-column sortable title="供应商编码" fixed="left" field="vendor" width="130"/>
                                            <vxe-table-column sortable title="供应商名称" field="vendorName" width="300"/>
                                            <vxe-table-column sortable title="销售收入" field="metric.netSales" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.netSales) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="销售收入占比" field="metric.netSalesPct" width="120">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.netSalesPct) }}%
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="销售收入得分" field="metric.netSalesScore" width="120">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.netSalesScore) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="销售收入加权" field="metric.netSalesScoreW" width="120">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.netSalesScoreW) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利" field="metric.grossProfit" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossProfit) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利占比" field="metric.grossProfitPct" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossProfitPct) }}%
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利得分" field="metric.grossProfitScore" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossProfitScore) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利加权" field="metric.grossProfitScoreW" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossProfitScoreW) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利率" field="metric.grossMargin" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossMargin) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利率得分" field="metric.grossMarginScore" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossMarginScore) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="毛利率加权" field="metric.grossMarginScoreW" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.grossMarginScoreW) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="到货率" field="metric.orderFillRate" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.orderFillRate) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="到货率得分" field="metric.orderFillRateScore" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.orderFillRateScore) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="到货率加权" field="metric.orderFillRateScoreW" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.orderFillRateScoreW) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="主观得分" field="metric.subjectiveScore" width="180">
                                                <template slot-scope="r">
                                                    <el-input-number size="mini" :min="0" :step="0.01" 
                                                        v-model="r.row.metric.subjectiveScore" 
                                                        @change="onSubjectChange(r.row)" />
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="主观加权" field="metric.subjectiveScoreW" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.subjectiveScoreW) }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column sortable title="总分" field="metric.totalScore" width="100">
                                                <template slot-scope="r">
                                                    {{ beautyNum(r.row.metric.totalScore) }}
                                                </template>
                                            </vxe-table-column>
                                        </x-table>
                                    </el-card>
                                </template>
                                <template slot="paneR">
                                    <el-card>
                                        <p class="weight-header">已生成报告</p>
                                        <x-table :minus="195" :data="generated" size="mini" v-loading="generatedLoading"
                                            show-overflow="tooltip" highlight-current-row ref="generated-report-list">
                                            <vxe-table-column title="报告名称" sortable field="name" min-width="120"/>
                                            <vxe-table-column title="起始时间" sortable field="beginDate" width="160">
                                                <template slot-scope="r">
                                                    {{ r.row.beginDate + '~' + r.row.endDate }}
                                                </template>
                                            </vxe-table-column>
                                            <vxe-table-column title="操作" width="60" align="center">
                                                <template slot-scope="r">
                                                    <i class="el-icon-s-data" 
                                                        style="color: #67C23A; cursor: pointer; font-size: 14px"
                                                        @click="viewReport(r.row)"/>
                                                    <i class="el-icon-close" 
                                                        style="color: #F56C6C; cursor: pointer; font-size: 14px; margin-left: 4px"
                                                        @click="deleteReport(r.row)"/>
                                                </template>
                                            </vxe-table-column>
                                        </x-table>
                                    </el-card>
                                </template>
                            </split-pane>
                        </el-col>
                        <!--el-col :span="12">
                            <el-card>
                                <p class="weight-header">预览报告</p>
                                <x-table :minus="195" :data="preview" size="mini" v-loading="previewLoading"
                                    show-overflow="tooltip" highlight-current-row ref="preview-report-list">
                                    <vxe-table-column fixed="left" type="seq" width="45"/>
                                    <vxe-table-column sortable title="供应商编码" fixed="left" field="vendor" width="130"/>
                                    <vxe-table-column sortable title="供应商名称" field="vendorName" width="300"/>
                                    <vxe-table-column sortable title="销售收入" field="metric.netSales" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.netSales) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="销售收入占比" field="metric.netSalesPct" width="120">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.netSalesPct) }}%
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="销售收入得分" field="metric.netSalesScore" width="120">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.netSalesScore) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="销售收入加权" field="metric.netSalesScoreW" width="120">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.netSalesScoreW) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利" field="metric.grossProfit" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossProfit) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利占比" field="metric.grossProfitPct" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossProfitPct) }}%
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利得分" field="metric.grossProfitScore" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossProfitScore) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利加权" field="metric.grossProfitScoreW" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossProfitScoreW) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利率" field="metric.grossMargin" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossMargin) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利率得分" field="metric.grossMarginScore" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossMarginScore) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="毛利率加权" field="metric.grossMarginScoreW" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.grossMarginScoreW) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="到货率" field="metric.orderFillRate" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.orderFillRate) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="到货率得分" field="metric.orderFillRateScore" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.orderFillRateScore) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="到货率加权" field="metric.orderFillRateScoreW" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.orderFillRateScoreW) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="主观得分" field="metric.subjectiveScore" width="180">
                                        <template slot-scope="r">
                                            <el-input-number size="mini" :min="0" :step="0.01" 
                                                v-model="r.row.metric.subjectiveScore" 
                                                @change="onSubjectChange(r.row)" />
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="主观加权" field="metric.subjectiveScoreW" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.subjectiveScoreW) }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column sortable title="总分" field="metric.totalScore" width="100">
                                        <template slot-scope="r">
                                            {{ beautyNum(r.row.metric.totalScore) }}
                                        </template>
                                    </vxe-table-column>
                                </x-table>
                            </el-card>
                        </el-col>
                        <el-col :span="7">
                            <el-card>
                                <p class="weight-header">已生成报告</p>
                                <x-table :minus="195" :data="generated" size="mini" v-loading="generatedLoading"
                                    show-overflow="tooltip" highlight-current-row ref="generated-report-list">
                                    <vxe-table-column title="报告名称" sortable field="name" />
                                    <vxe-table-column title="起始时间" sortable field="beginDate" width="160">
                                        <template slot-scope="r">
                                            {{ r.row.beginDate + '~' + r.row.endDate }}
                                        </template>
                                    </vxe-table-column>
                                    <vxe-table-column title="操作" width="60" align="center">
                                        <template slot-scope="r">
                                            <i class="el-icon-s-data" 
                                                style="color: #67C23A; cursor: pointer; font-size: 14px"
                                                @click="viewReport(r.row)"/>
                                            <i class="el-icon-close" 
                                                style="color: #F56C6C; cursor: pointer; font-size: 14px; margin-left: 4px"
                                                @click="deleteReport(r.row)"/>
                                        </template>
                                    </vxe-table-column>
                                </x-table>
                            </el-card>
                        </el-col-->
                    </el-row>
                </el-tab-pane>
                <el-tab-pane v-for="(t, i) in tabs" :key="'vdr-kpi-tab-' + i" closable v-loading="t.loading" 
                    :name="t.name" :label="t.label">
                    <el-button size="mini" type="success" icon="el-icon-download" 
                        @click="exportReport(t)">导出报告</el-button>
                    <el-tag size="mini" type="primary" style="margin: 0 8px 8px 8px;">{{ t.label }}</el-tag>
                    <el-tag size="mini" type="success">{{ t.sdate + ' ~ ' + t.edate }}</el-tag>
                    <x-table :data="t.report" :minus="125" :ref="`vdr-kpi-rpt-tbl-${i}`"
                        size="mini" show-overflow="tooltip" highlight-current-row>
                        <vxe-table-column fixed="left" type="seq" width="45"/>
                        <vxe-table-column sortable title="供应商编码" fixed="left" field="vendor" width="130"/>
                        <vxe-table-column sortable title="供应商名称" fixed="left" field="vendorName" width="300"/>
                        <vxe-table-column sortable title="销售收入" field="metric.netSales" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.netSales) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="销售收入占比" field="metric.netSalesPct" width="120">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.netSalesPct) }}%
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="销售收入得分" field="metric.netSalesScore" width="120">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.netSalesScore) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="销售收入加权" field="metric.netSalesScoreW" width="120">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.netSalesScoreW) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利" field="metric.grossProfit" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossProfit) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利占比" field="metric.grossProfitPct" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossProfitPct) }}%
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利得分" field="metric.grossProfitScore" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossProfitScore) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利加权" field="metric.grossProfitScoreW" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossProfitScoreW) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利率" field="metric.grossMargin" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossMargin) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利率得分" field="metric.grossMarginScore" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossMarginScore) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="毛利率加权" field="metric.grossMarginScoreW" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.grossMarginScoreW) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="到货率" field="metric.orderFillRate" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.orderFillRate) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="到货率得分" field="metric.orderFillRateScore" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.orderFillRateScore) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="到货率加权" field="metric.orderFillRateScoreW" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.orderFillRateScoreW) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="主观得分" field="metric.subjectiveScore" width="100">
                            <template slot-scope="r">
                                <!--el-input-number size="mini" :min="0" :step="0.1" 
                                    v-model="r.row.metric.subjectiveScore" /-->
                                {{ beautyNum(r.row.metric.subjectiveScore) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="主观加权" field="metric.subjectiveScoreW" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.subjectiveScoreW) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable title="总分" field="metric.totalScore" width="100">
                            <template slot-scope="r">
                                {{ beautyNum(r.row.metric.totalScore) }}
                            </template>
                        </vxe-table-column>
                    </x-table>
                </el-tab-pane>
            </el-tabs>
    </div>
</template>

<style lang="scss">
.vendor-kpi-main {
    background-color: #FFF;
    height: calc(100vh - 10px);
    overflow: auto;
    padding: 8px 8px 8px 8px;
    width: 100%;
    position: relative;

    .weight-header {
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
        width: 100%;
        text-align: center;
    }

    .el-form-item__label {
        font-size: 12px;
    }

    .el-steps.el-steps--simple {
        padding-top: 6px;
        padding-bottom: 6px;
        margin-bottom: 6px;
    }

    .el-step__title {
        font-size: 12px !important;
    }

    .vue-splitter-container {
        height: calc(100vh - 120px) !important;
    }
    .splitter-pane-resizer.vertical {
        background: url('../../assets/images/resizer_bar.png') no-repeat left center !important;
        border: none !important;

        &:hover {
            background-color: #E4E7ED;
        }
    }
}
</style>

<script>
import { doRequest, message, confirm, beautifyNumber, deepCopy, downloadFile } from '../../utils/utils'
import XTable from '../mixins/SmartVxeTable'
import splitPane from 'vue-splitpane'

export default {
    components: {
        'x-table': XTable,
        'split-pane': splitPane
    },
    data() {
        return {
            weights: [],
            originalWeights: [],
            dateRange: [],

            activeTab: 'report-config',
            tabs: [],
            tabMap: {},

            reports: [],
            reportMap: {},

            generated: [],
            generatedLoading: false,
            preview: [],
            previewLoading: false,
            previewSdate: '',
            previewEdate: '',

            reportName: '',
        }
    },
    mounted() {
        this.getWeights();
        this.getGenerated();
    },
    methods: {
        getWeights() {
            doRequest({
                url: '/v1/web/vendor/rating/config/list',
                loading: true,
            }, {
                obj: this,
                src: 'cfgs',
                dst: 'weights',
                success: res => {
                    this.originalWeights = deepCopy(this.weights);
                },
                fail: err => {
                    message('error', '获取权重配置失败')
                }
            })
        },
        updateWeight() {
            let t = 0;
            this.weights.forEach(e => {
                t += e.weight
            })
            if (t != 100) {
                message('warning', '请确认权重之和等于 100')
                return;
            }

            doRequest({
                url: '/v1/web/vendor/rating/config/update',
                method: 'POST',
                loading: true,
                data: {
                    cfgs: this.weights
                }
            }, {
                success: res => {
                    message('success', '权重更新成功');
                    this.originalWeights = deepCopy(this.weights);
                },
                fail: err => {
                    message('error', '更新权重失败, 请稍后再试');
                    this.weights = deepCopy(this.originalWeights);
                }
            })
        },
        viewReport(p) {
            let tabName = p.uuid;
            if (this.tabMap[tabName]) {
                this.activeTab = tabName;
                return;
            }

            let tab = {
                uuid: p.uuid,
                name: tabName,
                label: p.name,
                loading: true,
                sdate: p.beginDate,
                edate: p.endDate,
                i: this.tabs.length
            }
            this.tabs.push(tab);
            this.tabMap[tabName] = tab;
            this.activeTab = tabName;

            doRequest({
                url: `/v1/web/vendor/rating/report/detail/${p.uuid}`,
            }, {
                success: res => {
                    tab.report = res.report.vendorList.sort((a, b) => {
                        return b.metric.totalScore - a.metric.totalScore;
                    });
                },
                fail: err => {
                    message('error', '获取报表失败, 请稍后再试')
                },
                finally: _ => {
                    tab.loading = false;
                }
            })
        },
        beautyNum(num) {
            return beautifyNumber(num.toFixed(2));
        },
        tabClick(t) {
            let tab = this.tabMap[t.name];
            if (tab)
                this.$refs[`vdr-kpi-rpt-tbl-${tab.i}`][0].refresh();
            else {
                this.$refs['preview-report-list'].refresh();
                this.$refs['generated-report-list'].refresh();
            }
        },
        tabRemove(name) {
            let tab = this.tabMap[name], i = 0;
            this.tabs.splice(tab.i, 1);
            delete this.tabMap[name];

            if (this.activeTab == name) {
                if (this.tabs.length > 0) {
                    this.activeTab = this.tabs[0].name
                } else {
                    this.activeTab = 'report-config';
                }
            }

            this.tabs.forEach(e => {
                e.i = i;
                i += 1;
            })
        },
        getGenerated() {
            this.generatedLoading = true;
            doRequest({
                url: '/v1/web/vendor/rating/report/list'
            }, {
                obj: this,
                src: 'reportList',
                dst: 'generated',
                fail: err => {
                    message('error', '获取报告列表失败, 请稍后再试')
                },
                finally: _ => {
                    this.generatedLoading = false;
                }
            })
        },
        getPreview() {
            if (!this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的时间范围');
                return;
            }

            this.previewSdate = this.dateRange[0];
            this.previewEdate = this.dateRange[1];
            this.previewLoading = true;
            this.reportName = `${this.previewSdate}-${this.previewEdate}`
            doRequest({
                url: '/v1/web/vendor/rating/report',
                method: 'POST',
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1]
                }
            }, {
                obj: this,
                src: 'vendorList',
                dst: 'preview',
                fail: err => {
                    message('error', '生成预览报表失败, 请稍后再试')
                },
                finally: _ => {
                   this.previewLoading = false;
                }
            })
        },
        createNewReport() {
            if (!this.reportName || !this.previewSdate || !this.previewEdate) {
                message('warning', '请输入一个有效的报告名称并选择有交换的时间范围');
                return;
            }

            doRequest({
                url: '/v1/web/vendor/rating/report/add',
                method: 'POST', 
                loading: true,
                data: {
                    uuid: '',
                    name: this.reportName,
                    beginDate: this.previewSdate,
                    endDate: this.previewEdate,
                    vendorList: this.preview
                }
            }, {
                success: res => {
                    message('success', '报告提交成功');
                    this.getGenerated();
                    this.reportName = '';
                },
                fail: err => {
                    message('error', '报告提交失败, 请稍后再试');
                }
            })
        },
        deleteReport(p) {
            confirm('warning', '删除报告', `您确定要删除报告 ${p.name} 么?`, () => {
                this.generatedLoading = true;
                doRequest({
                    url: `/v1/web/vendor/rating/report/del/${p.uuid}`
                }, {
                    success: res => {
                        message('success', '报告删除成功');
                        this.getGenerated();
                    },
                    fail: err => {
                        message('error', '报告删除失败, 请稍后再试');
                        this.generatedLoading = false;
                    }
                })
            })
        },
        exportReport(p) {
            p.loading = true;
            downloadFile({
                url: `/v1/web/vendor/rating/report/export/${p.uuid}`,
                method: 'GET',
                loading: true
            }, {
                fail: err => {
                    this.$message({
                        type: 'error',
                        message: '导出报告失败, 请稍后再试'
                    });
                },
                finally: () => {
                    p.loading = false;
                }
            })
        },
        onSubjectChange(r) {
            if (r && r.metric) {
                if (isNaN(r.metric.subjectiveScore)) {
                    r.metric.subjectiveScore = 0;
                    message('error', '请输入合法的数值');
                }
                r.metric.subjectiveScoreW = r.metric.subjectiveScore * parseFloat(this.originalWeights[4].weight) / 100;
                r.metric.totalScore = r.metric.grossMarginScoreW + r.metric.grossProfitScoreW + 
                    r.metric.netSalesScoreW + r.metric.orderFillRateScoreW + r.metric.subjectiveScoreW;
            }
        }
    }
}
</script>