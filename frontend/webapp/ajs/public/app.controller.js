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
     * @name thymioCaptain.app.MainCtrl
     *
     * @description
     * Main controller
     */
    angular
        .module( 'thymioCaptain.app' )
        .controller( 'MainCtrl', MainCtrl );

    // --------------------------

    function MainCtrl( $rootScope, RestService, Action, History ){

        var self = this;
        self.actions = Action.actionsList();
        self.cardId = "card-id"; //TODO
        $rootScope.program = [];
        _init();

        //##------------ undo redo
        History.watch( 'program' );


        self.undo = function(){
            History.undo( 'program' );
        };

        self.redo = function(){
            History.redo( 'program' );
        };

        //##------------ drag and drop

        self.actionsDdConfig = {
            containment: '.grid-right',
            clone      : true
        };

        self.programDdConfig = {
            containment    : 'body',
            allowDuplicates: true
        };

        //##------------ rest

        self.save = function(){
            var prog = [];
            angular.forEach( $rootScope.program, function( action ){
                prog.push( action.toJson() );
            } );
            RestService.saveProgram( {id: self.cardId}, prog, function(){
                self.showToast( 'program saved' );
            }, _log );  // TODO errors
        };

        self.upload = function(){
            RestService.uploadProgram({id: self.cardId}, _log, _log);
        };

        self.run = function(){
            RestService.run({id: self.cardId}, _log, _log);
        };

        self.stop = function(){
            RestService.stop({id: self.cardId}, _log, _log);
        };

        //##------------ utils


        self.remove = function( array, index ){
            array.splice( index, 1 );
        };

        self.showToast = function( message ){
            $( '.mdl-js-snackbar' )[0].MaterialSnackbar.showSnackbar(
                {
                    message: message,
                    timeout: 2000
                }
            );
        };

        // ----------------------------------------------------


        function _init(){
            RestService.getProgram( {id: self.cardId}, function( data ){
                $rootScope.program = Action.fromJson( data );
            }, _log );
        }

        function _log( o ){
            console.log( o );
        }

    }

}());