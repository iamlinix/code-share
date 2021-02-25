<template>
    <div class="cnpc-no-sales-div">
        <el-dialog :visible.sync="configStationDlgShow" :close-on-click-modal="false" 
            append-to-body width="70%" @closed="stationSearch = ''">
            <p slot="title" style="text-align: center">
                为 “{{ materialToMonitor.name }}” 选择对应油站
            </p>
            <el-row :gutter="8">
                <el-col :span="10">
                <el-input size="small" style="width: 100%; margin-bottom: 10px" 
                    placeholder="输入代码或名称过滤"
                    prefix-icon="el-icon-search"
                    v-model="stationSearch"
                    clearable
                    @input="stationFilter"/>
                </el-col>
                <el-col :span="8" :offset="4">
                    <p style="font-weight: bold">
                    已选油站({{ selectedStations.length }}):
                    </p>
                </el-col>
            </el-row>
            <el-row :gutter="8">
                <el-col :span="10">
                    <x-table :data="filteredStations" :minus="300" size="mini" show-overflow="tooltip"
                        ref='station-selection' @checkbox-change="stationSelect" @checkbox-all="stationSelectAll"
                        border :default-sort="{
                            prop: 'orgCode'
                        }">
                        <vxe-table-column type="checkbox" width="35"/>
                        <vxe-table-column title="油站代码" width="100" field="orgCode" sortable/>
                        <vxe-table-column title="油站名称" field="orgText" sortable>
                            <span slot="header">
                                油站名称{{ `(${stationCheckBoxes.count}/${stationsShadow.length})` }}
                            </span>
                        </vxe-table-column>
                    </x-table>
                </el-col>
                <el-col :span="4" style="text-align: center; padding-top: 15%">
                    <br/>
                    <el-button round type="primary" @click="stationMoveRight">&gt;&gt;</el-button>
                    <br/>
                    <br/>
                    <el-button round type="primary" @click="stationMoveLeft">&lt;&lt;</el-button>
                </el-col>
                <el-col :span="10">
                    <x-table :data="selectedStations" :minus="300" border size="mini" show-overflow="tooltip"
                        @checkbox-change="stationSelectR" @checkbox-all="stationSelectR" :default-sort="{
                            prop: 'orgCode'
                        }">
                        <vxe-table-column type="checkbox" width="35"/>
                        <vxe-table-column title="油站代码" width="100" field="orgCode" sortable/>
                        <vxe-table-column title="油站名称" field="orgText" sortable/>
                    </x-table>
                </el-col>
            </el-row>
            <el-button v-if="updatingMonitor" style="width: 100%; margin-top: 10px" 
                type="primary" @click="updateMonitor">更新商品关注</el-button>
            <el-button v-else style="width: 100%; margin-top: 10px" 
                type="success" @click="submitMonitor">创建商品关注</el-button>
            
            <!--el-transfer :data="stations" filterable v-model="selectedStations"
                :props="{
                    key: 'orgCode',
                    label: 'orgText'
                }"
                :titles="['油站列表', '已选油站']"
                style="height: 500px">
            </el-transfer-->
        </el-dialog>
        <el-dialog :visible.sync="configDlgShow" width="55%" :close-on-click-modal="false" 
            @closed="materialCodeToMonitor = ''">
            <p slot="title" style="text-align: center">
                关注商品列表
            </p>
            <el-input size="small" style="width: 60vh" v-model="materialCodeToMonitor" 
                placeholder="请输入商品编码"/>
            <el-button size="small" type="primary" style="margin-left: 10px"
                @click="monitorMaterial" icon="el-icon-plus">添加关注商品</el-button>
            <!--s-table style="margin-top: 10px" border :data="monitorsShadow" :minus="300"
                size="small" :span-method="monitorSpan"
                @cell-mouse-enter="materialColEnter" @cell-mouse-leave="showMaterialHint = false">
                <el-table-column :index="0" label="商品编码" prop="material" width="150"/>
                <el-table-column :index="1" label="商品名称" prop="materialName" width="260">
                    <template slot-scope="r">
                        <p>
                            {{ r.row.materialName }}
                            <el-tooltip placement="top" effect="dark" content="删除关注商品" :hide-after="1000">
                                <i class="el-icon-error" 
                                    style="font-size: 18px; float: right; cursor: pointer; color: red"
                                    @click="clearPlantsForMaterial(r.row)"/>
                            </el-tooltip>
                            <el-tooltip placement="top" effect="dark" content="添加油站" :hide-after="1000">
                                <i class="el-icon-circle-plus" 
                                    style="font-size: 18px; float: right; cursor: pointer; color: #67C23A"
                                    @click="addPlantForMaterial(r.row)"/>
                            </el-tooltip>
                        </p>
                    </template>
                </el-table-column>
                <el-table-column label="关注油站" prop="orgText">
                    <template slot-scope="r">
                        <p>
                            {{ r.row.orgText }}
                            <el-tooltip placement="top" effect="dark" content="删除油站" :hide-after="1000">
                                <i class="el-icon-remove" 
                                    style="font-size: 18px; float: right; cursor: pointer; color: #F56C6C"
                                    @click="removePlantFromMaterial(r.row)"/>
                            </el-tooltip>
                        </p>
                    </template>
                </el-table-column>
            </s-table-->
            <x-table style="margin-top: 10px" border :data="monitorsShadow" :minus="300"
                size="small" show-overflow="tooltip" v-loading="monitorLoading">
                <vxe-table-column type="expand" width="45">
                    <template v-slot:content="{ row, rowIndex }">
                        <x-table style="border: 1px solid #EBEEF5" :data="row.plants" size="small" 
                            @checkbox-all="phaseOneStationSelectChange"
                            @checkbox-change="phaseOneStationSelectChange" show-overflow="tooltip"
                            :minus="600" border="inner">
                            <vxe-table-column type="checkbox" width="35"/>
                            <vxe-table-column type="seq" width="45"/>
                            <vxe-table-column title="油站编码" field="plant" width="100" sortable/>
                            <vxe-table-column title="油站名称" field="plantName" sortable>
                                <template slot-scope="s">
                                    <p>
                                        {{ s.row.plantName }}
                                        <el-tooltip placement="right" effect="dark" content="删除油站" :hide-after="1000">
                                            <i class="el-icon-remove" 
                                                style="font-size: 18px; float: right; cursor: pointer; color: #F56C6C"
                                                @click="removePlantFromMaterial(s.row, row)"/>
                                        </el-tooltip>
                                    </p>
                                </template>
                            </vxe-table-column>
                        </x-table>
                    </template>
                </vxe-table-column>
                <vxe-table-column :index="0" title="商品编码" field="material" width="150" sortable/>
                <vxe-table-column :index="1" title="商品名称" field="materialName" sortable>
                    <template slot-scope="r">
                        <p>
                            {{ r.row.materialName }}
                            <el-tooltip placement="right" effect="dark" content="删除关注商品" :hide-after="1000">
                                <i class="el-icon-error" 
                                    style="font-size: 18px; float: right; cursor: pointer; color: red"
                                    @click="clearPlantsForMaterial(r.row)"/>
                            </el-tooltip>
                            <el-tooltip placement="left" effect="dark" content="编辑油站列表" :hide-after="1000">
                                <i class="el-icon-circle-plus" 
                                    style="font-size: 18px; float: right; cursor: pointer; color: #67C23A"
                                    @click="addPlantForMaterial(r.row)"/>
                            </el-tooltip>
                        </p>
                    </template>
                </vxe-table-column>
                <vxe-table-column title="关注油站数量" width="150">
                    <template slot-scope="r">
                        {{ r.row.plants.length }}
                    </template>
                </vxe-table-column>
            </x-table>
        </el-dialog>
        <el-card>
            <span style="font-size: 14px">统计时间段 </span>
            <el-date-picker size="small"
                v-model="dateRange"
                type="daterange"
                range-separator="-"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="yyyy-MM-dd">
            </el-date-picker>
            <el-button size="small" style="margin-left: 12px" type="primary"
                @click="getNoSalesList">查询销售情况</el-button>
            <el-button size="small" style="margin-left: 12px; float: right" type="success"
                @click="configDlgShow = true"
                :disabled="stationLoading || monitorLoading"
                v-loading="stationLoading || monitorLoading">关注商品列表</el-button>
        </el-card>
        <!--s-table style="margin-top: 10px" :data="noSalesList" border stripe 
            :span-method="monitorSpan" :minus="115" size="small"
            @cell-mouse-enter="materialColEnter" @cell-mouse-leave="showMaterialHint = false">
            <el-table-column :index="0" label="关注商品编码" prop="material" width="150"/>
            <el-table-column :index="1" label="关注商品名称" prop="materialName" width="300"/>
            <el-table-column label="无销售油站名称" prop="orgText" />
        </s-table-->
        <x-table style="margin-top: 10px" :data="noSalesList" border stripe 
            :minus="115" size="small" show-overflow="tooltip">
            <vxe-table-column type="expand" width="45">
                <template v-slot:content="{ row, rowIndex }" align="center">
                    <el-button icon="el-icon-plus" size="small" v-if="!row.showGoodOnes" type="text" @click="row.showGoodOnes = true">显示有销售的站点</el-button>
                    <el-button icon="el-icon-minus" size="small" v-else type="text" @click="row.showGoodOnes = false">隐藏有销售的站点</el-button>
                    <transition name="zoom">
                        <div :key="row.material" v-if="row.showGoodOnes" style="animation-duration: 0.2s">
                            <x-table style="border: 1px solid #EBEEF5" :data="row.goodOnes" size="small" 
                                :border="false" :row-style="{color: '#67C23A'}" show-overflow="tooltip"
                                :minus="400">
                                <vxe-table-column type="seq" width="45"/>
                                <vxe-table-column title="油站编码" field="plant" width="100" sortable/>
                                <vxe-table-column title="油站名称" field="plantName" sortable/>
                            </x-table>
                            <el-divider ></el-divider>
                        </div>
                    </transition>
                    <x-table style="border: 1px solid #EBEEF5" :data="row.plants" size="small" 
                        :border="false" :row-style="{color: '#F56C6C'}"  show-overflow="tooltip" 
                        :minus="400">
                        <vxe-table-column type="seq" width="45"/>
                        <vxe-table-column title="油站编码" field="plant" width="100" sortable/>
                        <vxe-table-column title="油站名称" field="plantName" sortable/>
                    </x-table>
                </template>
            </vxe-table-column>
            <vxe-table-column title="关注商品编码" field="material" width="150" sortable/>
            <vxe-table-column title="关注商品名称" field="materialName" sortable/>
            <vxe-table-column title="无销售站点数" width="150">
                <template slot-scope="r">
                    {{ `${r.row.plants.length} / ${monitorNumber[r.row.material]}` }}
                </template>
            </vxe-table-column>
            <vxe-table-column title="操作" width="80">
                <template slot-scope="r">
                    <el-button type="primary" size="mini" round @click="exportPlants(r.row)">导出</el-button>
                </template>
            </vxe-table-column>
        </x-table>
        <span v-show="showMaterialHint" id="flying-span">
            {{ materialHint }}
        </span>
    </div>
