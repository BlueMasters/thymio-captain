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
        .filter( 'robotsFilter', robotsFilter );

    // --------------------------

    function robotsFilter(){

        return function( input, filterType ){

            var filtered = [];
            if(!filterType) return input;

            switch( filterType ){

                case "used":
                    angular.forEach( input, function( r ){
                        if( r.cardId ) filtered.push( r );
                    } );
                    break;

                case "free":
                    angular.forEach( input, function( r ){
                        if( !r.cardId ) filtered.push( r );
                    } );
                    break;

                default:
                    filtered = input;
            }

            return filtered;

        };

    }

}());