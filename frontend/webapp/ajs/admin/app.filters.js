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
        .filter( 'robotsFilter', robotsFilter )
        .filter( 'cardid', cardIdFilter );

    // ----------------------------------------------------

    function cardIdFilter(){

        // usage:
        // cardId = "0123456789ABC"
        // {{ cardId | cardid }}  => 012...ABC
        // {{ cardId | cardid:2 }}  => 01...bC
        // {{ cardId | cardid:200 }}  => 0123456789ABC
        return function( cardid, nbr ){

            var n = nbr || 3;
            return cardid && cardid.length > 3 + (n * 2) ?
            cardid.substr( 0, n ) + "..." + cardid.substr( -n ) : cardid;

        }
    }

    // ----------------------------------------------------

    function robotsFilter(){

        return function( input, filterType ){

            var filtered = [];
            if( !filterType ) return input;

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