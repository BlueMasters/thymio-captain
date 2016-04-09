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
     * @ngdoc overview
     * @name thymioCaptain.app
     * @description This module is the one responsible for the whole Thymio Captain App.
     * It is mainly composed of controllers.
     * @requires  thymioCaptain.rest
     * @requires  thymio.modals
     * @requires  ngAnimate
     * @requires  toaster
     *
     * @author Lucy Linder
     * @date     February 2016
     * @context  Thymio Captain
     */
    angular.module( 'thymioCaptain.app',
        // dependencies
        [
            'thymioCaptain.rest',
            'thymioCaptain.actions',
            'ngAnimate',
            'as.sortable',
            'decipher.history',
            'derlin.modals',
            'ui.materialize'
        ] );

}());