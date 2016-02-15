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

        self.actions = Action.actionsList();  // default actions available

        $rootScope.program = [];  // the program (set in the init())

        self.cardIdParam = {cardId: "jacques"}; //TODO

        _init();

        //##------------ functions

        self.canUndo = historyCanUndo;
        self.canRedo = historyCanRedo;
        self.undo = historyUndo;
        self.redo = histroyRedo;

        self.save = saveCardInfos;
        self.upload = uploadProgram;
        self.run = runProgram;
        self.stop = stopProgram;

        /* *****************************************************************
         * implementation
         * ****************************************************************/

        //##------------init

        function _init(){
            self.cardIdParam = {cardId: "jacques"};
            RestService.getCardData( self.cardIdParam, function( data ){
                $rootScope.program = Action.fromJson( data.program );
                History.watch( 'program' );
                console.log( "Initialisation done: ", data );
            }, _log );
        }

        //##------------ drag and drop

        self.actionsDdConfig = {
            containment: '.grid-right',
            clone      : true
        };

        self.programDdConfig = {
            containment    : 'body',
            allowDuplicates: true
        };

        //##------------undo/redo

        function historyUndo(){
            History.undo( 'program' );
        }

        function histroyRedo(){
            History.redo( 'program' );
        }

        function historyCanRedo(){
            return History.canRedo( 'program' );
        }

        function historyCanUndo(){
            return History.canUndo( 'program' );
        }

        //##------------ rest

        function saveCardInfos(){
            var prog = [];
            angular.forEach( $rootScope.program, function( action ){
                prog.push( action.toJson() );
            } );
            RestService.setCardData( self.cardIdParam, {program: prog}, function(){
                showToast( 'Prgramme sauvé!' );
            }, function(){
                showToast( 'ERREUR: le programme n\'a pu être sauvé' );
            } );  // TODO errors
        }

        function uploadProgram(){
            RestService.upload( self.cardIdParam, _log, _log );
        }

        function runProgram(){
            RestService.run( self.cardIdParam, _log, _log );
        }

        function stopProgram(){
            RestService.stop( self.cardIdParam, _log, _log );
        }

        //##------------ utils

        function showToast( message ){
            $( '.mdl-js-snackbar' )[0].MaterialSnackbar.showSnackbar(
                {
                    message: message,
                    timeout: 2000
                }
            );
        }


        function _log( o ){
            console.log( o );
        }

    }

}());