define([], function(){
    function filter (bytesFormatter) {
        return bytesFormatter.format;
    }

    filter.$inject = ['bytesFormatter'];

    return filter;
});
