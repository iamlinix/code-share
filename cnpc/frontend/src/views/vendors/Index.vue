<template>
     <div class="cnpc-vendor-global-config">
         <el-dialog :visible.sync="vendorConfigAddDlg"
            :title="updatingVendorConfig ? `${vendorConfig.vendorName}` : '添加供应商配置'"
            @closed="vendorConfigDlgClosed" width="50%">
             <el-form class="vendor-config-form" size="small" label-position="left" label-width="100px" inline>
                <!--el-form-item label="供应商名称">
                    <el-input v-model="vendorConfigName" style="width: 80%" />
                </el-form-item-->
                <el-form-item label="供应商编码">
                    <el-input v-model="vendorConfig.vendorCode" :disabled="updatingVendorConfig"/>
                </el-form-item>
                 <el-divider/>
                <!--el-form-item label="上月挂账">
                    <el-input-number v-model="vendorConfig.lastMonthUnpaid" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-form-item label="上月结余">
                    <el-input-number v-model="vendorConfig.lastMonthLeft" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-divider/-->
                <el-form-item label="本年采购">
                    <el-input-number v-model="vendorConfig.yearPurchase" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-form-item label="本年销售">
                    <el-input-number v-model="vendorConfig.yearSales" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-form-item label="本年结算">
                    <el-input-number v-model="vendorConfig.yearPaid" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-form-item label="本年挂账">
                    <el-input-number v-model="vendorConfig.yearUnpaid" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-form-item label="本年结余">
                    <el-input-number v-model="vendorConfig.yearSurplus" :disabled="vendorConfig.inactive"
                        :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-divider/>
                <el-form-item label="返点(%)">
                    <el-input-number v-model="vendorConfig.rebate" :disabled="vendorConfig.inactive"
                        :min="0.0" :max="100.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
                </el-form-item>
                <el-divider/>
                <el-form-item >
                    <el-checkbox v-model="vendorConfig.inactive"
                        class="vendor-config-checkbox-frame">非活跃供应商</el-checkbox>
                </el-form-item>
                <el-form-item >
                    <el-checkbox v-model="vendorConfig.isNew"
                        :disabled="vendorConfig.inactive"
                        class="vendor-config-checkbox-frame">新供应商</el-checkbox>
                </el-form-item>
                <el-button v-if="updatingVendorConfig" style="width: 100%" type="success"
                @click="updateVendorConfig">
                更新配置
            </el-button>
                <el-button v-else style="width: 100%" type="primary"
                @click="addVendorConfig">
                添加配置
            </el-button>
            </el-form>
        </el-dialog>
        <!--p class="title">供应商配置</p-->
        <div class="vendor-outter-div">
            <el-tabs :key="vendorTabKey" type="border-card" v-model="activeTab"
                @tab-remove="materialTabRemove" stretch @tab-click="tabClick">
                <el-tab-pane name="vendorConfig" >
                    <span slot="label"><i class="el-icon-s-tools"></i>供应商配置</span>
                    <el-row :gutter="6" style="margin-bottom: 8px">
                        <el-col :span="3">
                            <el-button size="mini" style="width: 100%"
                                icon="el-icon-download" type="success"
                                @click="downloadVendorConfigCSV">导出供应商配置</el-button>
                        </el-col>
                        <el-col :span="3">
                            <el-upload accept="*" :headers="{'token': getTokenVendor()}"
                                :action="uploadUrl" :show-file-list="false" :multiple="false"
                                :before-upload="vendorConfigUploadBefore"
                                :on-success="vendorConfigUploadSuccess"
                                :on-error="vendorConfigUploadError"
                                name="upload">
                                <el-button style="width: 100%" size="mini" icon="el-icon-upload2"
                                    type="primary">导入供应商配置</el-button>
                            </el-upload>
                        </el-col>
                        <el-col :span="3">
                            <el-button size="mini" style="width: 100%"
                                icon="el-icon-plus" type="warning"
                                @click="showVendorConfigDlg">添加供应商配置</el-button>
                        </el-col>
                        <el-col :span="3">
                            <el-button size="mini" style="width: 100%"
                                icon="el-icon-delete" type="danger"
                                @click="clearVendorConfig">清空供应商配置</el-button>
                        </el-col>
                    </el-row>

                    <x-table :stripe="true" size="medium" :data="vendorConfigList" :minus="135"
                        show-overflow="tooltip" ref="vendor-table"
                        :border="true" :loading="configTableLoading" :highlight-current-row="true">
                        <vxe-table-column width="45" type="seq" fixed="left" align="center"/>
                        <vxe-table-column sortable field="vendorCode" title="供应商编码" fixed="left" width="130"/>
                        <vxe-table-column sortable field="vendorName" title="供应商名称" width="300"/>
                        <!--vxe-table-column sortable field="lastMonthUnpaid" title="上月挂账" width="100"/-->
                        <!--vxe-table-column sortable field="lastMonthLeft" title="上月结余" width="100"/-->
                        <vxe-table-column sortable field="yearPurchase" title="本年采购" width="100"/>
                        <vxe-table-column sortable field="yearSales" title="本年销售" width="100"/>
                        <vxe-table-column sortable field="yearPaid" title="本年结算" width="100"/>
                        <vxe-table-column sortable field="yearUnpaid" title="本年挂账" width="100"/>
                        <vxe-table-column sortable field="yearSurplus" title="本年结余" width="100"/>
                        <vxe-table-column sortable field="rebate" title="返点(%)" width="100"/>
                        <vxe-table-column sortable field="inactive" title="是否活跃" width="100" align="center">
                            <template slot-scope="scope">
                                <i class="el-icon-error" v-if="scope.row.inactive"
                                    style="color: #F56C6C; font-size: 20px" />
                                <i class="el-icon-success" v-else
                                    style="color: #67C23A; font-size: 20px" />
                            </template>
                        </vxe-table-column>
                        <vxe-table-column sortable field="inactive" title="新供应商" width="100" align="center">
                            <template slot-scope="scope">
                                <i class="el-icon-success" v-if="scope.row.isNew"
                                    style="color: #67C23A; font-size: 20px" />
                                <i class="el-icon-error" v-else
                                    style="color: #F56C6C; font-size: 20px" />
                            </template>
                        </vxe-table-column>
                        <vxe-table-column title="操作" width="100" align="center" fixed="right">
                            <template slot-scope="scope">
                                <el-button size="mini" circle type="primary" icon="el-icon-edit"
                                    @click="showVendorConfigDlgUpdate(scope.row)"/>
                                <el-button size="mini" circle type="warning" icon="el-icon-close"
                                    @click="deleteVendorConfig(scope.row.vendorCode, scope.row.vendorName)"/>
                            </template>
                        </vxe-table-column>
                    </x-table>
                </el-tab-pane>
                <el-tab-pane name="vendorExtras">
                    <span slot="label"><i class="el-icon-collection-tag"></i>特例</span>
                    <el-row :gutter="6" style="margin-bottom: 8px">
                        <el-col :span="3">
                            <el-button size="mini" style="width: 100%"
                                icon="el-icon-download" type="success"
                                @click="downloadExtraCSV">导出特例</el-button>
                        </el-col>
                        <el-col :span="3">
                            <el-upload accept="*" :headers="{'token': getTokenVendor()}"
                                :action="uploadExtraUrl" :show-file-list="false" :multiple="false"
                                :before-upload="vendorConfigUploadBefore"
                                :on-success="extraUploadSuccess"
                                :on-error="vendorConfigUploadError"
                                name="upload">
                                <el-button style="width: 100%" size="mini" icon="el-icon-upload2"
                                    type="primary">导入特例</el-button>
                            </el-upload>
                        </el-col>
                        <el-col :span="3">
                            <el-popover placement="bottom" trigger="click" :width="300">
                                <el-button slot="reference" size="mini" style="width: 100%"
                                    icon="el-icon-plus" type="success">添加特例</el-button>
                                <el-input style="width: 100%; margin-bottom: 6px; margin-bottom: 6px"
                                    size="mini" placeholder="供应商代码" v-model="nvc"/>
                                <el-input style="width: 100%; margin-bottom: 6px"
                                    size="mini" placeholder="供应商名称" v-model="nvn" />
                                <el-input style="width: 100%; margin-bottom: 6px"
                                    size="mini" placeholder="商品代码" v-model="nmc" />
                                <el-input style="width: 100%; margin-bottom: 6px"
                                    size="mini" placeholder="商品名称" v-model="nmt" />
                                <el-input-number style="width: 100%; margin-bottom: 6px"
                                    size="mini" placeholder="返点" v-model="nmr"
                                    :step="0.1" :min="0"/>
                                <el-button type="primary" size="mini" style="width: 100%" @click="addExtra">
                                    提交
                                </el-button>
                            </el-popover>
                        </el-col>
                        <el-col :span="3">
                            <el-button size="mini" style="width: 100%"
                                icon="el-icon-delete" type="danger"
                                @click="clearExtra">清空特例</el-button>
                        </el-col>
                    </el-row>

                    <x-table :data="extras" :border="true" :highlight-current-row="true" show-overflow="tooltip"
                        :stripe="true" :minus="135" :loading="extraTableLoading"
                        ref="extra-table">
                        <vxe-table-column type="seq" width="45" align="center"/>
                        <vxe-table-column sortable title="供应商编码" width="130" field="vendor" />
                        <vxe-table-column sortable title="供应商名称" width="300" field="vendorName" />
                        <vxe-table-column sortable title="商品编码" width="200" field="material" />
                        <vxe-table-column sortable title="商品名称" width="300" field="materialTxt" />
                        <vxe-table-column sortable title="返点" width="120" field="rebate">
                            <template slot-scope="s">
                                <span v-if="!s.row.editting">
                                    {{ s.row.rebate }}
                                </span>
                                <el-input-number v-else :step="0.1" :min="0"
                                    size="mini" style="width: 100%" v-model="s.row.t" />
                            </template>
                        </vxe-table-column>
                        <vxe-table-column title="操作" width="120">
                            <template slot-scope="s">
                                <el-button v-if="s.row.editting"
                                    circle type="success" size="mini"
                                    icon="el-icon-refresh-right"
                                    @click="doExtraEditSet(s.row, false)"/>
                                <el-button v-else
                                    circle type="success" size="mini"
                                    icon="el-icon-edit" @click="doExtraEditSet(s.row, true)"/>
                                <el-button v-if="s.row.editting" circle type="primary"
                                    size="mini" icon="el-icon-check"
                                    @click="updateExtra(s.row.vendor,
                                        s.row.vendorName, s.row.material, s.row.materialTxt,
                                        s.row.t)"/>
                                <el-button v-else circle type="danger" size="mini"
                                    icon="el-icon-close"
                                    @click="deleteExtra(s.row.material, s.row.materialTxt)"/>
                            </template>
                        </vxe-table-column>
                    </x-table>
                </el-tab-pane>
                <!--el-tab-pane name="vendors" id="vendor-tab">
                    <span slot="label">
                        <el-tooltip effect="dark" placement="top"
                            :content="actualStartDate.length > 0 ? actualStartDate + ' 至 ' + actualEndDate : '尚未选择查询时间段'">
                            <span>
                                <i class="el-icon-s-home"></i>总体供应商数据
                            </span>
                        </el-tooltip>
                    </span>
                    <el-card style="margin-bottom: 12px">
                        <div class="vendors-form">
                            <el-form :inline="true">
                                <el-form-item label="时间段">
                                    <el-date-picker
                                        v-model="selectedTimespan"
                                        type="daterange"
                                        :picker-options="pickerOptions"
                                        range-separator="-"
                                        start-placeholder="开始日期"
                                        end-placeholder="结束日期"
                                        align="right"
                                        :clearable=false>
                                    </el-date-picker>
                                </el-form-item>
                                <el-form-item>
                                    <el-button type="primary" @click="getVendorList">查询供应商数据</el-button>
                                </el-form-item>
                            </el-form>
                        </div>
                    </el-card>
                    <el-button v-if="vendors.length > 0" icon="el-icon-download"
                        type="success" size="mini" style="margin-bottom: 8px"
                        @click="downloadCSVV2()">导出总体供应商数据</el-button>
                    <s-table :stripe="true" tableSize="small" :tableData="vendors" :minus="214"
                        :border="true">
                        <el-table-column width="45" type="index" align="center"/>
                        <el-table-column sortable prop="vendorName" width="270" label="供应商名称"/>
                        <el-table-column sortable prop="vendorCode" width="120" label="供应商编码" />
                        <el-table-column sortable prop="pconfZje" label="采购金额(含税)">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.pconfZje.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="invQty" label="销售数量">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.invQty) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="costWtax" label="销售成本(含税)">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.costWtax.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="openZinvCost" label="期初金额">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.openZinvCost.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="closeZinvCost" label="期末金额">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.closeZinvCost.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column label="操作" width="75">
                            <template slot-scope="scope">
                                <el-button round
                                    @click="getVendorMaterials(scope.row.vendorCode, scope.row.vendorName)"
                                    size="mini" type="primary">
                                    详情
                                </el-button>
                            </template>
                        </el-table-column>
                    </s-table>
                </el-tab-pane-->
                <el-tab-pane v-for="(i, n) in vendorMaterialTabs" :name="i.name"
                    :key="i.name" closable>
                    <span slot="label">
                        <el-tooltip effect="dark" placement="top"
                            :content="i.tag">
                            <span>
                                <i class="el-icon-s-goods"></i>
                                {{ i.vendorCode }}
                            </span>
                        </el-tooltip>
                    </span>
                    <el-button v-if="vendors.length > 0" icon="el-icon-download" type="success" size="mini"
                        style="margin-bottom: 8px;"
                        @click="downloadCSVV2(i.vendorCode, i.start, i.end)">导出供应商详细数据</el-button>
                    <s-table :stripe="true" size="medium" :data="vendorMaterials[i.name]"
                        :minus="120"
                        :border="true" :show-summary="true"
                        :summary-method="calcVendorSummary">
                        <el-table-column width="45" type="index" align="center"/>
                        <el-table-column sortable prop="materialName" width="300" title="商品名称"/>
                        <el-table-column sortable prop="materialCode" width="165" title="商品编码" />
                        <el-table-column sortable prop="pconfZje" title="采购金额(含税)">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.pconfZje.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="invQty" title="销售数量">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.invQty) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="costWtax" title="销售成本(含税)">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.costWtax.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="openZinvCost" title="期初金额">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.openZinvCost.toFixed(2)) }}
                            </template>
                        </el-table-column>
                        <el-table-column sortable prop="closeZinvCost" title="期末金额">
                            <template slot-scope="scope">
                                {{ beautyNum(scope.row.closeZinvCost.toFixed(2)) }}
                            </template>
                        </el-table-column>
                    </s-table>
                </el-tab-pane>
            </el-tabs>
        </div>
     </div>
