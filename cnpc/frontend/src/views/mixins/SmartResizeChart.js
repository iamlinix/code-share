/*
*/
export default {
    data: function() {
        return {
            smrschtIdToMonitor: {},
            smrschtLastWidth: -1,
            smrschtLastHeight: -1,
            smrschtCharts: [],
            smrschtDelay: 50,
            smrschtResizeCounter: 0,
            smrschtResizeMax: 3
        };
    },
    mounted: function() {
        window.addEventListener('resize', this.smrschtDoSmartDelayedResize);
    },
    methods: {
        smrschtMonitorResizeById: function(idToMonitor) {
            this.smrschtIdToMonitor[idToMonitor] = {
                lastWidth: -1,
                lastHeight: -1
            };
        },
        smrschtAddChart: function(chart) {
            this.smrschtCharts.push(chart);
        },
        smrschtDoSmartDelayedResize: function () {
            let ele = null, v = null, allIsGood = true;
            for (let k in this.smrschtIdToMonitor) {
                ele = document.getElementById(k);
                v = this.smrschtIdToMonitor[k];
                if (ele) {
                    if (ele.clientWidth !== v.lastWidth || ele.clientHeight !== v.lastHeight) {
                        v.lastWidth = ele.clientWidth;
                        v.lastHeight = ele.clientHeight;
                        allIsGood = false;
                        break;
                    }
                } else {
                    console.log(`cannot find element with id: ${k}`);
                }
            }

            if (allIsGood) {
                this.smrschtResizeCounter += 1;
                if (this.smrschtResizeCounter >= this.smrschtResizeMax) {
                    this.smrschtResizeCounter = 0;
                    if (this.smrschtCharts && this.smrschtCharts.length > 0) {
                        this.smrschtCharts.forEach(e => {
                            e.resize();
                        });
                    } else {
                        console.log("no charts to resize...");
                    }
                } else {
                    setTimeout(() => {
                        this.smrschtDoSmartDelayedResize();
                    }, this.smrschtDelay);
                }
            } else {
                setTimeout(() => {
                    this.smrschtDoSmartDelayedResize();
                }, this.smrschtDelay);
            }
            /*
            if (this.smrschtIdToMonitor && this.smrschtIdToMonitor.length > 0) {
                let ele = document.getElementById(this.smrschtIdToMonitor);
                if (ele) {
                    if (ele.clientWidth !== this.smrschtLastWidth || ele.clientHeight !== this.smrschtLastHeight) {
                        this.smrschtLastWidth = ele.clientWidth;
                        this.smrschtLastHeight = ele.clientHeight;
                        setTimeout(() => {
                            this.smrschtDoSmartDelayedResize();
                        }, this.smrschtDelay);
                    } else {
                        this.smrschtResizeCounter += 1;
                        if (this.smrschtResizeCounter >= this.smrschtResizeMax) {
                            this.smartResizeWidth = 0;
                            this.smartResizeHeight = 0;
                            this.smrschtResizeCounter = 0;
                            if (this.smrschtCharts && this.smrschtCharts.length > 0) {
                                this.smrschtCharts.forEach(element => {
                                    element.resize();
                                });
                            } else {
                                console.log("no charts to resize...");
                            }
                        } else {
                            setTimeout(() => {
                                this.smrschtDoSmartDelayedResize();
                            }, this.smrschtDelay);
                        }
                    }
                } else {
                    console.log(`cannot find element with id: ${this.smrschtIdToMonitor}`);
                }
            } else {
                console.log("no element has been assigned to monitor...");
            }
            */
        }
    },
};