define([], function () {

    function controller($scope, $websocket) {

        $scope.stateData = [];
        $scope.numOnline = 0;
        $scope.numOffline = 0;
        $scope.connected = false;
        $scope.filterCategory = 'output';
        $scope.filterKeyword = '';

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

                //update online/offline
                angular.forEach(data, function(row, key) {

                    if ($scope.filterCategory !== row.PluginCategory) {
                        return;
                    }

                    $scope.stateData.push(row);

                    if (row.HostUp) {
                        $scope.numOnline++;
                    } else {
                        $scope.numOffline++;
                    }
                });
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
