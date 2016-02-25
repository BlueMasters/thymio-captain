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

        self.progState = 0;  // zero if in sync with the saved program


        _getCardId();    //TODO
        // self.cardId = "test";

        self.cardIdParam = {cardId: self.cardId}; //TODO

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

        self.dial = showRunStopDialog;

        self.test = function(){
            console.log( 'test' );
        };

        /* *****************************************************************
         * implementation
         * ****************************************************************/

        //##------------init


        function _init(){
            RestService.getCardData( self.cardIdParam, function( data ){
                $rootScope.program = Action.fromJson( data.program );
                _initHistory();
                console.log( "Initialisation done: ", data );
            }, _log );
        }

        function _initHistory(){
            // progState == 0 only if the displayed program matches
            // the one saved in the cloud.
            var h = History.watch( 'program' );
            h.addChangeHandler( 'ch', function(){
                self.progState++;
            } );
            h.addUndoHandler( 'uh', function(){
                self.progState--;
            } );
            h.addRedoHandler( 'rh', function(){
                self.progState++;
            } );
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
                showToast( 'Programme sauvé!' );
                self.progState = 0;
            }, function(){
                showToast( 'ERREUR: le programme n\'a pu être sauvé' );
            } );  // TODO errors
        }

        function uploadProgram(){
            RestService.upload( self.cardIdParam, showRunStopDialog,
                function(){
                    showMessageDialog( "Pas de Thymio", "Tu n'as pas encore de Thymio attribué. Demande de l'aide à un" +
                        " animateur et réessaie" );
                } );
        }

        function runProgram(){
            RestService.run( self.cardIdParam, _log, _log );
        }

        function stopProgram(){
            RestService.stop( self.cardIdParam, _log, _log );
        }


        //##------------ dialogs and toasts

        function showMessageDialog( title, msg ){
            showDialog( {
                title   : title,
                text    : msg,
                positive: {
                    title: 'Ok'
                }
            } );
        }

        function showRunStopDialog(){
            showDialog( {
                customContent: '#runDialogContent',
                positive     : {
                    title  : 'fermer',
                    onClick: self.stop
                }
            } );
        }


        function showToast( message ){
            $( '.mdl-js-snackbar' )[0].MaterialSnackbar.showSnackbar(
                {
                    message: message,
                    timeout: 2000
                }
            );
        }

        //##------------ utils


        function _log( o ){
            console.log( o );
        }

        function _getCardId(){
            var m  = window.location.pathname.match('.*start/([^/#\?]*).*');
            if(m && m.length > 1){
                self.cardId = m[1];
            }else{
                $('body').html("ERREUR: pas de cardId. ACCES INTERDIT");
            }
        }

    }

}());