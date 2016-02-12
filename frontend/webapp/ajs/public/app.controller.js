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

    function MainCtrl( $scope, RestService, ActionsService ){

        var self = this;
        self.actions = ActionsService.actions();
        self.program = [];

        $scope.dropCallback = function(event, index, item, external){
            console.log(external);
           return external;
        };

        $scope.remove = function(array, index) {
            array = array.splice(index, 1);
            console.log("remove");
        }
        
    }

}());