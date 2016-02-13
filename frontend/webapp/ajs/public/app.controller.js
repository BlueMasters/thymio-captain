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
        $rootScope.program = [];
        History.watch('program');

        self.actionsDdConfig = {
            containment: '.grid-right',
            clone: true
        };

        self.programDdConfig = {
            containment: 'body',
            allowDuplicates: true
        };

        self.undo = function(){
            History.undo('program');
        };

        self.redo = function(){
            History.redo('program');
        };


        self.remove = function(array, index){
           array.splice(index,1);
        }
        
    }

}());