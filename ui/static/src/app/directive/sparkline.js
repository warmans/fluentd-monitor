define([], function(){

    function directive() {
        return {
            restrict: 'A',
            require: 'ngModel',
            link: function (scope, elem, attrs, ngModel) {
                //defaults
                var opts = {
                    type: 'line',
                    lineColor: 'rgba(255, 255, 255, 0.3)',
                    fillColor: 'rgba(255, 255, 255, 0.2)',
                    normalRangeColor: 'transparent',
                    normalRangeMin: 0,
                    normalRangeMax: 0,
                    spotColor: false,
                    minSpotColor: false,
                    maxSpotColor: false,
                    drawNormalOnTop: false,
                    height: 30,
                };

                opts.type = attrs.type || 'line';

                scope.$watch(attrs.ngModel, function () {
                    var model = ngModel.$viewValue;
                    if (!model || !angular.isArray(model)) { return; }

                    // Make sure we have an array of numbers
                    $(elem).sparkline(model, opts);
                });
            }
        }
    }
    directive.$inject=[];

    return directive;
});
