define([], function () {

    function controller($scope, $location, $websocket, filterFilter, pageHead, bytesFormatter) {

        $scope.rawStateData = [];
        $scope.filteredStateData = [];
        $scope.numOnline = 0;
        $scope.numOffline = 0;
        $scope.connected = false;
        $scope.filterKeyword = $location.search().filter || '';

        $scope.selectedRow = null;
        $scope.onRowSelect = function(row) {
            $scope.selectedRow = row;
        };

        $scope.chartData = [];
        $scope.chartOptions = {
            grid: { hoverable: true, borderWidth: 0},
            series: { shadowSize: 0, stack: true },
            lines : { lineWidth : 1, fill: true },
            legend: { show: false },
            tooltip: true,
            tooltipOpts: {
                content: "%s: %y",
                defaultTheme: false
            },
            xaxis: { show: false },
            yaxis: {
                min: 0,
                tickFormatter: bytesFormatter.format
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
                $scope.numWarnings = 0;
                $scope.stateData = [];
                $scope.chartData = [];

                $scope.rawStateData = data;

                //update online/offline
                angular.forEach(data, function(row, key) {
                    $scope.numWarnings += row.Warnings.length;

                    if (row.HostUp) {
                        $scope.numOnline++;
                    } else {
                        $scope.numOffline++;
                    }
                });

                //notify via tab title
                pageHead.setTitlePrefix($scope.numOffline == 0 ? '' : "("+$scope.numOffline+")");
                if ($scope.numOffline == 0) {
                    if ($scope.numWarnings == 0) {
                        pageHead.setFavicon('fluentd-logo-green-16.png');
                    } else {
                        pageHead.setFavicon('fluentd-logo-orange-16.png');
                    }
                } else {
                    pageHead.setFavicon('fluentd-logo-red-16.png');
                }
            });
        });

        $scope.$watch('rawStateData', function(newVal, oldVal) {

            //update URI with filter
            $location.search('filter', $scope.filterKeyword);

            //filter results
            $scope.filteredStateData = filterFilter(newVal, $scope.filterKeyword);

            //update overview graph
            angular.forEach($scope.filteredStateData, function(row, key) {
                //update chart with rows that have non-zero TotalQueued total
                if (row.BufferTotalQueuedSize.reduce(function(a, b) { return a + b; }) > 0) {
                    $scope.chartData.push({
                        label: row.ID,
                        data: row.BufferTotalQueuedSize.map(function(val, key) { return [key, val]; })
                    });
                }
            });

        });
    }

    controller.$inject=['$scope', '$location', '$websocket', 'filterFilter', 'pageHead', 'bytesFormatter'];

    return controller;
});