</template>

<style lang="scss">
.cnpc-no-sales-div {
    width: 100%;
    max-height: calc(100vh - 12px);
    background-color: white;
    padding: 8px 8px 8px 8px;
    overflow: auto;

    #flying-span {
        background-color: #303133;
        color: white;
        padding: 8px 12px 8px 12px;
        position: absolute;
        font-size: 13px;
        z-index: 9999;
    }

    p.no-sale-station {
        padding: 4px 4px 4px 4px;
        background-color: #DCDFE6;
        margin-bottom: 4px;
    }
}
</style>

<script>
import { doRequest, deepCopy, message, confirm } from '../../utils/utils'
import SmartTable from '../mixins/SmartMaxHeightTable';
import XTable from '../mixins/SmartVxeTable';
require('vue2-animate/dist/vue2-animate.min.css');
import { saveAs } from 'file-saver';

export default {
    components: {
        'x-table': XTable
    },
    mounted() {
        this.getStations();
        this.getMonitors();
    },
    data() {
        return {
            dateRange: [],
            configDlgShow: false,
            configStationDlgShow: false,

            stations: [],
            stationsShadow: [],
            filteredStations: [],
            selectedStations: [],
            stationToAdd: '',
            stationSearch: '',
            stationCheckBoxes: {
                count: 0
            },
            stationSelR: [],
            stationLoading: false,

            materials: [],
            filteredMaterials: [],

            monitors: [],
            monitorNumber: {},
            monitorsShadow: [],
            pmmcList: [],
            updatingMonitor: false,
            monitorLoading: false,

            materialCodeToMonitor: '',
            materialToMonitor: {},

            noSalesList: [],
            materialHint: '',
            showMaterialHint: false,
            hintTimeout: null,

            phaseOneStationSel: [],
        }
    },
    methods: {
        exportPlants(row) {
            let blob = '油站编码,油站名称,\n'
            row.plants.forEach(e => {
                blob += `${e.plant},${e.plantName},\n`
            })
            saveAs(new Blob([blob], {type: 'text/plain;charset=utf-8'}), `未销售-${row.materialName}.csv`)
        },
        alignNum(num) {
            let s = num + '';
            while (num.length < 3)
                s = " " + s;
            return s;
        },
        materialColEnter(row, column, cell, event) {
            if (column.index === 0 || column.index === 1) {
                if (this.hintTimeout) {
                    clearTimeout(this.hintTimeout);
                    this.hintTimeout = null;
                }
                this.materialHint = `${row.material} (${row.materialName})`;
                let e = document.getElementById('flying-span');
                e.style.left = event.x + 'px';
                e.style.top = event.y + 'px';
                this.showMaterialHint = true;
                this.hintTimeout = setTimeout(() => {
                    this.showMaterialHint = false;
                }, 2000);
            }
        },
        getStations() {
            this.stationLoading = true;
            doRequest({
                url: '/v1/web/org/plant',
                method: 'GET'
            }, {
                obj: this,
                src: 'orgList',
                dst: 'stations',
                finally: () => {
                    this.stationLoading = false;
                }
            })
        },
        getMonitors() {
            this.monitorLoading = true;
            doRequest({
                url: '/v1/web/plant/material/monitor/config/list',
                method: 'GET'
            }, {
                obj: this,
                src: 'pmmcList',
                dst: 'monitors',
                success: res => {
                    // let i = 0, j = 0, k = 0, l = 0, m = 0, v = null, s = null, monitors = [];
                    // for (i = 0; i < res.pmmcList.length; i ++) {
                    //     v = res.pmmcList[i];
                    //     k = 0;
                    //     m = l;
                    //     for (j = 0; j < v.plants.length; j ++) {
                    //         s = v.plants[j];
                    //         monitors.push({
                    //             material: v.material,
                    //             materialName: v.materialName,
                    //             orgCode: s.plant,
                    //             orgText: s.plantName,
                    //             status: s.status,
                    //             k: 0,
                    //         });
                    //         k ++;
                    //         l ++;
                    //     }
                    //     monitors[m].k = k;
                    //     monitors[m].plants = v.plants;
                    // }
                    // this.monitors = monitors;
                    if (this.monitors) {
                        let self = this;
                        this.monitors.forEach(e => {
                            self.monitorNumber[e.material] = e.plants.length;
                        })
                    } else {
                        this.monitors = []
                    }
                    
                    this.monitorsShadow = this.monitors;
                    //this.pmmcList = res.pmmcList;
                },
                finally: () => {
                    this.monitorLoading = false;
                }
            })
        },
        stationFilter(c) {
            let d = c.toLowerCase();
            this.filteredStations = this.stationsShadow.filter(v => {
                return v.orgCode.toLowerCase().includes(d) || v.orgText.toLowerCase().includes(d);
            });
            this.$nextTick(() => {
                this.$refs['station-selection'].setCheckboxRow(this.selectedStations, true);
            })
        },
        stationSelectAll({ records }) {
            let e = null;
            for (let i = 0; i < this.filteredStations.length; i ++) {
                e = this.filteredStations[i];
                if (records.includes(e)) {
                    if (!this.stationCheckBoxes[e.orgCode]) {
                        this.stationCheckBoxes[e.orgCode] = e;
                        this.stationCheckBoxes.count += 1;
                    }
                } else {
                    if (this.stationCheckBoxes[e.orgCode]) {
                        delete this.stationCheckBoxes[e.orgCode];
                        this.stationCheckBoxes.count -= 1;
                    }
                }
            }
        },
        stationSelect({ records, row }) {
            let b = false;
            for (let i = 0; i < records.length; i ++) {
                if (records[i] == row) {
                    b = true;
                    break;
                }
            }

            if (!b) {
                if (this.stationCheckBoxes[row.orgCode]) {
                    delete this.stationCheckBoxes[row.orgCode];
                    this.stationCheckBoxes.count -= 1;
                }
            } else {
                this.stationCheckBoxes[row.orgCode] = row;
                this.stationCheckBoxes.count += 1;
            }
        },
        stationSelectR({ records }) {
            this.stationSelR = records;
        },
        stationMoveRight() {
            for(let k in this.stationCheckBoxes) {
                if (k == 'count')
                    continue;
                this.selectedStations.push(this.stationCheckBoxes[k]);
                for (let i = 0; i < this.stationsShadow.length; i ++) {
                    if (this.stationsShadow[i].orgCode == this.stationCheckBoxes[k].orgCode) {
                        this.stationsShadow.splice(i, 1);
                        break;
                    }
                }
            }
            this.stationCheckBoxes = {count : 0};
            this.stationFilter(this.stationSearch);
        },
        stationMoveLeft() {
            for (let i = 0; i < this.stationSelR.length; i ++) {
                this.filteredStations.push(this.stationSelR[i]);
                for (let j = 0; j < this.selectedStations.length; j ++) {
                    if (this.selectedStations[j] == this.stationSelR[i]) {
                        this.selectedStations.splice(j, 1);
                    }
                }
            }
            this.stationSelR = [];
        },
        monitorSpan({ row, column, rowIndex, columnIndex }) {
            if (columnIndex <= 1) {
                if (row.k === 0) {
                    return {
                        rowspan: 0,
                        colspan: 0
                    }
                } else {
                    return {
                        rowspan: row.k,
                        colspan: 1,
                    }
                }
            }
        },
        monitorMaterial() {
            if (this.materialCodeToMonitor.length === 0) {
                // message({
                //     type: 'warning',
                //     message: '请输入需要添加的商品编码',
                //     showClose: true,
                // });
                message('warning', '请输入需要添加的商品编码');
                return;
            }

            for (let i = 0; i < this.monitorsShadow.length; i ++) {
                if (this.monitorsShadow[i].material.endsWith(this.materialCodeToMonitor)) {
                    message('info', `${this.materialCodeToMonitor} 已在关注列表中`);
                    return;
                }
            }

            doRequest({
                url: `/v1/web/basic/material/${this.materialCodeToMonitor.replace(/^0+/, '')}`,
                method: 'GET',
                loading: true
            }, {
                obj: this,
                src: 'material',
                dst: 'materialToMonitor',
                success: res => {
                    this.stationsShadow = deepCopy(this.stations);
                    this.filteredStations = this.stationsShadow;
                    this.selectedStations = [];
                    this.stationCheckBoxes = {count: 0};
                    this.updatingMonitor = false;
                    this.configStationDlgShow = true;
                },
                fail: err => {
                    message('error', `获取 ${this.materialCodeToMonitor} 商品信息失败`)
                }
            })
        },
        addPlantForMaterial(r) {
            this.stationsShadow = deepCopy(this.stations);
            this.stationsShadow = this.stationsShadow.filter(v => {
                let b = true;
                for (let i = 0; i < r.plants.length; i ++) {
                    if (v.orgCode == r.plants[i].plant) {
                        b = false;
                        break;
                    }
                }
                return b;
            });

            this.filteredStations = this.stationsShadow;
            this.selectedStations = [];
            for (let i = 0; i < r.plants.length; i ++) {
                this.selectedStations.push({
                    orgCode: r.plants[i].plant,
                    orgText: r.plants[i].plantName,
                })
            }
            this.stationCheckBoxes = {count: 0};
            this.materialToMonitor = {
                name: r.materialName,
                longCode: r.material,
            };
            this.updatingMonitor = true;
            this.configStationDlgShow = true;
        },
        removePlantFromMaterial(r, p) {
            if (this.phaseOneStationSel.includes(r)) {
                if (this.phaseOneStationSel.length === p.plants.length) {
                    confirm('warning', '删除关注油站', `您确定要删除所有的关注油站么?`,
                        () => {
                            doRequest({
                                url: '/v1/web/plant/material/monitor/config/del',
                                method: 'POST',
                                loading: true,
                                data: {
                                    material: p.material,
                                    all: true,
                                }
                            }, {
                                success: res => {
                                    message('success', '删除油站关注成功!');
                                    this.getMonitors();
                                },
                                fail: err => {
                                    message('error', '删除油站关注失败，请联系管理员或者稍后再试');
                                }
                            })
                        }
                    );
                } else {
                    this.$confirm(`您确定要删除 ${this.phaseOneStationSel.length} 个关注油站么？`, '删除关注油站', {
                        type: 'warning'
                    }).then(() => {
                        let plants = [];
                        for (let i = 0; i < this.phaseOneStationSel.length; i ++)
                            plants.push({
                                plant: this.phaseOneStationSel[i].plant
                            })
                        doRequest({
                            url: '/v1/web/plant/material/monitor/config/del',
                            method: 'POST',
                            loading: true,
                            data: {
                                material: p.material,
                                all: false,
                                plants: plants
                            }
                        }, {
                            success: res => {
                                message('success', '删除油站关注成功!');
                                this.getMonitors();
                            },
                            fail: err => {
                                message('error', '删除油站关注失败，请联系管理员或者稍后再试');
                            }
                        })
                    }).catch(() => {

                    })
                }
                
            } else {
                this.$confirm(`您确定要删除 ${r.plantName} 的关注么？`, '删除关注油站', {
                    type: 'warning'
                }).then(() => {
                    doRequest({
                        url: '/v1/web/plant/material/monitor/config/del',
                        method: 'POST',
                        loading: true,
                        data: {
                            material: p.material,
                            all: false,
                            plants: [{
                                plant: r.plant
                            }]
                        }
                    }, {
                        success: res => {
                            message('success', '删除油站关注成功!');
                            this.getMonitors();
                        },
                        fail: err => {
                            message('error', '删除油站关注失败，请联系管理员或者稍后再试');
                        }
                    })
                }).catch(() => {

                })
            }
            
        },
        clearPlantsForMaterial(r) {
            this.$confirm(`您确定要删除 ${r.materialName} 的所有关注么？`, '删除商品全部关注', {
                type: 'warning'
            }).then(() => {
                doRequest({
                    url: '/v1/web/plant/material/monitor/config/del',
                    method: 'POST',
                    loading: true,
                    data: {
                        material: r.material,
                        all: true,
                    }
                }, {
                    success: res => {
                        message('success', '删除商品关注成功');
                        this.getMonitors();
                    },
                    fail: err => {
                        message('error', '删除商品关注失败，请联系管理员或者稍后再试');
                    }
                })
            }).catch(() => {

            })
        },
        submitMonitor() {
            if (this.selectedStations.length <= 0) {
                message('error', '请至少选择一个油站');
                return;
            }

            let data = [];
            for (let i = 0; i < this.selectedStations.length; i ++) {
                data.push({
                    plant: this.selectedStations[i].orgCode,
                    status: 1
                })
            }

            doRequest({
                url: '/v1/web/plant/material/monitor/config/add',
                method: 'POST',
                data: {
                    material: this.materialToMonitor.longCode,
                    plants: data
                },
                loading: true
            }, {
                success: res => {
                    message('success', '添加商品关注成功!');
                    this.getMonitors();
                },
                fail: err => {
                    message('error', '添加商品关注失败，请联系管理员或者稍后再试')
                },
                finally: () => {
                    this.configStationDlgShow = false;
                }
            })
        },
        updateMonitor() {
            if (this.selectedStations.length <= 0) {
                message('error', '请至少选择一个油站');
                return;
            }

            let data = [];
            for (let i = 0; i < this.selectedStations.length; i ++) {
                data.push({
                    plant: this.selectedStations[i].orgCode,
                    status: 1
                })
            }

            doRequest({
                url: '/v1/web/plant/material/monitor/config/update',
                method: 'POST',
                data: {
                    material: this.materialToMonitor.longCode,
                    plants: data
                },
                loading: true
            }, {
                success: res => {
                    message('success', '更新商品关注成功');
                    this.getMonitors();
                },
                fail: err => {
                    message('error', '更新商品关注失败，请联系管理员或者稍后再试')
                },
                finally: () => {
                    this.configStationDlgShow = false;
                }
            })
        },
        getNoSalesList() {
            if (!this.dateRange || this.dateRange.length === 0) {
                message('warning', '请先选定时间范围')
                return;
            }

            doRequest({
                url: '/v1/web/plant/material/no-sales',
                method: 'POST',
                loading: true,
                data: {
                    beginDate: this.dateRange[0],
                    endDate: this.dateRange[1]
                }
            }, {
                // obj: this,
                // src: 'pmnsList',
                // dst: 'noSalesList',
                success: res => {
                    // let i = 0, j = 0, k = 0, l = 0, m = 0, v = null, s = null, x = [];
                    // for (i = 0; i < res.pmnsList.length; i ++) {
                    //     v = res.pmnsList[i];
                    //     k = 0;
                    //     m = l;
                    //     for (j = 0; j < v.plants.length; j ++) {
                    //         s = v.plants[j];
                    //         x.push({
                    //             material: v.material,
                    //             materialName: v.materialName,
                    //             orgCode: s.plant,
                    //             orgText: s.plantName,
                    //             status: s.status,
                    //             k: 0,
                    //         });
                    //         k ++;
                    //         l ++;
                    //     }
                    //     x[m].k = k;
                    //     x[m].plants = v.plants;
                    // }
                    // this.noSalesList = x;
                    let v = null, e = null, i = 0, j = 0, k = 0;
                    for (i = 0; i < res.pmnsList.length; i ++) {
                        v = res.pmnsList[i];
                        v.goodOnes = [];
                        v.showGoodOnes = false;
                        for (j = 0; j < this.monitors.length; j ++) {
                            if (v.material == this.monitors[j].material) {
                                e = this.monitors[j];
                                break;
                            }
                        }
                        for (j = 0; j < e.plants.length; j ++) {
                            for (k = 0; k < v.plants.length; k ++) {
                                if (e.plants[j].plant == v.plants[k].plant)
                                    break;
                            }
                            if (k === v.plants.length)
                                v.goodOnes.push(e.plants[j]);
                        }
                    }
                    this.noSalesList = res.pmnsList;
                }
            })
        },
        phaseOneStationSelectChange({ records  }) {
            this.phaseOneStationSel = records ;
        }
    }
}
</script>