</template>

<style lang="scss">
.cnpc-vendor-global-config {
    background-color: white;
    max-height: calc(100vh - 10px);
    overflow: auto;
    padding: 8px 8px 8px 8px;
    width: 100%;
    position: relative;

    p.title {
        text-align: center;
        font-size: 18px;
        font-weight: bold;
        padding: 10px 0 20px 0;
    }

    .vendors-form > .el-form > .el-form-item {
        margin-bottom: 0px;
    }

    .vendor-config-form {
        .el-divider--horizontal {
            margin: 0 0 18px 0 !important;
        }
    }

    .vendor-outter-div .el-table {
        font-size: 13px;
    }

    .vendor-config-checkbox-frame {
        background-color: #D9ECFF;
        border-radius: 7px;
        padding: 4px 12px 4px 12px;
        margin-right: 20px;
    }

    .extra-vendor-row {
        background-color: #409EFF;

        td {
            background-color: #409EFF !important;
        }
    }

    .extra-inner-table {
        th.is-leaf {
            border: 1px solid #EBEEF5
        }
    }
}

</style>

<script>
import { doRequest, doRequestv2, toggleKey, beautifyNumber, deepCopy, message, confirm } from '../../utils/utils';
import moment from 'moment';
import fileDownload from 'js-file-download';
import { getToken } from '../../utils/dataStorage.js'
import { Loading } from 'element-ui'
import Config from '../../config/index'
import SmartTable from '../mixins/SmartMaxHeightTable'
import XTable from '../mixins/SmartVxeTable'
require('vue2-animate/dist/vue2-animate.min.css');

