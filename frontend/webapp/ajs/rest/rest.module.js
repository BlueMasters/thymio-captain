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
     * @name thymio.rest
     * @requires $rootScope
     * @requires $httpProvider
     * @requires $resource
     * @description
     * This module handles the interaction with the server and the app front end.
     */
    angular
        .module( 'thymioCaptain.rest', ['ngResource', 'base64'] )
        .config( function( $httpProvider, $base64 ){

            // handle the from/to base64 (program argument only)
            $httpProvider.defaults.transformRequest.unshift( function( data, headerGetter ){
                if( data && data.program ){
                    data.program = $base64.encode( JSON.stringify( data.program ) );
                }
                return data;
            } );

            $httpProvider.defaults.transformResponse.push( function( data, headerGetter ){
                if( data && data.program ){
                    data.program = JSON.parse( $base64.decode( data.program ) );
                }
                return data;
            } );

        } );


}());
