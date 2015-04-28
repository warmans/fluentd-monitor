define([
    './controller/dashboard',
    './controller/head',

    './directive/state/grid',
    './directive/sparkline',

    './factory/head',
    './factory/bytes',

    './filter/bytes',
],
function (dashboardController, headController, statesGridDirective, sparklineDirective, headFactory, bytesFormatterFactory, bytesFilter) {

    var app = angular.module('monitor', ['ngRoute', 'ngWebsocket']);

    //module config
    app.config(['$routeProvider', function($routeProvider){
        $routeProvider
            .when('/dashboard', {
                templateUrl: '/ui/src/app/view/dashboard.html',
                controller: 'dashboardController',
                reloadOnSearch: false
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
    app.factory('bytesFormatter', bytesFormatterFactory);

    //register filters
    app.filter('bytes', bytesFilter);

});
