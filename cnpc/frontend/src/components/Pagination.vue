<template>
    <div class="pagination">
        <el-pagination
                @current-change="handleCurrentChange"
                :current-page="page"
                :page-size="pageSize"
                :pager-count="5"
                layout="total, prev, pager, next"
                :total="total"
                :hide-on-single-page="true"
                :background="true">
        </el-pagination>
    </div>
</template>

<script>
    //如果需要保存分页信息和搜索数据可以在组件销毁之前用页面关键key保存在本地
    //也可以使用 keep-alive保存页面
    import G from 'lodash/get'
    import { doRequest } from "@/utils/utils";

    export default {
        name: "Pagination",
        props: {
            responseObj: {
                type: Object,
                default: undefined,
            },
            params: {
                type: Object,
                default: undefined,
            },
            pageSize: {
                type: Number,
                default: 10,
            },
            hub: {
                type: Object,
                default: undefined,
            },
            refreshCommand: {
                type: String,
                default: undefined,
            },
            resetRefreshCommand: {
                type: String,
                default: undefined,
            },
            field: {
                type: String,
                default: undefined,
            }
        },
        data: () => {
            return {
                total: 0,
                page: 1,
            };
        },
        methods: {
            handleCurrentChange(val) {
                this.getPageData(val - 1);
                //console.log(`当前页: ${val}`);
            },
            getPageData(p, force) {
                doRequest({...this.params, params: {page: p}}, {
                    success: res => {
                        this.total = G(res, 'total', 0);
                        if (this.total > 0 && this.field) {
                            const arr = res[this.field];
                            if ((!arr || arr.length === 0) && p > 1) {
                                this.getPageData(p - 1);
                                return;
                            }
                        }

                        if (this.responseObj) {
                            if (this.responseObj.after) {
                                this.responseObj.after(res);
                            }

                            if ((this.page !== (p + 1) || force) && this.responseObj.pageChange) {
                                this.responseObj.pageChange(p + 1, force);
                            }
                        }

                        this.page = p + 1;
                    },
                    fail: err => {
                        console.log('pagination error', err);
                    }
                });
                /*
                this.requestFunc({...this.params, pageSize: this.pageSize, page: p}).then(r => {

                    this.total = G(r, 'total', 0);
                    this.page = p;
                    this.$emit('returnData', r);//this.$emit('returnData', G(r,'list',[]));
                }).catch(_ => {
                })
                */
            },
            refresh() {
                this.getPageData(this.page - 1, true);
            },
            resetAndRefresh() {
                this.page = 1;
                this.refresh();
            }
        },
        mounted: function () {
            this.getPageData(0);
            if (this.hub)
                if (this.refreshCommand) {
                    this.hub.$on(this.refreshCommand, this.refresh);
                }
            if (this.resetRefreshCommand) {
                this.hub.$on(this.resetRefreshCommand, this.resetAndRefresh);
            }

            if (this.responseObj) {
                if (this.responseObj.pageChange) {
                    setTimeout(_ => {
                        this.responseObj.pageChange('x')
                    }, 100);
                }
            }
        },
    }
</script>

<style scoped>
    .pagination {
        padding-top: 8px;
        padding-bottom: 8px;
        width: 50%;
        margin: 0 auto;
        text-align: center;
    }
</style>