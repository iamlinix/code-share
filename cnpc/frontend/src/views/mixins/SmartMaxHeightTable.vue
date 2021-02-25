<template>
<div style="width: 100%; height: 100%" ref="smart-table-frame">
    <el-table :max-height="tableMaxHeight" 
        :key="tableKey"
        v-loading="loading"
        v-bind="$attrs"
        v-on="$listeners"
        ref="innerTable">
        <slot></slot>
    </el-table>
</div>
</template>

<style scoped>

</style>

<script>
import { toggleKey } from '../../utils/utils';
const uuid = require('uuid');
// var elementResizeDetectorMaker = require("element-resize-detector");

export default {
    props: {
        minus: {
            type: Number,
            default: 0
        },
        loading: {
            type: Boolean,
            default: false,
        },
        showOverflowTooltip: {
            type: Boolean,
            default: true
        }
    },
    beforeUpdate() {
        this.$slots.default.forEach(v => {
            if (v.componentOptions)
                v.componentOptions.propsData.showOverflowTooltip = this.showOverflowTooltip;
        })
    },
    mounted() {
        this.tableKey = uuid.v1();
        // var erd = elementResizeDetectorMaker({
        //     strategy: "scroll"
        // });
        // erd.listenTo(this.$refs['smart-table-frame'], this.resizeTable);   
        window.addEventListener('resize', this.resizeTable);
        this.resizeTable();
    },
    data() {
        return {
            tableMaxHeight: 500,
            tableKey: 'table-key',

            lastWidth: 0,
            lastHeight: 0,
            counter: 0
        }
    },
    methods: {
        resizeTable(ev) {
            this.tableMaxHeight = window.innerHeight - this.minus;
            toggleKey(this, 'tableKey');
            
            // if (ev.clientHeight === this.lastHeight && ev.clientWidth === this.lastWidth) {
            //     if (this.counter === this.maxPeriod) {
            //         this.tableMaxHeight = window.innerHeight - this.minus;
            //         toggleKey(this, 'tableKey');
            //         this.counter = 0;
            //     } else {
            //         this.counter += 1;
            //         setTimeout(this.resizeTable, this.interval, ev);
            //     }
            // } else {
            //     setTimeout(this.resizeTable, this.interval, ev);
            // }
        },
        refresh() {
            toggleKey(this, 'tableKey');
        },
        doLayout() {
            this.$refs.innerTable.doLayout();
        },
        toggleRowSelection(row, b) {
            this.$refs.innerTable.toggleRowSelection(row, b);
        }
    }
}
</script>