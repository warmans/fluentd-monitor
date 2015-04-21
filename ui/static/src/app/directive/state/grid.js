define([], function(){

    function directive() {
        return {
            templateUrl: '/ui/src/app/directive/state/grid.html',
            restrict: 'E',
            scope: {
                gridRows: '=rows',
                onRowSelect: '=onRowSelect'
            },
            link: function postLink(scope, element, attrs) {
                scope.selectedRow = null;
                scope.selectRow = function(row) {
                    if (scope.selectedRow == row.ID) {
                        scope.selectedRow = null;
                        scope.onRowSelect(null);
                    } else {
                        scope.selectedRow = row.ID;
                        scope.onRowSelect(row);
                    }
                };
            }
        }
    }
    directive.$inject=[];

    return directive;
});
