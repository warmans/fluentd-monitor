define([], function(){

    function factory() {
        return {
            format: function(numBytes) {
                switch(true) {
                    case numBytes > 1024*1024*1024:
                        return (numBytes/(1024*1024*1024)).toFixed(2) + "GB";
                    case numBytes > (1024*1024):
                        return (numBytes/(1024*1024)).toFixed(2) + "MB";
                    case numBytes > 1024:
                        return (numBytes/1024).toFixed(2) + "KB";
                    default:
                        return numBytes + "B";
                }
            }
        };
    }

    factory.$inject=[];

    return factory;
});
