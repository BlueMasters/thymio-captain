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

    function MainCtrl( $rootScope, ModalService, RestService, Action, History ){

        var self = this;

        self.actions = Action.actionsList();  // default actions available

        $rootScope.program = [];  // the program (set in the init())

        self.progState = 0;  // zero if in sync with the saved program
        self.notes = "";
        self.savedNotes = ""; // last saved notes

        _getCardId();    //TODO
        //self.cardId = "test";
        self.cardIdParam = {cardId: self.cardId}; //TODO

        _init();

        //##------------ functions

        self.canUndo = historyCanUndo;
        self.canRedo = historyCanRedo;
        self.undo = historyUndo;
        self.redo = histroyRedo;

        self.notesDirty = isNotesFieldDirty;

        self.save = saveCardInfos;
        self.upload = uploadProgram;
        self.run = runProgram;
        self.stop = stopProgram;

        self.dial = showRunStopDialog;


        self.contentLoaded = updateMdl;

        /* *****************************************************************
         * implementation
         * ****************************************************************/

        //##------------init


        function _init(){
            RestService.getCardData( self.cardIdParam, function( data ){
                $rootScope.program = Action.fromJson( data.program );
                _initNotes( data.notes );
                _initHistory();
                _addConfirmDialogOnClose();
                console.log( "Initialisation done: ", data );
                updateMdl();
            }, _log );
        }

        function _initNotes( notes ){
            self.notes = self.savedNotes = notes;
            updateMdl();
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

        function _addConfirmDialogOnClose(){
            // add confirmation dialog on close
            $( window ).bind( 'beforeunload', function(){
                    if( self.progState != 0 ) return 'Attention! Certains de tes changements ne sont pas' +
                        ' sauvegardés. En quittant la page, ces derniers seront perdus!';
                }
            );
        }


        function updateMdl(){
            console.log( "update mdl" );
            componentHandler.upgradeAllRegistered();
            if( self.notes )
                $( '#notesArea' ).addClass( "is-dirty" );
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

        function isNotesFieldDirty(){
            return self.notes != self.savedNotes;
        }

        //##------------ rest

        function saveCardInfos(){
            RestService.setCardData( self.cardIdParam, {notes: self.notes, program: _createProg()}, function(){
                showToast( 'Programme sauvé!' );
                self.savedNotes = self.notes;
                self.progState = 0;
            }, function(){
                showToast( 'ERREUR: le programme n\'a pu être sauvé' );
            } );  // TODO errors
        }

        function uploadProgram(){
            if( $rootScope.program.length == 0 ){
                showMessageDialog( "Pas de programme", "Il faut d'abord que tu écrives un programme... " );
            }else{
                if( self.progState != 0 ){
                    // save prog before upload
                    RestService.setCardData( self.cardIdParam, {notes: self.notes, program: _createProg()}, function(){
                        self.progState = 0;
                        self.savedNotes = self.notes;
                        _uploadProgram();

                    }, function(){
                        showToast( 'ERREUR: le programme n\'a pu être sauvé' );
                    } );  // TODO errors

                }else{
                    _uploadProgram();
                }
            }
        }


        function _uploadProgram(){

            ModalService.showModal({
                framework: "mdl",
                html: "<div class='align-center'><image src='/vendor/loading.gif' /></div>",
                cancelable: false
            });

            RestService.upload( self.cardIdParam, showRunStopDialog,
                function(){
                    ModalService.close();
                    showMessageDialog( "Pas de Thymio", "Tu n'as pas encore de Thymio attribué. Demande de l'aide à un" +
                        " animateur et réessaie." );
                } );
        }

        function runProgram(){
            RestService.run( self.cardIdParam, runStopError, runStopError );
        }

        function stopProgram(){
            RestService.stop( self.cardIdParam, runStopError, runStopError );
        }

        function runStopError(){
            ModalService.showModal( {
                    framework  : "mdl",
                    html      : '<p class="run-error-icon"><i class="material-icons">error_outline</i></p><div class="align-center">Il semble que le robot ne répond' +
                    ' plus...</div>',
                    cancelable: true
                }
            );
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
            ModalService.showModal( {
                framework: "mdl",
                title    : "Ton robot est prêt!",
                html     : '<div class="align-center"><p>Utilise les boutons suivants pour le contrôler:</p>' +
                '<div><button class="run-btn mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--accent"ng-click="inputs.run();">vas-y !</button></div>' +
                '<div><button class="run-btn mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--accent"ng-click="inputs.stop();">arrête.</button></div>' +
                '</div>',
                inputs   : {
                    run : runProgram,
                    stop: stopProgram
                },
                positive : "fermer"

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

        function _createProg(){
            var prog = [];
            angular.forEach( $rootScope.program, function( action ){
                prog.push( action.toJson() );
            } );
            return prog;
        }

        function _log( o ){
            console.log( o );
        }

        function _getCardId(){
            var m = window.location.pathname.match( '.*start/([^/#\?]*).*' );
            if( m && m.length > 1 ){
                self.cardId = m[1];
            }else{
                $( 'body' ).html( "ERREUR: pas de cardId. ACCES INTERDIT" );
            }
        }

    }

    /* *****************************************************************
     * modal
     * ****************************************************************/

    function ModalCtrl(){
    }

}());