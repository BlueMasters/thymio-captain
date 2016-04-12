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

            this.color = function(){
                return ACTIONS[this.type].color;
            }
        }

        Action.actionsList = getAvailActions;
        Action.fromJson = fromJson;
        Action.prototype.toJson = asRestParam;

        // ----------------------------------------------------

        function asRestParam(){
            return {Action: this.type, Param: this.param};
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
        title: "avance",
        color: "#FFC000", // yellow
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
            descr: "--> mur"
        }, {
            id   : "UntilBlackFloor",
            descr: "--> sol noir"
        }, {
            id   : "UntilWhiteFloor",
            descr: "--> sol blanc"
        }]
    };

    ACTIONS["MoveBackward"] = {
        title: "recule",
        color: "#70AD47", // green
        args : ACTIONS["MoveForward"].args
    };

    ACTIONS["Turn"] = {
        title: "tourne",
        color: "#E76D19", // orange
        args : [{
            id   : "Right45",
            descr: "45° droite"
        }, {
            id   : "Left45",
            descr: "45° gauche"
        }, {
            id   : "Right90",
            descr: "90° droite"
        }, {
            id   : "Left90",
            descr: "90° gauche"
        }, {
            id   : "Right180",
            descr: "de 180°"
        }, {
            id   : "Right135",
            descr: "135° droite"
        }, {
            id   : "Left135",
            descr: "135° gauche"
        }]
    };

    ACTIONS["FollowLine"] = {
        title: "suit la ligne",
        color: "#41719C", // blue
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
            descr: "--> mur"
        }]
    };

    ACTIONS["SetTopColor"] = {
        title: "couleur haut",
        color: "#AA00FF", // violet
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
        title: "couleur bas",
        color: "#FF4081", // pink
        args : ACTIONS["SetTopColor"].args
    };


})();