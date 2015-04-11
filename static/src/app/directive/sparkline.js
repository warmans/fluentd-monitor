define([], function(){

    function directive() {
        return {
            restrict: 'A',
            require: 'ngModel',
            link: function (scope, elem, attrs, ngModel) {
                //defaults
                var opts = {
                    type: 'line',
                    lineColor: '#ccc',
                    fillColor: '#ccc',
                    normalRangeColor: 'transparent',
                    normalRangeMin: 0,
                    normalRangeMax: 0,
                    spotColor: false,
                    minSpotColor: false,
                    maxSpotColor: false,
                    drawNormalOnTop: false
                };

                //TODO: Use $eval to get the object
                opts.type = attrs.type || 'line';

                scope.$watch(attrs.ngModel, function () {
                    render();
                });

                scope.$watch(attrs.opts, function(){
                    render();
                });

                var render = function () {
                    var model;
                    if(attrs.opts) angular.extend(opts, angular.fromJson(attrs.opts));
                    // Trim trailing comma if we are a string
                    angular.isString(ngModel.$viewValue) ? model = ngModel.$viewValue.replace(/(^,)|(,$)/g, "") : model = ngModel.$viewValue;
                    if (!model) { return; }

                    // Make sure we have an array of numbers
                    var data;
                    data = angular.isArray(model) ? model : model.split(',');

                    var change = 0;
                    var last = 0;
                    var prev = 0;
                    if (data.length > 2) {
                        var last = data[data.length-1];
                        var prev = data[data.length-2];
                        change = (last - prev);
                    }

                    //add the sparkline
                    elem.html('<span class="spark"></span> '+last+'<span title="'+prev+' to '+last+' change" class="change '+(change > 0 ? 'inc' : 'dec')+'">'+(change > 0 ? '+'+change : change)+'</span>');

                    var spark = elem.children()[0];
                    $(spark).sparkline(data, opts);
                };
            }
        }
    }
    directive.$inject=[];

    return directive;
});
