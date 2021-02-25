<template>
    <div style="width: 100%">
        <vxe-grid :max-height="tableMaxHeight" 
            :key="tableKey"
            v-bind="$attrs"
            v-on="$listeners">
            <slot></slot>
        </vxe-grid>
    </div>
</template>

<style scoped>

</style>

<script>
import { toggleKey } from '../../utils/utils';
const uuid = require('uuid');

export default {
    props: {
        minus: {
            type: Number,
            default: 0
        }
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
    }
}
</script>