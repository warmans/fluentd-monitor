define([], function () {

    function controller($scope, $websocket) {

        $scope.stateData = [];
        $scope.numOnline = 0;
        $scope.numOffline = 0;
        $scope.connected = false;
        $scope.filterKeyword = '';

        $scope.chartData = [];
        $scope.chartOptions = {
            grid: { hoverable: true, borderWidth: 0},
            series: { shadowSize: 0, stack: true },
            lines : { lineWidth : 1, fill: true },
            legend: { show: false },
            tooltip: true,
            tooltipOpts: {},
            xaxis: { show: false },
            yaxis: {
                min: 0,
                tickFormatter: function(val, axis) {
                switch(true) {
                    case val > 1000000000:
                        return (val/1000000000).toFixed(2) + "B";
                    case val > 1000000:
                        return (val/1000000).toFixed(2) + "M";
                    case val > 1000:
                        return (val/1000).toFixed(2) + "K";
                    default:
                        return val;
                }
                  return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
                }
            },
            colors: ['#0E86CC']
        };

        var ws = $websocket.$new({
            url: "ws://"+location.host+"/ws",
            protocols: []
        });

        ws.$on('$open', function () {
            console.log('connected');
            $scope.$apply(function() {
                $scope.connected = true;
            });
        });

        ws.$on('$close', function () {
            console.log('disconnected');
            $scope.$apply(function() {
                $scope.connected = false;
            });
        });

        ws.$on('$message', function (data) {
            $scope.$apply(function() {

                //reset
                $scope.numOnline = 0;
                $scope.numOffline = 0;
                $scope.stateData = [];
                $scope.chartData = [];

                //update online/offline
                angular.forEach(data, function(row, key) {

                    //update table
                    $scope.stateData.push(row);

                    //update chart
                    $scope.chartData.push({label: row.ID, data: row.BufferTotalQueuedSize.map(function(val, key) { return [key, val]; })});

                    if (row.HostUp) {
                        $scope.numOnline++;
                    } else {
                        $scope.numOffline++;
                    }
                });

                console.log($scope.chartData);
            });
        });

        $scope.selectedRow = null;
        $scope.onRowSelect = function(row) {
            $scope.selectedRow = row;
        };
    }

    controller.$inject=['$scope', '$websocket'];

    return controller;
});
