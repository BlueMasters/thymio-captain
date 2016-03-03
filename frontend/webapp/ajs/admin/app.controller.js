/*
 * @author   Lucy Linder <lucy.derlin@gmail.com>
 * @date     February 2016
 * @context  Thymio Captain
 *
 * Copyright 2016 BlueMagic. All rights reserved.
 * Use of this source code is governed by an Apache 2 license
 * that can be found in the LICENSE file.
 */
(function(){

    /**
     * @ngdoc controller
     * @name thymioCaptain.admin.MainCtrl
     *
     * @description
     * Main controller
     */
    angular
        .module( 'thymioCaptain.admin' )
        .controller( 'MainCtrl', MainCtrl );

    // --------------------------

    function MainCtrl( $scope, RestService ){

        // disgusting fix for mdl + ng-include
        $scope.contentLoaded = function(){
            componentHandler.upgradeAllRegistered();
        };


        var self = this;

        self.cardId = null;
        self.currentRobotUrl = null;


        self.ping = ping;
        self.associate = associate;
        self.dissociate = dissociate;

        _getCardId();
        _init();


        /* *****************************************************************
         * implementation
         * ****************************************************************/

        function _init(){
            RestService.getRobots( function( data ){
                self.robots = data;
                // check if card already associated
                if( self.cardId ){
                    var cr = self.robots.filter( function( r ){
                        return r.cardId == self.cardId;
                    } );
                    self.currentRobotUrl = cr.length > 0 ? cr[0].url : null;
                }
            } );
        }

        function associate( robot ){
            if( !self.cardId || self.currentRobotUrl ){
                createShowToast( "No card or card already associated..." )();
                return;
            }
            RestService.associateThymio( {cardId: self.cardId, name: robot.name},
                function(){
                    _init();
                    createShowToast( "Associated !" )();
                }, createShowToast( "Something went wrong.", true ) );
        }

        function dissociate( robot ){
            if( !robot.cardId ){
                createShowToast( "Not associated..." )();
                return;
            }
            RestService.dissociateThymio( {name: robot.name}, function(){
                _init();
                createShowToast( "Dissociated !" )();
            }, createShowToast( "Something went wrong.", true ) );
        }


        function ping( robot ){
            RestService.pingRobot( {name: robot.name}, createShowToast( "PING successful" ), //
                createShowToast( "PING failed.", true ) );
        }

        // ----------------------------------------------------

        function createShowToast( msg, logAdditional ){
            function show( data ){
                if( logAdditional && data ){
                    console.log( data );
                    msg += "\n" + data.status + " " + data.statusText;
                }

                $( '.mdl-js-snackbar' )[0].MaterialSnackbar.showSnackbar(
                    {
                        message: msg,
                        timeout: 2000
                    }
                );
            }

            return show;
        }


        function _getCardId(){
            var m  = window.location.pathname.match('.*start/([^/#\?]*).*');
            if(m && m.length > 1){
                self.cardId = m[1];
            }
        }
    }

}());