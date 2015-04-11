define([
    './controller/dashboard',
    './directive/state/grid'
],
function (dashboardController, statesGrid) {
    var app = angular.module('monitor', ['ngRoute', 'ngWebsocket', 'angularGrid']);

    //module config
    app.config(['$routeProvider', function($routeProvider){
        $routeProvider
            .when('/dashboard', {
                templateUrl: '/static/src/app/view/dashboard.html',
                controller: 'dashboardController'
            })
            .otherwise({redirectTo: '/dashboard'});
    }]);

    //register controllers
    app.controller('dashboardController', dashboardController);

    //register directives
    app.directive('stateGrid', statesGrid);

});
