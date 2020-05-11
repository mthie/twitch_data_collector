var twitchCollector = {
    updateHooks: [],
    data: null,


    update: function(data) {
        this.data = data;
        var that = this;
        this.updateHooks.forEach(function(v) {
            v(that.data);
        });
    },

    registerHook: function(kind, callback) {
        switch (kind) {
            case "update":
                this.updateHooks.push(callback);
                break;
        }
    },

    start: function() {
        this.getData();
        var that = this;
        window.setInterval(function() {
            that.getData()
        }, 20000);
    },

    getData: function() {
        let xmlhttp = new XMLHttpRequest();
        let url = "/data";
        var that = this;

        xmlhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                let data = JSON.parse(this.responseText);
                if (data !== null) {
                    that.update(data);
                }
            }
        };
        xmlhttp.open("GET", url, true);
        xmlhttp.send();
    }
}