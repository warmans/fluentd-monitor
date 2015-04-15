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

                var unWatch = scope.$watch(attrs.ngModel, function () {
                    var model = ngModel.$viewValue;
                    if (!model || !angular.isArray(model)) { return; }

                    // Make sure we have an array of numbers
                    $(elem).sparkline(model, opts);

                    //seems there is a bug where failing to call this might cause a memory leak
                    $.sparkline_display_visible();

                    unWatch();
                });

                var unOn = scope.$on('$destroy', function() {
                    //prevent leaks from events hanging around
                    $(elem).off();
                    $(elem).remove();
                    unOn();
                });
            }
        }
    }
    directive.$inject=[];

    return directive;
});
