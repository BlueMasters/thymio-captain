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

    function MainCtrl( $rootScope, RestService ){

        var self = this;

        var robots = [];

        _init();


        /* *****************************************************************
         * implementation
         * ****************************************************************/

        function _init(){
            RestService.getRobots(function(data){
                self.robots =  data;
            });
        }


    }

}());