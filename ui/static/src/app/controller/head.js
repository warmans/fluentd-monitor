define([], function () {

    function controller($scope, pageHead) {
        $scope.head = pageHead;
    }

    controller.$inject=['$scope', 'pageHead'];

    return controller;
});
