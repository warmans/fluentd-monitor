define([
    './controller/dashboard',
    './controller/head',

    './directive/state/grid',
    './directive/sparkline',

    './factory/head',
],
function (dashboardController, headController, statesGridDirective, sparklineDirective, headFactory) {
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
    app.controller('headController', headController);

    //register directives
    app.directive('stateGrid', statesGridDirective);
    app.directive('sparkline', sparklineDirective);

    //register factories
    app.factory('pageHead', headFactory);

});
