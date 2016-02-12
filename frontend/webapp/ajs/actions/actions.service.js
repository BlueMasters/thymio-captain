(function(){
    angular.
        module( 'thymioCaptain.actions' )
        .service( 'ActionsService', ActionService );
    // --------------------------


    function ActionService(){
        return {
            actions: getActionsList,
            encode: encodeList
        };


        // ----------------------------------------------------
        function encodeList(actionsList){
            var program = [];
            angular.forEach(actionsList, function(action){
            });
        }
        // ----------------------------------------------------


        function getActionsList(){
            return [{
                id   : "MoveForward",
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
            }, {
                id   : "MoveBackward",
                title: "reculer",
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
            }, {
                id   : "Turn",
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
            }, {
                id   : "FollowLine",
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
            }, {
                id   : "SetTopColor",
                title: "changer la couleur du dessus",
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
            }, {
                id   : "SetBottomColor",
                title: "changer la couleur du dessous",
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
            }];
        }
    }

})();