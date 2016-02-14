(function(){
    angular.
        module( 'thymioCaptain.actions' )
        .factory( 'Action', ActionFactory );
    // --------------------------


    function ActionFactory(){

        function Action( type, param ){
            this.type = type;
            this.param = param || null;

            this.args = function(){
                return ACTIONS[this.type].args;
            };

            this.title = function(){
                return ACTIONS[this.type].title;
            };
        }

        Action.actionsList = getAvailActions;
        Action.fromJson = fromJson;
        Action.prototype.asJson = asJson;

        // ----------------------------------------------------

        function asJson(){
            return JSON.stringify( {Action: this.type, Param: this.param} );
        }

        function getAvailActions(){
            var array = [];
            for( var key in ACTIONS ){
                array.push( new Action( key ) );
            }

            return array;
        }


        function fromJson( obj ){
            if( typeof obj === "string" ) obj = JSON.parse( obj );
            var array = [];
            for( var i in obj ){
                array.push( new Action( obj[i].Action, obj[i].Param ) );
            }

            return array;
        }

        // ----------------------------------------------------

        return Action;

    }

    // ----------------------------------------------------

    var ACTIONS = [];

    ACTIONS["MoveForward"] = {
        title: "avancer",
        args : [{
            id   : "10cm",
            descr: "de 10 cm"
        }, {
            id   : "20cm",
            descr: "de 20 cm"
        }, {
            id   : "50cm",
            descr: "de 50 cm"
        }, {
            id   : "UntilWall",
            descr: "jusqu'au mur"
        }, {
            id   : "UntilBlackFloor",
            descr: "jusqu'à ce que le sol soit noir"
        }, {
            id   : "UntilWhiteFloor",
            descr: "jusqu'à ce que le sol soit blanc"
        }]
    };

    ACTIONS["MoveBackward"] = {
        title: "reculer",
        args : ACTIONS["MoveForward"].args
    };

    ACTIONS["Turn"] = {
        title: "tourner",
        args : [{
            id   : "Right45",
            descr: "de 45° sur la droite"
        }, {
            id   : "Right90",
            descr: "de 90° sur la droite"
        }, {
            id   : "Right135",
            descr: "de 135° sur la droite"
        }, {
            id   : "Right180",
            descr: "de 180°"
        }, {
            id   : "Left45",
            descr: "de 45° sur la gauche"
        }, {
            id   : "Left90",
            descr: "de 90° sur la gauche"
        }, {
            id   : "Left135",
            descr: "de 135° sur la gauche"
        }]
    };

    ACTIONS["FollowLine"] = {
        title: "suivre la ligne",
        args : [{
            id   : "10cm",
            descr: "sur 10 cm"
        }, {
            id   : "20cm",
            descr: "sur 20 cm"
        }, {
            id   : "50cm",
            descr: "sur 50 cm"
        }, {
            id   : "UntilWall",
            descr: "jusqu'au mur"
        }]
    };

    ACTIONS["SetTopColor"] = {
        title: "couleur du dessus",
        args : [{
            id   : "off",
            descr: "éteindre"
        }, {
            id   : "red",
            descr: "rouge"
        }, {
            id   : "blue",
            descr: "bleu"
        }, {
            id   : "green",
            descr: "vert"
        }, {
            id   : "pink",
            descr: "rose"
        }, {
            id   : "orange",
            descr: "orange"
        }, {
            id   : "white",
            descr: "blanc"
        }]
    };

    ACTIONS["SetBottomColor"] = {
        title: "couleur du dessous",
        args : ACTIONS["SetTopColor"].args
    };


})();