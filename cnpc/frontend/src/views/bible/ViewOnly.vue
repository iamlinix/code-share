<template>
    <div class="cnpc-bible">
        <el-dialog :visible.sync="bibleDialogVisible" @closed="selectedBible = {}"
            :close-on-click-modal="false">
            <p slot="title" style="text-align: center">
                {{ isEditing ? selectedBible.materialTxt : '添加天书条目' }}
            </p>
            <el-form label-position="right" label-width="100px" size="mini">
                <el-form-item label="商品编码">
                    <el-input v-model="selectedBible.material" size="mini"/>
                </el-form-item>
                <el-form-item label="商品名称">
                    <el-input v-model="selectedBible.materialTxt" size="mini"/>
                </el-form-item>
                <el-form-item label="供应商编码">
                    <el-input v-model="selectedBible.vendorCode" size="mini"/>
                </el-form-item>
                <el-form-item label="供应商名称">
                    <el-input v-model="selectedBible.vendorName" size="mini"/>
                </el-form-item>
                <el-row :gutter="6">
                    <el-col :span="8">
                        <el-form-item label="商品大类">
                            <el-input v-model="selectedBible.mainClass" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="16">
                        <el-form-item label="大类描述">
                            <el-input v-model="selectedBible.mainClassTxt" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="8">
                        <el-form-item label="商品中类">
                            <el-input v-model="selectedBible.subClass" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="16">
                        <el-form-item label="中类描述">
                            <el-input v-model="selectedBible.subClassTxt" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="8">
                        <el-form-item label="量纲">
                            <el-input v-model="selectedBible.groes" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="8">
                        <el-form-item label="计量单位">
                            <el-input v-model="selectedBible.baseUom" size="mini"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="8">
                        <el-form-item label="起始日期">
                            <el-date-picker v-model="selectedBible.validSDate" size="mini" type="date"
                                value-format="yyyy-MM-dd" style="width: 100%"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="采购价">
                            <el-input-number style="width: 100%" v-model="selectedBible.purchasePrice" 
                                size="mini" :min="0.0" :step="0.01"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="零售价">
                            <el-input-number style="width: 100%" v-model="selectedBible.salePrice" 
                                size="mini" :min="0.0" :step="0.01"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="状态">
                            <el-input-number style="width: 100%" v-model="selectedBible.status" 
                                size="mini" :min="0" :max="1"/>
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="库存">
                            <el-input-number style="width: 100%" v-model="selectedBible.stock" 
                                size="mini" :min="0"/>
                        </el-form-item>
                    </el-col>
                </el-row>
            </el-form>
            <el-button style="width: 100%" :type="selectedBible.materialTxt ? 'success' : 'primary'"
                @click="createBible">
                {{ selectedBible.materialTxt ? '更新条目' : '创建条目' }}
            </el-button>
        </el-dialog>
        <p class="title">天书</p>
        <!--el-row :gutter="6" style="margin-bottom: 8px">
            <el-col :span="3">
                <el-button size="mini" style="width: 100%" 
                    icon="el-icon-download" type="success" 
                    @click="downloadBibleCSV">导出天书</el-button>
            </el-col>
            <el-col :span="3">
                <el-upload accept="*" :headers="{'token': getTokenBible()}" 
                    :action="uploadUrl" :show-file-list="false" :multiple="false"
                    :before-upload="bibleUploadBefore"
                    :on-success="bibleUploadSuccess"
                    :on-error="bibleUploadError"
                    name="upload">
                    <el-button style="width: 100%" size="mini" 
                        icon="el-icon-upload2" type="primary">导入天书</el-button>
                </el-upload>
            </el-col>
            <el-col :span="3">
                <el-button size="mini" style="width: 100%" 
                    icon="el-icon-plus" type="warning" 
                    @click="showBibleDlg">添加天书内容</el-button>
            </el-col>
            <el-col :span="3">
                <el-button size="mini" style="width: 100%" 
                    icon="el-icon-delete" type="danger" 
                    @click="clearBible">清空天书</el-button>
            </el-col>
        </el-row-->
        <!--s-table :minus="125" v-loading="tableLoading" :data="bibles" :stripe="true" 
            size="medium" :border="true">
            <el-table-column type="index" width="70" align="center" fixed />
            <el-table-column width="160" label="商品编码" prop="material" sortable/> 
            <el-table-column width="250" label="商品名称" prop="materialTxt" sortable/>   
            <el-table-column width="120" label="供应商编码" prop="vendorCode" sortable/> 
            <el-table-column width="250" label="供应商名称" prop="vendorName" sortable/>  
            <el-table-column width="110" label="商品大类" prop="mainClass" sortable/> 
            <el-table-column width="120" label="大类名称" prop="mainClassTxt" sortable/> 
            <el-table-column width="110" label="商品中类" prop="subClass" sortable/> 
            <el-table-column width="120" label="中类名称" prop="subClassTxt" sortable/> 
            <el-table-column width="110" label="起始日期" prop="validSDate" sortable/> 
            <el-table-column width="100" label="量纲" prop="groes" sortable/> 
            <el-table-column width="100" label="单位" prop="basUom" sortable/> 
            <el-table-column width="100" label="采购价" prop="purchasePrice" sortable/> 
            <el-table-column width="100" label="零售价" prop="salePrice" sortable/> 
            <el-table-column width="110" label="商品状态" prop="status" sortable/> 
            <el-table-column width="100" label="库存" prop="stock" sortable/> 
        </s-table-->

        <s-table  showOverflow="tooltip" :loading="tableLoading"
            :fixed-header="true" :data="filteredBibles" :columns="tableColumns"
            :minus="120" :border="true" :stripe="true" size="small"
            :highlight-current-row="true" :highlight-hover-row="true">
        </s-table>
    </div>
