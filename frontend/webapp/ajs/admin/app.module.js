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
     * @name thymioCaptain.admin
     * @description This module is the one responsible for the whole Thymio Captain App.
     * It is mainly composed of controllers.
     *
     * @author Lucy Linder
     * @date     February 2016
     * @context  Thymio Captain
     */
    angular.module( 'thymioCaptain.admin',
        // dependencies
        [
            'thymioCaptain.rest',
            'ngAnimate'
        ] );

}());