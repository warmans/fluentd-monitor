define([
    './controller/dashboard',
    './directive/state/grid',
    './directive/sparkline'
],
function (dashboardController, statesGridDirective, sparklineDirective) {
    var app = angular.module('monitor', ['ngRoute', 'ngWebsocket']);

    //module config
    app.config(['$routeProvider', function($routeProvider){
        $routeProvider
            .when('/dashboard', {
                templateUrl: '/ui/src/app/view/dashboard.html',
                controller: 'dashboardController'
            })
            .otherwise({redirectTo: '/dashboard'});
    }]);

    //register controllers
    app.controller('dashboardController', dashboardController);

    //register directives
    app.directive('stateGrid', statesGridDirective);
    app.directive('sparkline', sparklineDirective);

});