export default {
    components: {
        's-table': SmartTable,
        'x-table': XTable
    },
    data() {
        return {
            vendors: [],
            vendorTabName: '总体供应商数据',
            selectedVendor: null,
            vendorMaterials: {},
            vendorTableHeight: 0,
            vendorTableKey: 'vendorListTableKey',
            vendorConfigKey: 'vendorConfigKey',
            materialTableKey: 'materialListTableKey',
            vendorTabKey: 'vendorTabKey',
            activeTab: 'vendorConfig',
            vendorMaterialTabs: [],
            vendorMaterialKeys: {},
            vendorConfigList: [],
            configTableLoading: false,
            extraTableLoading: false,
            vendorConfigAddDlg: false,
            updatingVendorConfig: false,
            vendorConfig: {
                vendorCode: '',
                vendorName: '',
                yearPurchase: 0.0,
                yearSales: 0.0,
                yearPaid: 0.0,
                yearUnpaid: 0.0,
                yearSurplus: 0.0,
                rebate: 0.0,
                inactive: false,
                isNew: false,
            },
            localLoading: null,
            actualStartDate: '',
            actualEndDate: '',
            pickerOptions: {
                shortcuts: [{
                    text: '最近7天',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '最近30天',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();

                        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '一月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(0, 1);
                        end.setMonth(0, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '二月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(1, 1);
                        end.setMonth(2, 1);
                        end.setTime(end.getTime() - 3600 * 1000 * 24);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '三月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(2, 1);
                        end.setMonth(2, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '四月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(3, 1);
                        end.setMonth(3, 30);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '五月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(4, 1);
                        end.setMonth(4, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '六月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(5, 1);
                        end.setMonth(5, 30);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '七月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(6, 1);
                        end.setMonth(6, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '八月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(7, 1);
                        end.setMonth(7, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '九月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(8, 1);
                        end.setMonth(8, 30);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '十月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(9, 1);
                        end.setMonth(9, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '十一月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(10, 1);
                        end.setMonth(10, 30);
                        picker.$emit('pick', [start, end]);
                    }
                }, {
                    text: '十二月',
                    onClick(picker) {
                        const end = new Date();
                        const start = new Date();
                        start.setMonth(11, 1);
                        end.setMonth(11, 31);
                        picker.$emit('pick', [start, end]);
                    }
                }]
            },
            selectedTimespan: [],

            extras: [],
            nvc: '',
            nvn: '',
            nmc: '',
            nmt: '',
            nmr: 0,
            selectedExtraVendor: {}
        }
    },
    methods: {
        getTokenVendor() {
            return getToken();
        },
        getVendorList() {
            if (this.selectedTimespan.length > 0) {
                doRequest({
                    loading: true,
                    url: '/v1/web/erp/vendors',
                    method: 'post',
                    data: {
                        beginDate: moment(this.selectedTimespan[0]).format('YYYY-MM-DD'),
                        endDate: moment(this.selectedTimespan[1]).format('YYYY-MM-DD')
                    }
                }, {
                    obj: this,
                    src: 'vendorsList',
                    dst: 'vendors',
                    success: _ => {
                        this.actualStartDate = moment(this.selectedTimespan[0]).format('YYYY-MM-DD');
                        this.actualEndDate = moment(this.selectedTimespan[1]).format('YYYY-MM-DD');
                        //this.vendorTabName = `总体供应商数据 (${this.actualStartDate}-${this.actualEndDate})`;
                        this.activeTab = 'vendors';
                    }
                })
            } else {
                this.$message({
                    type: "warning",
                    message: "请先选择一个时间段"
                })
            }
        },
        getVendorMaterials(vendorCode, vendorName) {
            let tabName = vendorCode + this.actualStartDate + this.actualEndDate;
            if (tabName in this.vendorMaterialKeys) {
                this.activeTab = tabName;
                return;
            }

            doRequest({
                loading: true,
                url: '/v1/web/erp/materials',
                method: 'post',
                data: {
                    beginDate: this.actualStartDate,
                    endDate: this.actualEndDate,
                    vendor: vendorCode
                }
            }, {
                success: res => {
                    this.vendorMaterials[tabName] = res.materialsList;
                    this.vendorMaterialKeys[tabName] = true;
                    this.vendorMaterialTabs.push({
                        label: vendorName,
                        name: tabName,
                        vendorCode: vendorCode,
                        tag: `${vendorName}, ${this.actualStartDate} 至 ${this.actualEndDate}`,
                        start: this.actualStartDate,
                        end: this.actualEndDate,
                    });
                    this.activeTab = tabName;
                },
                fail: err => {
                    console.log('get material error:', err);
                }
            })
        },
        vendorListDbClick(r, c, e) {
            if (!(r.vendorCode in this.vendorMaterialKeys)) {
                this.getVendorMaterials(r.vendorCode, r.vendorName);
            }
        },
        tabClick(tab) {
            if (tab.name == 'vendorConfig') {
                this.$refs['vendor-table'].refresh();
            } else if (tab.name == 'vendorExtras') {
                this.$refs['extra-table'].refresh();
            }
        },
        materialTabRemove(name) {
            if (name in this.vendorMaterialKeys) {
                delete this.vendorMaterialKeys[name];
                delete this.vendorMaterials[name];
            }

            for (let i = 0; i < this.vendorMaterialTabs.length; i ++) {
                if (this.vendorMaterialTabs[i].name === name) {
                    this.vendorMaterialTabs.splice(i, 1);
                    break;
                }
            }

            this.activeTab = 'vendors';

            return true;
        },
        downloadCSV(csvName, names, labels, data) {
            let csv = '';
            labels.forEach((a, i) => {
                csv += `="${a}",`;
            })
            csv += '\n';
            data.forEach((a, i) => {
                names.forEach((b, j) => {
                    if (j < 2)
                        csv += `="${a[b]}",`;
                    else
                        csv += `"${a[b]}",`;

                });
                csv += '\n';
            });

            let aLink = document.createElement('a');
            aLink.download = `${csvName}-${moment(this.selectedTimespan[0]).format('YYYYMMDD')}-${moment(this.selectedTimespan[1]).format('YYYYMMDD')}.csv`;
            aLink.href = 'data:text/csv;charset=UTF-8,' + encodeURIComponent(csv);

            var event = new MouseEvent('click');
            aLink.dispatchEvent(event);
        },
        downloadCSVV2(vendorCode, start, end) {
            if (vendorCode) {
                // material request
                doRequestv2({
                    url: '/v1/web/erp/materials/export',
                    method: 'post',
                    responseType: 'blob',
                    loading: true,
                    data: {
                        beginDate: start,
                        endDate: end,
                        vendor: vendorCode
                    }
                }, {
                    success: res => {
                        let fn = res.headers['content-disposition'].substring('attachment;filename='.length);
                        if (fn && res.data) {
                            fileDownload(res.data, fn);
                        }
                    },
                    fail: err => {
                        console.log(err);
                    }
                })
            } else {
                // vendor request
                doRequestv2({
                    url: '/v1/web/erp/vendors/export',
                    method: 'post',
                    responseType: 'blob',
                    loading: true,
                    data: {
                        beginDate: this.actualStartDate,
                        endDate: this.actualEndDate,
                    }
                }, {
                    success: res => {
                        let fn = res.headers['content-disposition'].substring('attachment;filename='.length);
                        if (fn && res.data) {
                            fileDownload(res.data, fn);
                        }
                    },
                    fail: err => {
                        console.log(err);
                    }
                })
            }
        },
        getVendorConfigList() {
            this.configTableLoading = true;
            doRequest({
                url: '/v1/web/erp/vendor/config/list',
                method: 'get'
            }, {
                obj: this,
                src: 'cfgs',
                dst: 'vendorConfigList',
                finally: _ => {
                    this.configTableLoading = false;
                }
            })
        },
        showVendorConfigDlg(updating) {
            this.updatingVendorConfig = false;
            this.vendorConfigAddDlg = true;
        },
        showVendorConfigDlgUpdate(config) {
            this.updatingVendorConfig = true;
            this.vendorConfig = deepCopy(config);
            this.vendorConfigAddDlg = true;
        },
        addVendorConfig() {
            // if (this.vendorConfigName.length <= 0 || this.vendorConfigCode.length <= 0) {
            if (this.vendorConfig.vendorCode.length <= 0) {
                message('warning', '请输入有效的供应商信息');
                return;
            }

            doRequest({
                url: '/v1/web/erp/vendor/config/add',
                method: 'post',
                loading: true,
                data: this.vendorConfig
            }, {
                success: _ => {
                    this.$message({
                        type: 'success',
                        message: '添加供应商配置成功'
                    });
                    this.getVendorConfigList();
                    this.vendorConfigAddDlg = false;
                },
                fail: err => {
                    console.log(err);
                    this.$message({
                        type: 'error',
                        message: '添加供应商配置失败，请稍后再试'
                    });
                }
            })
        },
        clearVendorConfig() {
            this.$confirm(`您确定要删除所有供应商配置么？`, '清空供应商配置', {
                type: 'warning'
            }).then(() => {
                doRequest({
                    url: '/v1/web/erp/vendor/config/clear',
                    method: 'get'
                }, {
                    success: res => {
                        this.$message({
                            type: 'success',
                            message: '清空供应商配置成功！'
                        });
                        this.getVendorConfigList();
                    },
                    fail: err => {
                        console.log(err);
                        this.$message({
                            type: 'error',
                            message: '清空供应商配置失败，请稍后再试！'
                        })
                    }
                })
            }).catch(() => {

            })
        },
        updateVendorConfig() {
            if (this.vendorConfig.vendorCode.length <= 0) {
                this.$message({
                    type: 'warning',
                    message: '请输入有效的供应商信息'
                });
                return;
            }

            doRequest({
                url: '/v1/web/erp/vendor/config/update',
                method: 'post',
                loading: true,
                data: this.vendorConfig,
            }, {
                success: _ => {
                    this.$message({
                        type: 'success',
                        message: '更新供应商配置成功'
                    });
                    this.getVendorConfigList();
                    this.vendorConfigAddDlg = false;
                },
                fail: err => {
                    console.log(err);
                    this.$message({
                        type: 'error',
                        message: '更新供应商配置失败，请稍后再试'
                    });
                }
            })
        },
        deleteVendorConfig(vendorCode, vendorName) {
            this.$confirm(`您确定要删除供应商 ${vendorName} 的配置么？`, '删除供应商配置', {
                type: 'warning'
            }).then(() => {
                doRequest({
                    url: `/v1/web/erp/vendor/config/del/${vendorCode}`,
                    method: 'get'
                }, {
                    success: res => {
                        this.$message({
                            type: 'success',
                            message: '删除供应商配置成功！'
                        });
                        this.getVendorConfigList();
                    },
                    fail: err => {
                        console.log(err);
                        this.$message({
                            type: 'error',
                            message: '删除供应商配置失败，请稍后再试！'
                        })
                    }
                })
            }).catch(() => {

            })
        },
        downloadVendorConfigCSV() {
            doRequestv2({
                url: '/v1/web/erp/vendor/config/export',
                method: 'get',
                responseType: 'blob',
                loading: true
            }, {
                success: res => {
                    let fn = res.headers['content-disposition'].substring('attachment;filename='.length);
                    if (fn && res.data) {
                        fileDownload(res.data, fn);
                    }
                },
                fail: err => {
                    message('error', '下载供应商配置失败')
                    console.log(err);
                }
            })
        },
        vendorConfigUploadBefore(file) {
            if (!this.localLoading) {
                this.localLoading = Loading.service({
                    text: `${file.name} 上传中...`
                })
            }
        },
        vendorConfigUploadSuccess(res, file, fileList) {
            if (this.localLoading) {
                this.localLoading.close();
                this.localLoading = null;
            }

            this.$message({
                type: "success",
                message: `${file.name} 上传成功`
            });
            this.getVendorConfigList();
        },
        extraUploadSuccess(res, file, fileList) {
            if (this.localLoading) {
                this.localLoading.close();
                this.localLoading = null;
            }

            this.$message({
                type: "success",
                message: `${file.name} 上传成功`
            });
            this.getExtras();
        },
        vendorConfigUploadError(err, file, fileList) {
            if (this.localLoading) {
                this.localLoading.close();
                this.localLoading = null;
            }

            this.$message({
                type: "error",
                message: `${file.name} 上传失败，请稍后再试`
            });
        },
        beautyNum(num) {
            return beautifyNumber(num);
        },
        vendorConfigDlgClosed() {
            this.vendorConfig.vendorCode = '';
            this.vendorConfig.rebate = 0.0;
            this.vendorConfig.inactive = false;
            this.vendorConfig.isNew = false;
            this.vendorConfig.lastMonthUnpaid = 0.0;
            this.vendorConfig.lastMonthLeft = 0.0;
            this.vendorConfig.yearPurchase = 0.0;
            this.vendorConfig.yearSales = 0.0;
            this.vendorConfig.yearPaid = 0.0;
            this.vendorConfig.yearUnpaid = 0.0;
        },
        calcVendorSummary(param) {
            const { columns, data } = param;
            const sums = [];
            columns.forEach((column, index) => {
                if (index <= 1) {
                    sums[index] = '';
                    return;
                }

                if (index === 2) {
                    sums[index] = '合计';
                    return;
                }

                const values = data.map(item => Number(item[column.property]));
                if (!values.every(value => isNaN(value))) {
                    sums[index] = values.reduce((prev, curr) => {
                        const value = Number(curr);
                        if (!isNaN(value)) {
                            return prev + curr;
                        } else {
                            return prev;
                        }
                    }, 0);
                    sums[index] = this.beautyNum(sums[index].toFixed(2));
                }
            });

            return sums;
        },
        getExtras() {
            this.extraTableLoading = true;
            doRequest({
                url: '/v1/web/erp/vendor/ex-matl/config/list',
                method: 'GET'
            }, {
                obj: this,
                src: 'cfgs',
                dst: 'extras',
                finally: _ => {
                    this.extraTableLoading = false;
                }
            })
        },
        resetExtras() {

            this.extras.forEach(e => {
                e.editting = false;
            })

        },
        doExtraEditSet(row, b) {
            this.$set(row, 'editting', b);
        },
        addExtra() {
            if (this.nvc.length <= 0 || this.nvn.length <= 0 || this.nmc.length <= 0 || this.nmt.length <= 0) {
                message('error', '请输入有效的供应商与商品信息');
                return;
            }
            doRequest({
                url :'/v1/web/erp/vendor/ex-matl/config/add',
                method: 'POST',
                loading: true,
                data: {
                    vendor: this.nvc,
                    vendorName: this.nvn,
                    material: this.nmc,
                    materialTxt: this.nmt,
                    rebate: this.nmr
                }
            }, {
                success: res => {
                    message('success', '添加特例成功');
                    this.getExtras();
                    this.nvc = '';
                    this.nvn = '';
                    this.nmc = '';
                    this.nmt = '';
                    this.nmr = 0;
                },
                fail: err => {
                    message('error', '添加特例失败');
                    console.log(err);
                }
            })
        },
        deleteExtra(material, materialTxt) {
            confirm('warning', '删除特例', `确定要删除 ${materialTxt} 的特例配置么?`, () => {
                doRequest({
                    url: `/v1/web/erp/vendor/ex-matl/config/del/${material}`,
                    method: 'GET',
                    loading: true
                }, {
                    success: res => {
                        this.getExtras();
                    }
                })
            })
        },
        updateExtra(vendor, vendorName, material, materialTxt, rebate) {
            doRequest({
                url: '/v1/web/erp/vendor/ex-matl/config/update',
                method: 'POST',
                loading: true,
                data: {
                    material: material,
                    materialTxt: materialTxt,
                    vendor: vendor,
                    vendorName: vendorName,
                    rebate: rebate
                }
            }, {
                success: res => {
                    this.getExtras()
                },
                fail: err => {
                    message('error', '更新特例失败')
                }
            })

        },
        downloadExtraCSV() {
            doRequestv2({
                url: '/v1/web/erp/vendor/ex-matl/config/export',
                method: 'get',
                responseType: 'blob',
                loading: true
            }, {
                success: res => {
                    let fn = res.headers['content-disposition'].substring('attachment;filename='.length);
                    if (fn && res.data) {
                        fileDownload(res.data, fn);
                    }
                },
                fail: err => {
                    message('error', '下载供应商配置失败')
                    console.log(err);
                }
            })
        },
        clearExtra() {
            confirm('warning', '清空特例', '确定要清空所有特例么?', () => {
                doRequest({
                    url: '/v1/web/erp/vendor/ex-matl/config/clear',
                    method: 'GET',
                    loading: true
                }, {
                    success: res => {
                        message('success', '清空特例成功')
                        this.extras = [];
                    },
                    fail: err => {
                        console.log(err);
                        message('error', '清空特例失败')
                    }
                })
            })
        },
        onResize() {
            // this.vendorTableHeight = window.innerHeight - 255;
            // toggleKey(this, 'vendorTabKey');
            //toggleKey(this, 'vendorConfigKey');
        }
    },
    mounted: function() {
        // this.vendorTableHeight = window.innerHeight - 255;
        // toggleKey(this, 'vendorTabKey');
        // //toggleKey(this, 'vendorConfigKey');
        // this.$nextTick(() => {
        //     window.addEventListener('resize', this.onResize);
        // })
        this.getExtras();
        this.getVendorConfigList();
    },
    computed: {
        uploadUrl: function() {
            return Config.apiUrl + '/v1/web/erp/vendor/config/import';
        },
        uploadExtraUrl: function() {
            return Config.apiUrl + '/v1/web/erp/vendor/ex-matl/config/import';
        }
    },
}
</script>