</template>

<style lang="scss">
.cnpc-bible {
    background-color: #FFF;
    height: calc(100vh - 10px);
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

    table.vgt-table {
        font-size: 14px;
    }

    .el-input-group__append {
        padding: 0 10px;
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
import SmartTable from '../mixins/SmartMaxHeightVxeTable'
import VXETable from 'vxe-table';

export default {
    components: {
        // eslint-disable-next-line vue/no-unused-components
        's-table': SmartTable
    },
    beforeMount() {
        let self = this;
        VXETable.renderer.add('btnGroup', {
            renderDefault (h, renderOpts, params) {
                let { row, column } = params
                let { events } = renderOpts
                return [
                    <el-button size="mini" circle type="primary" icon="el-icon-edit" 
                        onClick={() => events.editBible(row)}/>,
                    <el-button size="mini" circle type="warning" icon="el-icon-close" 
                        onClick={() => events.deleteBible(row)}/>
                ]
          }
        });
        VXETable.renderer.add('searchHeader', {
            renderHeader (h, renderOpts, params) {
                let { columnIndex } = params
                let { events, model } = renderOpts
                
                return [
                    <el-input size="mini" placeholder="输入编码或名称检索" v-model={ self[model] }
                        style="width: 90%" clearable onChange={ events.searchBible }>
                        <el-button slot="append" icon="el-icon-search" size="mini"
                            onClick={ events.searchBible }></el-button>
                    </el-input>
                ]
                
          }
        });
    },
    mounted() {
        this.getBibles();
    },
    data() {
        return {
            materialSearch: '',
            vendorSearch: '',
            filteredBibles: [],
            bibles: [],
            tableLoading: false,
            isEditing: false,
            selectedBible: {},
            bibleDialogVisible: false,
            tableColumns: [{
                type: 'seq',
                width: 70,
                fixed: 'left',
                align: 'center'
            }, {
                field: 'material',
                title: '商品编码',
                width: 160,
                fixed: 'left',
                sortable: true
            }, {
                field: 'materialTxt',
                title: '商品名称',
                width: 260,
                sortable: true,
                editRender: {
                    name: 'searchHeader',
                    events: {
                        searchBible: this.doSearch
                    },
                    model: 'materialSearch'
                }
            }, {
                field: 'vendorCode',
                title: '供应商编码',
                width: 120,
                sortable: true
            }, {
                field: 'vendorName',
                title: '供应商名称',
                width: 260,
                sortable: true,
                editRender: {
                    name: 'searchHeader',
                    events: {
                        searchBible: this.doSearch
                    },
                    model: 'vendorSearch'
                }
            }, {
                field: 'mainClass',
                title: '商品大类',
                width: 100,
                sortable: true
            }, {
                field: 'mainClassTxt',
                title: '大类描述',
                width: 120,
                sortable: true
            }, {
                field: 'subClass',
                title: '商品中类',
                width: 100,
                sortable: true
            }, {
                field: 'subClassTxt',
                title: '中类描述',
                width: 180,
                sortable: true
            }, {
                field: 'groes',
                title: '量纲',
                width: 100,
                sortable: true
            }, {
                field: 'baseUom',
                title: '计量单位',
                width: 100,
                sortable: true
            }, {
                field: 'purchasePrice',
                title: '采购价格',
                width: 100,
                sortable: true
            }, {
                field: 'salePrice',
                title: '零售价格',
                width: 100,
                sortable: true
            }, {
                field: 'status',
                title: '商品状态',
                width: 100,
                sortable: true
            }, {
                field: 'stock',
                title: '商品库存',
                width: 100,
                sortable: true
            }, {
                field: 'validSDate',
                title: '起始日期',
                width: 100,
                sortable: true
            }, {
                field: 'createdTime',
                title: '修改日期',
                width: 100,
                sortable: true
            }]
        }
    },
    methods: {
        getBibles() {
            this.tableLoading = true;
            doRequest({
                url: '/v1/web/erp/bible/list?all=1',
                method: 'GET'
            }, {
                obj: this,
                src: 'biList',
                dst: 'bibles',
                success: res => {
                    this.filteredBibles = deepCopy(this.bibles);
                },
                fail: err => {
                    message('error', '加载天书失败，请刷新页面再试');
                },
                finally: () => {
                    this.tableLoading = false;
                }
            })
        },
        downloadBibleCSV() {
            doRequestv2({
                url: '/v1/web/erp/bible/export',
                method: 'GET',
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
                    message('error', '下载天书失败')
                    console.log(err);
                }
            })
        },
        bibleUploadBefore(file) {
            if (!this.localLoading) {
                this.localLoading = Loading.service({
                    text: `${file.name} 上传中...`
                })
            }
        },
        bibleUploadSuccess(file) {
            if (this.localLoading) {
                this.localLoading.close();
                this.localLoading = null;
            }

            message('success', `${file.name} 上传成功`);
            this.getBibles();
        },
        bibleUploadError(err, file, fileList) {
            if (this.localLoading) {
                this.localLoading.close();
                this.localLoading = null;
            }

            message("error", `${file.name} 上传失败，请稍后再试`);
        },
        getTokenBible() {
            return getToken();
        },
        clearBible() {
            let self = this;
            confirm('warning', '清空天书', '您确定要清空所有的天书内容么?', () => {
                doRequest({
                    url: '/v1/web/erp/bible/clear',
                    method: 'GET',
                    loading: true
                }, {
                    success: res => {
                        message('success', '天书清空成功');
                        this.bibles = [];
                        this.filteredBibles = [];
                    },
                    fail: err => {
                        message('error', '天书清空失败');
                        console.log(err);
                    }
                })
            })
        },
        showBibleDlg() {
            this.isEditing = false;
            this.bibleDialogVisible = true;
        },
        deleteBible(row) {
            let self = this;
            confirm('warning', '删除天书', `确定删除 ${row.material}(${row.materialTxt}) 的天书条目么?`, () => {
                doRequest({
                    url: `/v1/web/erp/bible/del/${row.material}`,
                    method: 'GET',
                    loading: true
                }, {
                    success: res => {
                        message('success', '删除天书条目成功');
                        self.getBibles();
                    },
                    fail: err => {
                        console.log(err);
                        message('error', '删除天书条目失败')
                    }
                })
            })
        },
        editBible(row) {
            this.isEditing = true;
            this.selectedBible = deepCopy(row);
            this.bibleDialogVisible = true;
        },
        createBible() {
            doRequest({
                url: `/v1/web/erp/bible/${this.isEditing ? 'update' : 'add'}`,
                method: 'POST',
                loading: true,
                data: this.selectedBible
            }, {
                success: res => {
                    message('success', this.isEditing ? '天书更新成功' : '天书创建成功');
                    this.bibleDialogVisible = false;
                    this.getBibles();
                },
                fail: err => {
                    console.log(err);
                    message('error', this.isEditing ? '天书更新失败' : '天书创建失败')
                }
            });
        },
        doSearch() {
            this.filteredBibles = this.bibles.filter(e => {
                let acc = 0;
                if (this.materialSearch.length > 0)
                    acc += (e.material.toLowerCase().includes(this.materialSearch.toLowerCase()) || 
                        e.materialTxt.toLowerCase().includes(this.materialSearch.toLowerCase())) ? 1: 0
                else
                    acc += 1;

                if (this.vendorSearch.length > 0)
                    acc += (e.vendorCode.toLowerCase().includes(this.vendorSearch.toLowerCase()) || 
                        e.vendorName.toLowerCase().includes(this.vendorSearch.toLowerCase())) ? 1: 0
                else
                    acc += 1;

                return acc === 2;

            })
        }
    },
    computed: {
        uploadUrl: function() {
            return Config.apiUrl + '/v1/web/erp/bible/import';
        }
    },
}
</script>