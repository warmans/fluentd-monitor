define([], function(){

    function factory($sce) {

        var title = 'Fluentd Monitor';

        var prefix ='';

        return {
            getTitle: function() {
                return $sce.trustAsHtml(prefix + " " + title);
            },
            setTitlePrefix: function(newPrefix) {
                prefix = newPrefix;
            },
            setTitle: function(newTitle) {
                title = newTitle;
            }
        };
    }

    factory.$inject=['$sce'];

    return factory;
});
