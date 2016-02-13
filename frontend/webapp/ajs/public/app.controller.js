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

    function MainCtrl( RestService, Action ){

        var self = this;
        self.actions = Action.actionsList();
        self.program = [];

        self.actionsDdConfig = {
            containment: '.grid-right',
            clone: true
        };

        self.programDdConfig = {
            containment: 'body',
            allowDuplicates: true
        };


        self.remove = function(array, index){
           array.splice(index,1);
        }
        
    }

}());