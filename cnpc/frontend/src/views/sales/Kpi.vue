<template>
    <div class="cnpc-sales-kpi">
        <el-card>
            <el-form inline label-width="70px" size="mini" label-position="right">
                <el-form-item label="起始日期">
                    <el-date-picker type="daterange" v-model="dateRange" range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd" />
                </el-form-item>
                <el-form-item label="层级">
                    <el-select v-model="orgLevel" placeholder="请选择统计层级" clearable>
                        <el-option value="0" label="分公司" />
                        <el-option value="1" label="加油站" />
                    </el-select>
                </el-form-item>
                <el-form-item label="单位">
                    <el-select v-model="orgCode" placeholder="请选择要统计的单位" clearable>
                        <el-option value="0" label="分公司" />
                        <el-option value="1" label="加油站" />
                    </el-select>
                </el-form-item>
            </el-form>
        </el-card>
    </div>
</template>

<style lang="scss">
.cnpc-sales-kpi {
    width: 100%;
    max-height: calc(100vh - 12px);
    background-color: #FFF;
    padding: 8px 8px 8px 8px;
    overflow: auto;

    .el-form-item {
        margin-bottom: 0;
    }
}
</style>

<script>
import { doRequest, allRequests, message, confirm } from '../../utils/utils'

export default {
    data() {
        return {
            dateRange: [],
            orgLevel: '',
            orgCode: '',

            companies: [],
            stations: []
        }
    },
    mounted() {
        allRequests((x) => {
            console.log(x)
        }, [{
            url:'/v1/web/org/branch'
        }, {
            'url': '/v1/web/org/plant'
        }])
    },
    methods: {

    }
}
</script>