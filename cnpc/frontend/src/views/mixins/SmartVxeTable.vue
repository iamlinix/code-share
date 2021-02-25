<template>
    <div style="width: 100%" v-resize:debounce.50="resizeTable" :id="divId">
        <vxe-table :max-height="tableMaxHeight" 
            :key="tableKey"
            v-bind="$attrs"
            v-on="$listeners"
            ref="xTableInstance">
            <slot></slot>
        </vxe-table>
    </div>
</template>

<style scoped>

</style>

<script>
import { toggleKey } from '../../utils/utils';
import resize from 'vue-resize-directive'
const uuid = require('uuid');

export default {
    directives: {
        resize
    },
    props: {
        minus: {
            type: Number,
            default: 0
        }
    },
    beforeMount() {
        this.divId = uuid.v1();
        this.tableKey = uuid.v1();
    },
    mounted() {
        // window.addEventListener('resize', this.resizeTable);
        this.resizeTable({ width: 0, height: 0});
        let ev = document.getElementById('mainContainer');
        if (ev) {
            this.lastMainContainerWidth = ev.clientWidth;
            this.lastMainContainerHeight = ev.clientHeight;
        }
    },
    data() {
        return {
            divId: '',
            tableMaxHeight: 500,
            tableKey: 'table-key',
            lastWidth: -1,
            lastHeight: -1,
            lastMainContainerWidth: -1,
            lastMainContainerHeight: -1,
        }
    },
    methods: {
        resizeTable(ele) {
            // FIXME
            if (ele.clientWidth != this.lastWidth) {
                this.tableMaxHeight = window.innerHeight - this.minus;
                toggleKey(this, 'tableKey');
                this.lastWidth = ele.clientWidth;
                this.lastHeight = ele.clientHeight;
            }
            // let ev = document.getElementById('mainContainer');
            // if (ev) {
            //     if (this.lastMainContainerWidth != ev.clientWidth || 
            //         this.lastMainContainerHeight != ev.clientHeight) {
            //         this.lastMainContainerWidth = ev.clientWidth;
            //         this.lastMainContainerHeight = ev.clientHeight;
            //         if (ele.clientWidth != this.lastWidth || ele.clientHeight != this.lastHeight) {
            //             this.tableMaxHeight = window.innerHeight - this.minus;
            //             toggleKey(this, 'tableKey');
            //             this.lastWidth = ele.clientWidth;
            //             this.lastHeight = ele.clientHeight;
            //         }
            //     }
            // }
        },
        refresh(force = false) {
            let ele = document.getElementById(this.divId)
            if (force || ele.clientWidth != this.lastWidth || ele.clientHeight != this.lastHeight) {
                this.tableMaxHeight = window.innerHeight - this.minus;
                toggleKey(this, 'tableKey');
                this.lastWidth = ele.clientWidth;
                this.lastHeight = ele.clientHeight;
            }
        },
        setCheckboxRow(a, b) {
            this.$refs['xTableInstance'].setCheckboxRow(a, b);
        }
    }
}
</script>