define([], function(){

    function factory($sce) {

        var title = 'Fluentd Monitor';
        var prefix ='';
        var favicon = 'fluentd-logo-black-16.png';

        return {
            getTitle: function() {
                return $sce.trustAsHtml(prefix + " " + title);
            },
            setTitlePrefix: function(newPrefix) {
                prefix = newPrefix;
            },
            setTitle: function(newTitle) {
                title = newTitle;
            },
            getFavicon: function() {
                return favicon;
            },
            setFavicon: function(newFavicon) {
                favicon = newFavicon;
            }
        };
    }

    factory.$inject=['$sce'];

    return factory;
});
