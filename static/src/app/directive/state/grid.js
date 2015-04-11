define([], function(){

    function directive($http) {
        return {
            templateUrl: '/static/src/app/directive/state/grid.html',
            restrict: 'E',
            scope: {
                gridRows: '=rows'
            },
            link: function postLink(scope, element, attrs) {
            }
        }
    }
    directive.$inject=[];

    return directive;
});
